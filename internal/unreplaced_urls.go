package internal

import (
	"fmt"
	"slices"
	"strings"
)

type unreplacedURLsProblems []unreplacedURLsProblem

func (ps unreplacedURLsProblems) Text() string {
	var sb strings.Builder

	sb.WriteString("UNREPLACED URLS\n\n")

	if len(ps) == 0 {
		sb.WriteString("        No problems found\n")

		return sb.String()
	}

	for i, p := range ps {
		fmt.Fprintf(&sb, "%6d. %s\n", i+1, p.replacement.from)
	}

	return sb.String()
}

func (ps unreplacedURLsProblems) Len() int {
	return len(ps)
}

type unreplacedURLsProblem struct {
	replacement replacement
}

func checkUnreplacedURLs(rs []replacement) unreplacedURLsProblems {
	ps := make([]unreplacedURLsProblem, 0, len(rs))

	for _, r := range rs {
		if r.from != r.to {
			continue
		}

		ps = append(ps, unreplacedURLsProblem{replacement: r})
	}

	return slices.Clip(ps)
}
