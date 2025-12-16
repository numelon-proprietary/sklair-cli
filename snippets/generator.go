package snippets

import "golang.org/x/net/html"

var Generator = &html.Node{
	Type: html.ElementNode,
	Data: "meta",
	Attr: []html.Attribute{
		{Key: "name", Val: "generator"},
		{Key: "content", Val: "Sklair - https://sklair.numelon.com"}, // TODO: include sklair version here
	},
}
