package internal

import (
	"fmt"
	"slices"
	"strings"
)

type unusedDomainMapItemsProblems []unusedDomainMapItemsProblem

func (ps unusedDomainMapItemsProblems) Text() string {
	var sb strings.Builder

	sb.WriteString("UNUSED DOMAIN MAP ITEMS\n\n")

	if len(ps) == 0 {
		sb.WriteString("        No problems found\n")

		return sb.String()
	}

	for i, p := range ps {
		fmt.Fprintf(&sb, "%6d. %s\n", i+1, p.domainMapItem.Source)
		fmt.Fprintf(&sb, "%9s -> %s\n", "", p.domainMapItem.Destination)
	}

	return sb.String()
}

func (ps unusedDomainMapItemsProblems) Len() int {
	return len(ps)
}

type unusedDomainMapItemsProblem struct {
	domainMapItem DomainMapItem
}

func checkUnusedDomainMapItems(dm DomainMap, rs []replacement) unusedDomainMapItemsProblems {
	used := make([]*DomainMapItem, 0, len(rs))

	for _, r := range rs {
		if r.culprit == nil {
			continue
		}

		used = append(used, r.culprit)
	}

	ps := make([]unusedDomainMapItemsProblem, 0, len(dm))
out:
	for i := range dm {
		for u := range used {
			if dm[i] == *used[u] {
				continue out
			}
		}

		ps = append(ps, unusedDomainMapItemsProblem{domainMapItem: dm[i]})
	}

	return slices.Clip(ps)
}
