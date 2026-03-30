package goop

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Attrs returns a map containing all attributes
func (r Root) Attrs() map[string]string {
	if r.Pointer.Type != html.ElementNode {
		if debug {
			panic("Not an ElementNode")
		}
		return nil
	}
	if len(r.Pointer.Attr) == 0 {
		return nil
	}
	return getKeyValue(r.Pointer.Attr)
}

// Text returns the string inside a non-nested element
func (r Root) Text() string {
	k := r.Pointer.FirstChild
checkNode:
	if k != nil && k.Type != html.TextNode {
		k = k.NextSibling
		if k == nil {
			if debug {
				panic("No text node found")
			}
			return ""
		}
		goto checkNode
	}
	if k != nil {
		r, _ := regexp.Compile(`^\s+$`)
		if ok := r.MatchString(k.Data); ok {
			k = k.NextSibling
			if k == nil {
				if debug {
					panic("No text node found")
				}
				return ""
			}
			goto checkNode
		}
		return k.Data
	}
	return ""
}

// HTML returns the HTML code for the specific element
func (r Root) HTML() string {
	var buf bytes.Buffer
	if err := html.Render(&buf, r.Pointer); err != nil {
		return ""
	}
	return buf.String()
}

// FullText returns the string inside even a nested element
func (r Root) FullText() string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.Type == html.ElementNode {
			f(n.FirstChild)
		}
		if n.NextSibling != nil {
			f(n.NextSibling)
		}
	}

	f(r.Pointer.FirstChild)

	return buf.String()
}

// attributeAndValueEquals reports when the html.Attribute attr has the same attribute name and value as from
// provided arguments
func attributeAndValueEquals(attr html.Attribute, attribute, value string) bool {
	return attr.Key == attribute && attr.Val == value
}

// attributeContainsValue reports when the html.Attribute attr has the same attribute name as from provided
// attribute argument and compares if it has the same value in its values parameter
func attributeContainsValue(attr html.Attribute, attribute, value string) bool {
	if attr.Key == attribute {
		for _, attrVal := range strings.Fields(attr.Val) {
			if attrVal == value {
				return true
			}
		}
	}
	return false
}

// Returns a key pair value (like a dictionary) for each attribute
func getKeyValue(attributes []html.Attribute) map[string]string {
	var keyvalues = make(map[string]string)
	for i := 0; i < len(attributes); i++ {
		_, exists := keyvalues[attributes[i].Key]
		if exists == false {
			keyvalues[attributes[i].Key] = attributes[i].Val
		}
	}
	return keyvalues
}
