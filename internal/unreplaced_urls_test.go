package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_checkUnreplacedURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rs   []replacement
		want unreplacedURLsProblems
	}{
		{
			name: "empty",
			rs:   nil,
			want: []unreplacedURLsProblem{},
		},
		{
			name: "unreplaced",
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "foo.com/example",
					culprit: nil,
				},
			},
			want: []unreplacedURLsProblem{
				{
					replacement: replacement{
						from:    "foo.com/example",
						to:      "foo.com/example",
						culprit: nil,
					},
				},
			},
		},
		{
			name: "unreplaced_multiple",
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "foo.com/example",
					culprit: nil,
				},
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
			},
			want: []unreplacedURLsProblem{
				{
					replacement: replacement{
						from:    "foo.com/example",
						to:      "foo.com/example",
						culprit: nil,
					},
				},
				{
					replacement: replacement{
						from:    "bar.com/example",
						to:      "bar.com/example",
						culprit: nil,
					},
				},
			},
		},
		{
			name: "no_problem_single",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
			},
			want: []unreplacedURLsProblem{},
		},
		{
			name: "no_problem",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo1.com",
						Destination: "bar1.com",
					},
				},
				{
					from: "baz.com/example",
					to:   "qux.com/example",
					culprit: &DomainMapItem{
						Source:      "baz.com",
						Destination: "qux.com",
					},
				},
			},
			want: []unreplacedURLsProblem{},
		},
		{
			name: "multiple_problems",
			rs: []replacement{
				{
					from: "foo1.com/example",
					to:   "bar1.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
				{
					from:    "foo.com/example",
					to:      "foo.com/example",
					culprit: nil,
				},
				{
					from: "baz1.com/example",
					to:   "qux1.com/example",
					culprit: &DomainMapItem{
						Source:      "baz.com",
						Destination: "qux.com",
					},
				},
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
			},
			want: []unreplacedURLsProblem{
				{
					replacement: replacement{
						from:    "foo.com/example",
						to:      "foo.com/example",
						culprit: nil,
					},
				},
				{
					replacement: replacement{
						from:    "bar.com/example",
						to:      "bar.com/example",
						culprit: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := checkUnreplacedURLs(tt.rs)

			ss := cmpopts.SortSlices(func(a, b unreplacedURLsProblem) bool {
				return a.replacement.from < b.replacement.from
			})

			au := cmp.AllowUnexported(unreplacedURLsProblem{}, replacement{})

			if diff := cmp.Diff(tt.want, got, ss, au); diff != "" {
				t.Errorf("checkUnreplacedURLs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
