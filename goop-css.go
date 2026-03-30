package goop

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// CSSSelector represents a parsed CSS selector
type CSSSelector struct {
	Type       string // "element", "class", "id", "attribute"
	Value      string // selector value
	Attribute  string // for attribute selectors
	Operator   string // for attribute selectors (=, *=, ^=, $=)
	Combinator string // " ", ">", "+", "~"
	Pseudo     string // pseudo-class name
	Index      int    // for nth-child selectors
}

// CSS finds the first element matching the CSS selector
func (r Root) CSS(selector string) Root {
	timer := startTimer("CSS: "+selector, DebugVerbose)
	defer timer.finish()

	selectors, err := parseCSSSelector(selector)
	if err != nil {
		debugLog(DebugBasic, "CSS selector parse error: %v", err)
		return Root{Error: newError(ErrElementNotFound, "invalid CSS selector: "+selector)}
	}

	element := findCSSMatch(r.Pointer, selectors)
	if element == nil {
		logDOMOperation("CSS", selector, 0)
		return Root{Error: newError(ErrElementNotFound, "no element found for selector: "+selector)}
	}

	logDOMOperation("CSS", selector, 1)
	return Root{Pointer: element, NodeValue: element.Data}
}

// CSSAll finds all elements matching the CSS selector
func (r Root) CSSAll(selector string) []Root {
	timer := startTimer("CSSAll: "+selector, DebugVerbose)
	defer timer.finish()

	selectors, err := parseCSSSelector(selector)
	if err != nil {
		debugLog(DebugBasic, "CSS selector parse error: %v", err)
		return []Root{}
	}

	elements := findAllCSSMatches(r.Pointer, selectors)
	logDOMOperation("CSSAll", selector, len(elements))

	results := make([]Root, len(elements))
	for i, elem := range elements {
		results[i] = Root{Pointer: elem, NodeValue: elem.Data}
	}
	return results
}

// parseCSSSelector parses a CSS selector string into a slice of selector parts
func parseCSSSelector(selector string) ([]CSSSelector, error) {
	var selectors []CSSSelector
	parts := strings.Fields(selector) // split by whitespace for combinators

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Handle combinators
		combinator := " " // descendant (default)
		if strings.Contains(part, ">") {
			combinator = ">"
			part = strings.Replace(part, ">", "", -1)
		} else if strings.Contains(part, "+") {
			combinator = "+"
			part = strings.Replace(part, "+", "", -1)
		} else if strings.Contains(part, "~") {
			combinator = "~"
			part = strings.Replace(part, "~", "", -1)
		}

		// Parse the selector part
		sel, err := parseCSSPart(part)
		if err != nil {
			return nil, err
		}
		sel.Combinator = combinator
		selectors = append(selectors, sel)
	}

	return selectors, nil
}

// parseCSSPart parses a single CSS selector part
func parseCSSPart(part string) (CSSSelector, error) {
	var sel CSSSelector

	// ID selector (#id)
	if strings.HasPrefix(part, "#") {
		sel.Type = "id"
		sel.Value = strings.TrimPrefix(part, "#")
		return sel, nil
	}

	// Class selector (.class)
	if strings.HasPrefix(part, ".") {
		sel.Type = "class"
		sel.Value = strings.TrimPrefix(part, ".")
		return sel, nil
	}

	// Attribute selector [attr], [attr=value], [attr*=value], etc.
	if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
		attrContent := strings.TrimSuffix(strings.TrimPrefix(part, "["), "]")

		// Simple attribute selector [attr]
		if !strings.ContainsAny(attrContent, "=^$*") {
			sel.Type = "attribute"
			sel.Attribute = attrContent
			sel.Operator = "exists"
			return sel, nil
		}

		// Attribute with operator [attr=value], [attr*=value], etc.
		for _, op := range []string{"*=", "^=", "$=", "="} {
			if strings.Contains(attrContent, op) {
				parts := strings.SplitN(attrContent, op, 2)
				if len(parts) == 2 {
					sel.Type = "attribute"
					sel.Attribute = strings.TrimSpace(parts[0])
					sel.Value = strings.TrimSpace(parts[1])
					sel.Operator = op
					return sel, nil
				}
			}
		}

		return sel, fmt.Errorf("invalid attribute selector: %s", part)
	}

	// Pseudo-class selectors (:first-child, :last-child, :nth-child(n))
	if strings.HasPrefix(part, ":") {
		pseudo := strings.TrimPrefix(part, ":")

		// Handle :nth-child(n)
		if strings.HasPrefix(pseudo, "nth-child(") && strings.HasSuffix(pseudo, ")") {
			numStr := strings.TrimSuffix(strings.TrimPrefix(pseudo, "nth-child("), ")")
			index, err := parseNthChild(numStr)
			if err != nil {
				return sel, fmt.Errorf("invalid nth-child selector: %s", part)
			}
			sel.Type = "pseudo"
			sel.Pseudo = "nth-child"
			sel.Index = index
			return sel, nil
		}

		// Simple pseudo-classes
		sel.Type = "pseudo"
		sel.Pseudo = pseudo
		return sel, nil
	}

	// Element selector
	sel.Type = "element"
	sel.Value = part
	return sel, nil
}

// parseNthChild parses nth-child expressions like "2", "odd", "even", "3n+1"
func parseNthChild(expr string) (int, error) {
	expr = strings.TrimSpace(expr)

	// Handle "odd" and "even"
	if expr == "odd" {
		return 1, nil // first odd element
	}
	if expr == "even" {
		return 2, nil // first even element
	}

	// Handle simple numbers
	if num, err := strconv.Atoi(expr); err == nil {
		if num < 1 {
			return 0, fmt.Errorf("nth-child index must be >= 1")
		}
		return num, nil
	}

	// Handle "an+b" format (simplified - only handle simple cases)
	if strings.Contains(expr, "n") {
		parts := strings.Split(expr, "n")
		if len(parts) == 2 {
			a, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			b, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			if a == 1 && b > 0 {
				return b, nil
			}
		}
	}

	return 0, fmt.Errorf("unsupported nth-child expression: %s", expr)
}

// findCSSMatch finds the first element matching CSS selectors
func findCSSMatch(node *html.Node, selectors []CSSSelector) *html.Node {
	if len(selectors) == 0 {
		return nil
	}

	// For now, implement simple selector matching without combinators
	// This is a simplified version - full CSS selector support would be more complex
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if matchesCSSSelector(child, selectors[0]) {
			if len(selectors) == 1 {
				return child
			}
			// Recursively search for remaining selectors
			if result := findCSSMatch(child, selectors[1:]); result != nil {
				return result
			}
		}
		// Continue searching in child elements
		if result := findCSSMatch(child, selectors); result != nil {
			return result
		}
	}

	return nil
}

// findAllCSSMatches finds all elements matching CSS selectors
func findAllCSSMatches(node *html.Node, selectors []CSSSelector) []*html.Node {
	var results []*html.Node

	if len(selectors) == 0 {
		return results
	}

	// Simple implementation - traverse all nodes and check if they match the first selector
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && matchesCSSSelector(n, selectors[0]) {
			if len(selectors) == 1 {
				results = append(results, n)
			} else {
				// For complex selectors, we'd need to check the full selector chain
				// Simplified: just add if it matches the first selector
				results = append(results, n)
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(node)
	return results
}

// matchesCSSSelector checks if a node matches a CSS selector
func matchesCSSSelector(node *html.Node, selector CSSSelector) bool {
	switch selector.Type {
	case "element":
		return selector.Value == "" || node.Data == selector.Value

	case "id":
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == selector.Value {
				return true
			}
		}
		return false

	case "class":
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				classes := strings.Fields(attr.Val)
				for _, class := range classes {
					if class == selector.Value {
						return true
					}
				}
			}
		}
		return false

	case "attribute":
		for _, attr := range node.Attr {
			if attr.Key == selector.Attribute {
				switch selector.Operator {
				case "exists":
					return true
				case "=":
					return attr.Val == selector.Value
				case "*=":
					return strings.Contains(attr.Val, selector.Value)
				case "^=":
					return strings.HasPrefix(attr.Val, selector.Value)
				case "$=":
					return strings.HasSuffix(attr.Val, selector.Value)
				}
			}
		}
		return false

	case "pseudo":
		return matchesPseudoSelector(node, selector)

	default:
		return false
	}
}

// matchesPseudoSelector checks if a node matches a pseudo-class selector
func matchesPseudoSelector(node *html.Node, selector CSSSelector) bool {
	switch selector.Pseudo {
	case "first-child":
		return isFirstChild(node)
	case "last-child":
		return isLastChild(node)
	case "nth-child":
		return isNthChild(node, selector.Index)
	default:
		return false
	}
}

// isFirstChild checks if node is the first child element
func isFirstChild(node *html.Node) bool {
	if node.Parent == nil {
		return false
	}

	for child := node.Parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			return child == node
		}
	}
	return false
}

// isLastChild checks if node is the last child element
func isLastChild(node *html.Node) bool {
	if node.Parent == nil {
		return false
	}

	var lastElement *html.Node
	for child := node.Parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			lastElement = child
		}
	}
	return lastElement == node
}

// isNthChild checks if node is the nth child element
func isNthChild(node *html.Node, index int) bool {
	if node.Parent == nil {
		return false
	}

	count := 0
	for child := node.Parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			count++
			if child == node {
				return count == index
			}
		}
	}
	return false
}
