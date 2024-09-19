package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_checkDuplicatedTos(t *testing.T) {
	tests := []struct {
		name string
		rs   []replacement
		want duplicatedTosProblems
	}{
		{
			name: "empty",
			rs:   nil,
			want: []duplicatedTosProblem{},
		},
		{
			name: "single",
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
			want: []duplicatedTosProblem{},
		},
		{
			name: "no_culprit_single",
			rs: []replacement{
				{
					from:    "foo.com/example",
					to:      "foo.com/example",
					culprit: nil,
				},
			},
			want: []duplicatedTosProblem{},
		},
		{
			name: "no_culprit_multiple",
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
			want: []duplicatedTosProblem{},
		},
		{
			name: "no_problem",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
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
				{
					from:    "quux.com/example",
					to:      "quux.com/example",
					culprit: nil,
				},
			},
			want: []duplicatedTosProblem{},
		},
		{
			name: "duplicated",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
				{
					from: "baz.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "baz.com",
						Destination: "bar.com",
					},
				},
			},
			want: []duplicatedTosProblem{
				{
					replacements: []replacement{
						{
							from: "foo.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "foo.com",
								Destination: "bar.com",
							},
						},
						{
							from: "baz.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "baz.com",
								Destination: "bar.com",
							},
						},
					},
				},
			},
		},
		{
			name: "duplicated_no_culprit",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
			},
			want: []duplicatedTosProblem{
				{
					replacements: []replacement{
						{
							from: "foo.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "foo.com",
								Destination: "bar.com",
							},
						},
						{
							from:    "bar.com/example",
							to:      "bar.com/example",
							culprit: nil,
						},
					},
				},
			},
		},
		{
			name: "duplicated_exactly",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
			},
			want: []duplicatedTosProblem{
				{
					replacements: []replacement{
						{
							from: "foo.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "foo.com",
								Destination: "bar.com",
							},
						},
						{
							from: "foo.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "foo.com",
								Destination: "bar.com",
							},
						},
					},
				},
			},
		},
		{
			name: "duplicated_exactly_no_culprit",
			rs: []replacement{
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
				{
					from:    "bar.com/example",
					to:      "bar.com/example",
					culprit: nil,
				},
			},
			want: []duplicatedTosProblem{
				{
					replacements: []replacement{
						{
							from:    "bar.com/example",
							to:      "bar.com/example",
							culprit: nil,
						},
						{
							from:    "bar.com/example",
							to:      "bar.com/example",
							culprit: nil,
						},
					},
				},
			},
		},
		{
			name: "multiple_problems",
			rs: []replacement{
				{
					from: "foo.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "foo.com",
						Destination: "bar.com",
					},
				},
				{
					from: "baz.com/example2",
					to:   "bar.com/example2",
					culprit: &DomainMapItem{
						Source:      "baz.com",
						Destination: "bar.com",
					},
				},
				{
					from: "quux.com/example2",
					to:   "bar.com/example2",
					culprit: &DomainMapItem{
						Source:      "quux.com",
						Destination: "bar.com",
					},
				},
				{
					from: "good.com/example",
					to:   "better.com/example",
					culprit: &DomainMapItem{
						Source:      "good.com",
						Destination: "better.com",
					},
				},
				{
					from:    "quux.com/example",
					to:      "quux.com/example",
					culprit: nil,
				},
				{
					from: "nil.com/example",
					to:   "quux.com/example",
					culprit: &DomainMapItem{
						Source:      "nil.com",
						Destination: "quux.com",
					},
				},
				{
					from: "baz.com/example",
					to:   "bar.com/example",
					culprit: &DomainMapItem{
						Source:      "baz.com",
						Destination: "bar.com",
					},
				},
				{
					from: "good.com/example2",
					to:   "better.com/example2",
					culprit: &DomainMapItem{
						Source:      "good.com",
						Destination: "better.com",
					},
				},
			},
			want: []duplicatedTosProblem{
				{
					replacements: []replacement{
						{
							from: "foo.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "foo.com",
								Destination: "bar.com",
							},
						},
						{
							from: "baz.com/example",
							to:   "bar.com/example",
							culprit: &DomainMapItem{
								Source:      "baz.com",
								Destination: "bar.com",
							},
						},
					},
				},
				{
					replacements: []replacement{
						{
							from: "baz.com/example2",
							to:   "bar.com/example2",
							culprit: &DomainMapItem{
								Source:      "baz.com",
								Destination: "bar.com",
							},
						},
						{
							from: "quux.com/example2",
							to:   "bar.com/example2",
							culprit: &DomainMapItem{
								Source:      "quux.com",
								Destination: "bar.com",
							},
						},
					},
				},
				{
					replacements: []replacement{
						{
							from:    "quux.com/example",
							to:      "quux.com/example",
							culprit: nil,
						},
						{
							from: "nil.com/example",
							to:   "quux.com/example",
							culprit: &DomainMapItem{
								Source:      "nil.com",
								Destination: "quux.com",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkDuplicatedTos(tt.rs)

			ss := cmpopts.SortSlices(func(a, b duplicatedTosProblem) bool {
				return a.replacements[0].to < b.replacements[0].to
			})

			au := cmp.AllowUnexported(duplicatedTosProblem{}, replacement{})

			if diff := cmp.Diff(tt.want, got, ss, au); diff != "" {
				t.Errorf("checkDuplicatedTos() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
