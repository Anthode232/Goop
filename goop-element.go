package goop

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// Find finds the first occurrence of the given tag name,
// with or without attribute key and value specified,
// and returns a struct with a pointer to it
func (r Root) Find(args ...string) Root {
	temp, ok := findOnce(r.Pointer, args, false, false)
	if ok == false {
		if debug {
			panic("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")
		}
		return Root{Error: newError(ErrElementNotFound, fmt.Sprintf("element `%s` with attributes `%s` not found", args[0], strings.Join(args[1:], " ")))}
	}
	return Root{Pointer: temp, NodeValue: temp.Data}
}

// FindAll finds all occurrences of the given tag name,
// with or without key and value specified,
// and returns an array of structs, each having
// the respective pointers
func (r Root) FindAll(args ...string) []Root {
	temp := findAllofem(r.Pointer, args, false)
	if len(temp) == 0 {
		if debug {
			panic("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")
		}
		return []Root{}
	}
	pointers := make([]Root, 0, len(temp))
	for i := 0; i < len(temp); i++ {
		pointers = append(pointers, Root{Pointer: temp[i], NodeValue: temp[i].Data})
	}
	return pointers
}

// FindStrict finds the first occurrence of the given tag name
// only if all the values of the provided attribute are an exact match
func (r Root) FindStrict(args ...string) Root {
	temp, ok := findOnce(r.Pointer, args, false, true)
	if ok == false {
		if debug {
			panic("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")
		}
		return Root{nil, "", newError(ErrElementNotFound, fmt.Sprintf("element `%s` with attributes `%s` not found", args[0], strings.Join(args[1:], " ")))}
	}
	return Root{Pointer: temp, NodeValue: temp.Data}
}

// FindAllStrict finds all occurrences of the given tag name
// only if all the values of the provided attribute are an exact match
func (r Root) FindAllStrict(args ...string) []Root {
	temp := findAllofem(r.Pointer, args, true)
	if len(temp) == 0 {
		if debug {
			panic("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")
		}
		return []Root{}
	}
	pointers := make([]Root, 0, len(temp))
	for i := 0; i < len(temp); i++ {
		pointers = append(pointers, Root{Pointer: temp[i], NodeValue: temp[i].Data})
	}
	return pointers
}

// FindNextSibling finds the next sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindNextSibling() Root {
	nextSibling := r.Pointer.NextSibling
	if nextSibling == nil {
		if debug {
			panic("No next sibling found")
		}
		return Root{Error: newError(ErrNoNextSibling, "no next sibling found")}
	}
	return Root{Pointer: nextSibling, NodeValue: nextSibling.Data}
}

// FindPrevSibling finds the previous sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindPrevSibling() Root {
	prevSibling := r.Pointer.PrevSibling
	if prevSibling == nil {
		if debug {
			panic("No previous sibling found")
		}

		return Root{Error: newError(ErrNoPreviousSibling, "no previous sibling found")}
	}
	return Root{Pointer: prevSibling, NodeValue: prevSibling.Data}
}

// FindNextElementSibling finds the next element sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindNextElementSibling() Root {
	nextSibling := r.Pointer.NextSibling
	if nextSibling == nil {
		if debug {
			panic("No next element sibling found")
		}
		return Root{Error: newError(ErrNoNextElementSibling, "no next element sibling found")}
	}
	if nextSibling.Type == html.ElementNode {
		return Root{Pointer: nextSibling, NodeValue: nextSibling.Data}
	}
	p := Root{Pointer: nextSibling, NodeValue: nextSibling.Data}
	return p.FindNextElementSibling()
}

// FindPrevElementSibling finds the previous element sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindPrevElementSibling() Root {
	prevSibling := r.Pointer.PrevSibling
	if prevSibling == nil {
		if debug {
			panic("No previous element sibling found")
		}
		return Root{Error: newError(ErrNoPreviousElementSibling, "no previous element sibling found")}
	}
	if prevSibling.Type == html.ElementNode {
		return Root{Pointer: prevSibling, NodeValue: prevSibling.Data}
	}
	p := Root{Pointer: prevSibling, NodeValue: prevSibling.Data}
	return p.FindPrevElementSibling()
}

// Children returns all direct children of this DOME element.
func (r Root) Children() []Root {
	child := r.Pointer.FirstChild
	var children []Root
	for child != nil {
		children = append(children, Root{Pointer: child, NodeValue: child.Data})
		child = child.NextSibling
	}
	return children
}

func matchElementName(n *html.Node, name string) bool {
	return name == "" || name == n.Data
}

// Using depth first search to find the first occurrence and return
func findOnce(n *html.Node, args []string, uni bool, strict bool) (*html.Node, bool) {
	if uni == true {
		if n.Type == html.ElementNode && matchElementName(n, args[0]) {
			if len(args) > 1 && len(args) < 4 {
				for i := 0; i < len(n.Attr); i++ {
					attr := n.Attr[i]
					searchAttrName := args[1]
					searchAttrVal := args[2]
					if (strict && attributeAndValueEquals(attr, searchAttrName, searchAttrVal)) ||
						(!strict && attributeContainsValue(attr, searchAttrName, searchAttrVal)) {
						return n, true
					}
				}
			} else if len(args) == 1 {
				return n, true
			}
		}
	}
	uni = true
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p, q := findOnce(c, args, true, strict)
		if q != false {
			return p, q
		}
	}
	return nil, false
}

// Using depth first search to find all occurrences and return
func findAllofem(n *html.Node, args []string, strict bool) []*html.Node {
	var nodeLinks = make([]*html.Node, 0, 10)
	var f func(*html.Node, []string, bool)
	f = func(n *html.Node, args []string, uni bool) {
		if uni == true {
			if n.Type == html.ElementNode && matchElementName(n, args[0]) {
				if len(args) > 1 && len(args) < 4 {
					for i := 0; i < len(n.Attr); i++ {
						attr := n.Attr[i]
						searchAttrName := args[1]
						searchAttrVal := args[2]
						if (strict && attributeAndValueEquals(attr, searchAttrName, searchAttrVal)) ||
							(!strict && attributeContainsValue(attr, searchAttrName, searchAttrVal)) {
							nodeLinks = append(nodeLinks, n)
						}
					}
				} else if len(args) == 1 {
					nodeLinks = append(nodeLinks, n)
				}
			}
		}
		uni = true
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, args, true)
		}
	}
	f(n, args, false)
	return nodeLinks
}
