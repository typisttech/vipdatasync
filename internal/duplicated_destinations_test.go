package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_checkDuplicatedDestinations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dm   DomainMap
		want duplicatedDestinationsProblems
	}{
		{
			name: "empty",
			dm:   nil,
			want: []duplicatedDestinationsProblem{},
		},
		{
			name: "single",
			dm:   DomainMap{{"foo", "bar"}},
			want: []duplicatedDestinationsProblem{},
		},
		{
			name: "no_problem",
			dm: DomainMap{
				{"foo", "bar"},
				{"baz", "qux"},
			},
			want: []duplicatedDestinationsProblem{},
		},
		{
			name: "duplicated",
			dm: DomainMap{
				{"foo1", "bar"},
				{"foo2", "bar"},
			},
			want: []duplicatedDestinationsProblem{
				{
					domainMap: DomainMap{
						{"foo1", "bar"},
						{"foo2", "bar"},
					},
				},
			},
		},
		{
			name: "duplicated_exactly",
			dm: DomainMap{
				{"foo", "bar"},
				{"foo", "bar"},
			},
			want: []duplicatedDestinationsProblem{
				{
					domainMap: DomainMap{
						{"foo", "bar"},
						{"foo", "bar"},
					},
				},
			},
		},
		{
			name: "multiple_problems",
			dm: DomainMap{
				{"foo1", "bar1"},
				{"baz", "qux"},
				{"baz1", "qux2"},
				{"foo2", "bar2"},
				{"foo3", "bar1"},
				{"foo4", "bar2"},
			},
			want: []duplicatedDestinationsProblem{
				{
					domainMap: DomainMap{
						{"foo1", "bar1"},
						{"foo3", "bar1"},
					},
				},
				{
					domainMap: DomainMap{
						{"foo2", "bar2"},
						{"foo4", "bar2"},
					},
				},
			},
		},
		{
			name: "www_good",
			dm: DomainMap{
				{"www.foo", "bar"},
				{"foo", "bar"},
			},
			want: []duplicatedDestinationsProblem{},
		},
		{
			name: "www_bad",
			dm: DomainMap{
				{"foo", "bar"},
				{"www.foo", "bar"},
			},
			want: []duplicatedDestinationsProblem{
				{
					domainMap: DomainMap{
						{"foo", "bar"},
						{"www.foo", "bar"},
					},
				},
			},
		},
		{
			name: "www_extra",
			dm: DomainMap{
				{"www.foo", "bar"},
				{"foo", "bar"},
				{"baz", "bar"},
			},
			want: []duplicatedDestinationsProblem{
				{
					domainMap: DomainMap{
						{"www.foo", "bar"},
						{"foo", "bar"},
						{"baz", "bar"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := checkDuplicatedDestinations(tt.dm)

			ss := cmpopts.SortSlices(func(a, b duplicatedDestinationsProblem) bool {
				return a.domainMap[0].Destination < b.domainMap[0].Destination
			})
			au := cmp.AllowUnexported(duplicatedDestinationsProblem{})

			if diff := cmp.Diff(tt.want, got, ss, au); diff != "" {
				t.Errorf("checkDuplicatedDestinations() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
