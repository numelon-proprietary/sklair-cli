package priorities

import (
	"golang.org/x/net/html"
)

type SegmentType int

const (
	Charset SegmentType = iota

	HttpEquiv
	Viewport // theme-color is also treated as Viewport

	ResourceHint // preload, prefetch, dns-prefetch
	PreventFOUC  // private

	Title

	// TODO: add "Icon", and it is just below title

	Stylesheet
	Script

	Analytics // marked only by the user

	DontCare

	Generator // private
)

func StrToSegment(s string) SegmentType {
	switch s {
	case "charset":
		return Charset
	case "http-equiv":
		return HttpEquiv
	case "viewport":
		return Viewport
	case "resource-hint":
		return ResourceHint
		//case "prevent-fouc":
		//	return PreventFOUC // private
	case "title":
		return Title
	case "stylesheet":
		return Stylesheet
	case "script":
		return Script
	case "analytics":
		return Analytics
	case "dont-care":
		return DontCare
	//case "generator":
	//	return Generator // private

	default:
		return -1
	}
}

func ClassifySegment(n *html.Node) SegmentType {
	if n.Data == "meta" {
		for _, a := range n.Attr {
			if a.Key == "charset" {

				return Charset

			} else if a.Key == "http-equiv" {

				return HttpEquiv

			} else if a.Key == "name" {

				switch a.Val {
				case "viewport", "theme-color":
					return Viewport
				}

			}
		}
	} else if n.Data == "link" {
		for _, a := range n.Attr {
			// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/rel/dns-prefetch
			if a.Key == "rel" {
				switch a.Val {
				case "preload", "preconnect", "prefetch", "dns-prefetch", "modulepreload":
					return ResourceHint
				default:
					return Stylesheet // yes, just treat it as a stylesheet even if rel != stylesheet
				}
			}
		}
	} else if n.Data == "title" {
		return Title
	} else if n.Data == "script" {
		return Script
	}

	return DontCare
}
