package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_checkUnusedDomainMapItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dm   DomainMap
		rs   []replacement
		want unusedDomainMapItemsProblems
	}{
		{
			name: "empty",
			dm:   nil,
			rs:   nil,
			want: []unusedDomainMapItemsProblem{},
		},
		{
			name: "no_replacement",
			dm: DomainMap{
				{"foo.com", "bar.com"},
			},
			rs: nil,
			want: []unusedDomainMapItemsProblem{
				{domainMapItem: DomainMapItem{"foo.com", "bar.com"}},
			},
		},
		{
			name: "unused",
			dm: DomainMap{
				{"foo.com", "bar.com"},
			},
			rs: []replacement{
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
			},
			want: []unusedDomainMapItemsProblem{
				{domainMapItem: DomainMapItem{"foo.com", "bar.com"}},
			},
		},
		{
			name: "unused_multiple",
			dm: DomainMap{
				{"foo.com", "bar.com"},
				{"bax.com", "qux.com"},
			},
			rs: []replacement{
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
			},
			want: []unusedDomainMapItemsProblem{
				{domainMapItem: DomainMapItem{"foo.com", "bar.com"}},
				{domainMapItem: DomainMapItem{"bax.com", "qux.com"}},
			},
		},
		{
			name: "no_problem_single",
			dm: DomainMap{
				{"foo.com", "bar.com"},
			},
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "bar.com/example",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
			},
			want: []unusedDomainMapItemsProblem{},
		},
		{
			name: "no_problem_repeated",
			dm: DomainMap{
				{"foo.com", "bar.com"},
			},
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "bar.com/example",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
				{
					from:    "foo.com/example2",
					to:      "bar.com/example2",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
			},
			want: []unusedDomainMapItemsProblem{},
		},
		{
			name: "no_problem",
			dm: DomainMap{
				{"foo.com", "bar.com"},
				{"baz.com", "qux.com"},
			},
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "bar.com/example",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
				{
					from:    "foo.com/example2",
					to:      "bar.com/example2",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
				{
					from:    "baz.com/example",
					to:      "qux.com/example",
					culprit: &DomainMapItem{"baz.com", "qux.com"},
				},
			},
			want: []unusedDomainMapItemsProblem{},
		},

		{
			name: "multiple_problems",
			dm: DomainMap{
				{"foo.com", "bar.com"},
				{"quux1.com", "quux2.com"},
				{"baz.com", "qux.com"},
				{"quux3.com", "quux4.com"},
			},
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "bar.com/example",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
				{
					from:    "foo.com/example2",
					to:      "bar.com/example2",
					culprit: &DomainMapItem{"foo.com", "bar.com"},
				},
				{
					from:    "baz.com/example",
					to:      "qux.com/example",
					culprit: &DomainMapItem{"baz.com", "qux.com"},
				},
			},
			want: []unusedDomainMapItemsProblem{
				{domainMapItem: DomainMapItem{"quux1.com", "quux2.com"}},
				{domainMapItem: DomainMapItem{"quux3.com", "quux4.com"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := checkUnusedDomainMapItems(tt.dm, tt.rs)

			ss := cmpopts.SortSlices(func(a, b unusedDomainMapItemsProblem) bool {
				return a.domainMapItem.Source < b.domainMapItem.Source
			})

			au := cmp.AllowUnexported(unusedDomainMapItemsProblem{})

			if diff := cmp.Diff(tt.want, got, ss, au); diff != "" {
				t.Errorf("checkUnusedDomainMapItems() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
