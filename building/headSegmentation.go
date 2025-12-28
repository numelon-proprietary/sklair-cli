package building

import (
	"errors"
	"fmt"
	"sklair/building/directives"
	"sklair/building/priorities"

	"golang.org/x/net/html"
)

type HeadSegment struct {
	Nodes             []*html.Node
	TreatAsTag        priorities.SegmentType // the tag to treat this segment as, if IsOrderingBarrier is true
	IsOrderingBarrier bool                   // whether this segment is an ordering barrier
}

func SegmentHead(head *html.Node) ([]*HeadSegment, error) {
	var segments []*HeadSegment

	var currentBlock *HeadSegment
	removeMode := false
	for c := head.FirstChild; c != nil; c = c.NextSibling {
		// --------------------------------------------------
		// remove directives
		// --------------------------------------------------
		//htmlUtilities.DumpNode(c)
		if directives.IsRemoveStart(c) {
			if removeMode {
				return nil, errors.New("nested remove directive")
			}
			removeMode = true
			continue
		}

		if directives.IsRemoveEnd(c) {
			if !removeMode {
				return nil, errors.New("remove directive declared without a start")
			}
			removeMode = false
			continue
		}

		if removeMode {
			continue
		}
		//if removeMode && directives.IsRemoveEnd(c) {
		//	removeMode = false
		//}
		//if removeMode && directives.IsRemoveStart(c) {
		//	return nil, errors.New("nested sklair:remove")
		//}
		//if removeMode || directives.IsRemoveStart(c) {
		//	removeMode = true
		//	continue
		//}

		// --------------------------------------------------
		// ordering barriers
		// --------------------------------------------------
		if ok, tag := directives.IsOBStart(c); ok {
			//logger.Debug("Ordering barrier start")
			//fmt.Println("Hi")
			//htmlUtilities.DumpNode(c)
			//fmt.Printf(logger.Red)
			if currentBlock != nil {
				fmt.Println(currentBlock)
				return nil, errors.New("nested ordering barrier")
			}
			//fmt.Printf("%s`%s`\n", logger.Red, tag)
			segmentType := priorities.StrToSegment(tag)
			if segmentType == -1 {
				return nil, fmt.Errorf(`unknown ordering barrier type "%s"`, tag)
			}
			currentBlock = &HeadSegment{
				TreatAsTag:        segmentType,
				IsOrderingBarrier: true,
			}
			continue
		}

		if directives.IsOBEnd(c) {
			if currentBlock == nil {
				return nil, errors.New("end of an ordering barrier declared without a start")
			}
			segments = append(segments, currentBlock)
			currentBlock = nil
			continue
		}

		// --------------------------------------------------
		// regular node handling
		// --------------------------------------------------
		if currentBlock != nil {
			currentBlock.Nodes = append(currentBlock.Nodes, c)
		} else if c.Type == html.ElementNode {
			segments = append(segments, &HeadSegment{
				Nodes:             []*html.Node{c},
				TreatAsTag:        priorities.ClassifySegment(c),
				IsOrderingBarrier: false,
			})
		}
	}

	if currentBlock != nil {
		return nil, errors.New("unclosed ordering barrier")
	}
	if removeMode {
		return nil, errors.New("unclosed remove directive")
	}

	return segments, nil
}
