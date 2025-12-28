package building

import (
	"sklair/htmlUtilities"
	"sort"
)

// TODO: if no charset tag found, add one default utf8

func OptimiseHead(segmented []*HeadSegment) []*HeadSegment {
	segmented = deduplicateSegmented(segmented)
	orderByPriority(segmented)
	return segmented
}

func deduplicateSegmented(segments []*HeadSegment) []*HeadSegment {
	seenHead := make(map[uint64]struct{})

	out := segments[:0]

	for _, s := range segments {
		if s.IsOrderingBarrier {
			out = append(out, s)
			continue
		}

		key := htmlUtilities.WeakHashNode(s.Nodes[0])
		if key == 0 {
			out = append(out, s)
			continue
		}

		if _, seen := seenHead[key]; seen {
			continue // drop the duplicate entirely, i.e. don't append it back to the reused array (see "out := segments[:0]")
		}

		seenHead[key] = struct{}{}
		out = append(out, s)
	}

	return out
}

func orderByPriority(segments []*HeadSegment) {
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].TreatAsTag < segments[j].TreatAsTag
	})
}
