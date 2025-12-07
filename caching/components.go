package caching

import (
	"bytes"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

type Component struct {
	Node    *html.Node
	Dynamic bool // whether the component contains any dynamic <lua> tags
}

type ComponentCache struct {
	Static  map[string]*Component
	Dynamic map[string]*Component
}

func Cache(source string, name string) (*Component, bool, error) {
	path := filepath.Join(source, name+".html")

	//if _, err := os.Stat(path); err != nil {
	//	return nil, err
	//}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}

	// this is VERY naive but it actually works, we simply check for an opening lua tag
	hasLua := bytes.Contains(f, []byte("<lua"))
	component, err := html.Parse(bytes.NewReader(f))
	if err != nil {
		return nil, false, err
	}

	// even though components are usually bare (without doctype, head, body, etc), we still need to find the "body" (bc parsed)
	//body := component.FirstChild
	//for body != nil && body.Data != "html" {
	//	body = body.NextSibling
	//}
	//if body != nil {
	//	body = body.FirstChild
	//	for body != nil && body.Data != "body" {
	//		body = body.NextSibling
	//	}
	//}
	//
	//if (body != nil) && (body.FirstChild != nil) {
	//	component = body.FirstChild
	//}

	// old code before refactoring
	/*
		component, err := html.Parse(bytes.NewReader(f))
						if err != nil {
							logger.Error("Could not parse component %s : %s", componentPath, err.Error())
							return
						}

						// even though components are usually bare (without doctype, head, body, etc), we still need to find the "body" (bc parsed)
						body := component.FirstChild
						for body != nil && body.Data != "html" {
							body = body.NextSibling
						}
						if body != nil {
							body = body.FirstChild
							for body != nil && body.Data != "body" {
								body = body.NextSibling
							}
						}

						if body != nil {
							parent := node.Parent
							if parent != nil {
								for child := body.FirstChild; child != nil; child = child.NextSibling {
									parent.InsertBefore(htmlUtilities.Clone(child), node)
								}
								parent.RemoveChild(node)
							}
						}
	*/

	// TODO: make a new struct for components which includes a head section and a body section
	// for head, perform deduplication when multiple components in same document share head stuff
	// for body, just insert as usual

	return &Component{component, hasLua}, hasLua, nil
}

/*
for _, node := range toReplace {

			componentPath := filepath.Join("components", node.Data+".html")

			if _, err := os.Stat(componentPath); err != nil {
				logger.Error("Could not find component %s : %s", componentPath, err.Error())
				return
			} else {
				f, err := os.ReadFile(componentPath)
				if err != nil {
					logger.Error("Could not read component %s : %s", componentPath, err.Error())
					return
				}

				component, err := html.Parse(bytes.NewReader(f))
				if err != nil {
					logger.Error("Could not parse component %s : %s", componentPath, err.Error())
					return
				}

				// even though components are usually bare (without doctype, head, body, etc), we still need to find the "body" (bc parsed)
				body := component.FirstChild
				for body != nil && body.Data != "html" {
					body = body.NextSibling
				}
				if body != nil {
					body = body.FirstChild
					for body != nil && body.Data != "body" {
						body = body.NextSibling
					}
				}

				if body != nil {
					parent := node.Parent
					if parent != nil {
						for child := body.FirstChild; child != nil; child = child.NextSibling {
							parent.InsertBefore(htmlUtilities.Clone(child), node)
						}
						parent.RemoveChild(node)
					}
				}
			}
		}
*/
