package goop

import "golang.org/x/net/html"

// Root is a structure containing a pointer to an html node, the node value, and an error variable to return an error if one occurred
type Root struct {
	Pointer   *html.Node
	NodeValue string
	Error     error
}
