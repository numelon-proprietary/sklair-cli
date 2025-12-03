package main

import (
	"bytes"
	"os"
	"path/filepath"
	"sklair/caching"
	"sklair/discovery"
	"sklair/htmlUtilities"
	"sklair/logger"
	"strings"

	"golang.org/x/net/html"
)

const ComponentsDir = "components"
const SrcDir = "src"

func main() {
	logger.InitShared(logger.LevelDebug)

	// TODO: add a function to logger which has a cool processing animation or something
	logger.Info("Discovering documents...")
	scanned, err := discovery.DocumentDiscovery(SrcDir)
	if err != nil {
		logger.Error("Could not scan documents : %s", err.Error())
		return
	}
	logger.Info("Discovering components...")
	components, err := discovery.ComponentDiscovery(ComponentsDir)
	if err != nil {
		logger.Error("Could not scan components : %s", err.Error())
		return
	}

	componentCache := caching.ComponentCache{}

	for _, filePath := range scanned.HtmlFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Error("Could not read file %s : %s", filePath, err.Error())
			return
		}

		logger.Debug("File %s : %s", filePath, string(content))

		doc, err := html.Parse(bytes.NewReader(content))
		if err != nil {
			logger.Error("Could not parse file %s : %s", filePath, err.Error())
			return
		}

		var toReplace []*html.Node

		for node := range doc.Descendants() {
			if node.Type == html.ElementNode {
				tag := strings.ToLower(node.Data)

				if !htmlTags[tag] {
					_, dynamicExists := componentCache.Dynamic[tag]
					_, staticExists := componentCache.Static[tag]

					if !(dynamicExists || staticExists) {
						componentSrc, exists := components[tag]
						if !exists {
							logger.Info("Non-standard tag found in HTML and no component present : %s, assuming JS tag", tag)
							continue
						}

						logger.Info("Processing and caching tag %s...", tag)
						c, dynamic, err := caching.Cache(ComponentsDir, componentSrc)
						if err != nil {
							logger.Error("Could not cache component %s : %s", componentSrc, err.Error())
							return
						}

						if dynamic {
							componentCache.Dynamic[tag] = c
						} else {
							componentCache.Static[tag] = c
						}
					}

					toReplace = append(toReplace, node)
				}
			}
		}

		// TODO: LEFT OFF HERE. AFTER ALL NODES DISCOVERED ETC NEED TO REPLACE
		// TODO: in the future, hash component file contents and construct local cache in .sklair directory
		// but how would we "cache" a html.Node struct?? lol

		logger.Info("Found %d tags to replace in %s", len(toReplace), filePath)

		for _, node := range toReplace {

			componentPath := filepath.Join(ComponentsDir, node.Data+".html")

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

		newWriter := bytes.NewBuffer(nil)
		err = html.Render(newWriter, doc)
		if err != nil {
			logger.Error("Could not render output : %s", err.Error())
			return
		}

		err = os.WriteFile("./src/output.html", newWriter.Bytes(), 0644)
		if err != nil {
			logger.Error("Could not write output : %s", err.Error())
			return
		}

		logger.Info("Saved to output.html")
	}
}
