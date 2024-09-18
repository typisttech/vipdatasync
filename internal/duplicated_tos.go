package internal

import (
	"fmt"
	"slices"
	"strings"
)

type duplicatedTosProblems []duplicatedTosProblem

func (ps duplicatedTosProblems) Text() string {
	var sb strings.Builder

	sb.WriteString("DUPLICATED TOS\n\n")

	if len(ps) == 0 {
		sb.WriteString("        No problems found\n")

		return sb.String()
	}

	for i, p := range ps {
		fmt.Fprintf(&sb, "%6d. %s\n", i+1, p.replacements[0].to)

		for _, r := range p.replacements {
			fmt.Fprintf(&sb, "%8s - %s\n", "", r.from)

			if r.culprit == nil {
				continue
			}

			fmt.Fprintf(&sb, "%10s %s\n", "", r.culprit.Source)
			fmt.Fprintf(&sb, "%12s -> %s\n", "", r.culprit.Destination)
		}
	}

	return sb.String()
}

func (ps duplicatedTosProblems) Len() int {
	return len(ps)
}

type duplicatedTosProblem struct {
	replacements []replacement
}

func checkDuplicatedTos(rs []replacement) duplicatedTosProblems {
	byTo := make(map[string][]replacement, len(rs))

	for _, r := range rs {
		byTo[r.to] = append(byTo[r.to], r)
	}

	ps := make([]duplicatedTosProblem, 0, len(byTo))

	for _, prs := range byTo {
		if len(prs) <= 1 {
			continue
		}

		ps = append(ps, duplicatedTosProblem{replacements: prs})
	}

	return slices.Clip(ps)
}
