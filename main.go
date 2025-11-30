package main

import (
	"bytes"
	"fmt"
	"os"
	"sklair/logger"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	logger.InitShared(logger.LevelDebug)

	content, err := os.ReadFile("test.html")
	if err != nil {
		panic(err)
	}

	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		panic(err)
	}

	var toReplace []*html.Node

	for node := range doc.Descendants() {
		if node.Type == html.ElementNode {
			tag := strings.ToLower(node.Data)

			if !htmlTags[tag] {
				toReplace = append(toReplace, node)
			}
		}
	}

	fmt.Println(toReplace)
	logger.Info("Found %d tags to replace", len(toReplace))
}
