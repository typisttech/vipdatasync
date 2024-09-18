package internal

import (
	"fmt"
	"slices"
	"strings"
)

type duplicatedDestinationsProblems []duplicatedDestinationsProblem

func (ps duplicatedDestinationsProblems) Text() string {
	var sb strings.Builder

	sb.WriteString("DUPLICATED DESTINATIONS\n\n")

	if len(ps) == 0 {
		sb.WriteString("        No problems found\n")

		return sb.String()
	}

	for i, p := range ps {
		fmt.Fprintf(&sb, "%6d. %s\n", i+1, p.domainMap[0].Destination)

		for _, item := range p.domainMap {
			fmt.Fprintf(&sb, "%8s - %s\n", "", item.Source)
		}
	}

	return sb.String()
}

func (ps duplicatedDestinationsProblems) Len() int {
	return len(ps)
}

type duplicatedDestinationsProblem struct {
	domainMap DomainMap
}

func checkDuplicatedDestinations(dm DomainMap) duplicatedDestinationsProblems {
	byDest := make(map[string][]DomainMapItem, len(dm))

	for _, dmi := range dm {
		byDest[dmi.Destination] = append(byDest[dmi.Destination], dmi)
	}

	ps := make([]duplicatedDestinationsProblem, 0, len(byDest))

	for _, pdm := range byDest {
		if len(pdm) <= 1 {
			continue
		}

		if len(pdm) == 2 && pdm[0].Source == "www."+pdm[1].Source {
			continue
		}

		ps = append(ps, duplicatedDestinationsProblem{domainMap: pdm})
	}

	return slices.Clip(ps)
}
