package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	normalDomainMap = DomainMap{ //nolint:gochecknoglobals
		{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
		{"www.example.hk/en", "example-com-staging.go-vip.net/hk-en"},
		{"www.example.hk", "example-com-staging.go-vip.net/hk"},
		{"www.example.com.uk", "example-com-staging.go-vip.net/uk"},
		{"vip-example.com", "example-com-staging.go-vip.net"},
		{"example.go-vip.net", "example-com-staging.go-vip.net"},
	}

	reversedDomainMap = DomainMap{ //nolint:gochecknoglobals
		{"example.go-vip.net", "example-com-staging.go-vip.net"},
		{"vip-example.com", "example-com-staging.go-vip.net"},
		{"www.example.com.uk", "example-com-staging.go-vip.net/uk"},
		{"www.example.hk", "example-com-staging.go-vip.net/hk"},
		{"www.example.hk/en", "example-com-staging.go-vip.net/hk-en"},
		{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
	}
)

func TestNewDomainMapFromConfigFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    DomainMap
		wantErr bool
	}{
		{
			name:    "not_exist",
			path:    "testdata/config/not/exist/config.yml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			path:    "testdata/config/empty.yml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty_map",
			path:    "testdata/config/empty-map.yml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not_yaml",
			path:    "testdata/config/not-yaml.yml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "reversed",
			path:    "testdata/config/reversed.yml",
			want:    reversedDomainMap,
			wantErr: false,
		},
		{
			name:    "normal",
			path:    "testdata/config/normal.yml",
			want:    normalDomainMap,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewDomainMapFromConfigFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NewDomainMapFromConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewDomainMapFromConfigFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDomainMap_replace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dm   DomainMap
		urls URLs
		want []replacement
	}{
		{
			name: "single",
			dm: DomainMap{
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
			},
			urls: URLs{"https://www.example.hk/zh"},
			want: []replacement{
				{
					from:    "https://www.example.hk/zh",
					to:      "https://example-com-staging.go-vip.net/hk-zh",
					culprit: &DomainMapItem{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
				},
			},
		},
		{
			name: "extraneous_map",
			dm: DomainMap{
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
				{"www.example.hk", "example-com-staging.go-vip.net/hk"},
			},
			urls: URLs{"https://www.example.hk/zh"},
			want: []replacement{
				{
					from:    "https://www.example.hk/zh",
					to:      "https://example-com-staging.go-vip.net/hk-zh",
					culprit: &DomainMapItem{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
				},
			},
		},
		{
			name: "order_matters_exact",
			dm: DomainMap{
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh-1"},
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh-2"},
			},
			urls: URLs{"https://www.example.hk/zh"},
			want: []replacement{
				{
					from:    "https://www.example.hk/zh",
					to:      "https://example-com-staging.go-vip.net/hk-zh-1",
					culprit: &DomainMapItem{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh-1"},
				},
			},
		},
		{
			name: "order_matters_subset",
			dm: DomainMap{
				{"www.example.hk", "example-com-staging.go-vip.net/hk"},
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
			},
			urls: URLs{"https://www.example.hk/zh"},
			want: []replacement{
				{
					from:    "https://www.example.hk/zh",
					to:      "https://example-com-staging.go-vip.net/hk/zh",
					culprit: &DomainMapItem{"www.example.hk", "example-com-staging.go-vip.net/hk"},
				},
			},
		},
		{
			name: "subset",
			dm: DomainMap{
				{"www.example.hk", "example-com-staging.go-vip.net/hk"},
			},
			urls: URLs{"https://www.example.hk/en"},
			want: []replacement{
				{
					from:    "https://www.example.hk/en",
					to:      "https://example-com-staging.go-vip.net/hk/en",
					culprit: &DomainMapItem{"www.example.hk", "example-com-staging.go-vip.net/hk"},
				},
			},
		},
		{
			name: "no_match",
			dm: DomainMap{
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
			},
			urls: URLs{"https://www.example.com.uk"},
			want: []replacement{
				{
					from:    "https://www.example.com.uk",
					to:      "https://www.example.com.uk",
					culprit: nil,
				},
			},
		},
		{
			name: "no_match_multiple",
			dm: DomainMap{
				{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
			},
			urls: URLs{"https://vip-example.com/de", "https://vip-example.com/fr"},
			want: []replacement{
				{
					from:    "https://vip-example.com/de",
					to:      "https://vip-example.com/de",
					culprit: nil,
				},
				{
					from:    "https://vip-example.com/fr",
					to:      "https://vip-example.com/fr",
					culprit: nil,
				},
			},
		},
		{
			name: "same_map",
			dm: DomainMap{
				{"vip-example.com", "example-com-staging.go-vip.net"},
			},
			urls: URLs{"https://vip-example.com/de", "https://vip-example.com/fr"},
			want: []replacement{
				{
					from:    "https://vip-example.com/de",
					to:      "https://example-com-staging.go-vip.net/de",
					culprit: &DomainMapItem{"vip-example.com", "example-com-staging.go-vip.net"},
				},
				{
					from:    "https://vip-example.com/fr",
					to:      "https://example-com-staging.go-vip.net/fr",
					culprit: &DomainMapItem{"vip-example.com", "example-com-staging.go-vip.net"},
				},
			},
		},
		{
			name: "normal",
			dm:   normalDomainMap,
			urls: normalURLs,
			want: []replacement{
				{
					from:    "https://www.example.hk/zh",
					to:      "https://example-com-staging.go-vip.net/hk-zh",
					culprit: &DomainMapItem{"www.example.hk/zh", "example-com-staging.go-vip.net/hk-zh"},
				},
				{
					from:    "https://www.example.hk/en",
					to:      "https://example-com-staging.go-vip.net/hk-en",
					culprit: &DomainMapItem{"www.example.hk/en", "example-com-staging.go-vip.net/hk-en"},
				},
				{
					from:    "https://www.example.hk",
					to:      "https://example-com-staging.go-vip.net/hk",
					culprit: &DomainMapItem{"www.example.hk", "example-com-staging.go-vip.net/hk"},
				},
				{
					from:    "https://www.example.com.uk",
					to:      "https://example-com-staging.go-vip.net/uk",
					culprit: &DomainMapItem{"www.example.com.uk", "example-com-staging.go-vip.net/uk"},
				},
				{
					from:    "https://vip-example.com/de",
					to:      "https://example-com-staging.go-vip.net/de",
					culprit: &DomainMapItem{"vip-example.com", "example-com-staging.go-vip.net"},
				},
				{
					from:    "https://vip-example.com",
					to:      "https://example-com-staging.go-vip.net",
					culprit: &DomainMapItem{"vip-example.com", "example-com-staging.go-vip.net"},
				},
				{
					from:    "https://example.go-vip.net/fr",
					to:      "https://example-com-staging.go-vip.net/fr",
					culprit: &DomainMapItem{"example.go-vip.net", "example-com-staging.go-vip.net"},
				},
				{
					from:    "https://example.go-vip.net",
					to:      "https://example-com-staging.go-vip.net",
					culprit: &DomainMapItem{"example.go-vip.net", "example-com-staging.go-vip.net"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.dm.replace(tt.urls)

			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(replacement{})); diff != "" {
				t.Errorf("DomainMap.replace() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
