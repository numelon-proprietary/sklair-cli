package directives

import (
	"strings"

	"golang.org/x/net/html"
)

func IsRemoveStart(n *html.Node) bool {
	return n.Type == html.CommentNode && strings.TrimSpace(n.Data) == "sklair:remove"
}

func IsRemoveEnd(n *html.Node) bool {
	return n.Type == html.CommentNode && strings.TrimSpace(n.Data) == "sklair:remove-end"
}
