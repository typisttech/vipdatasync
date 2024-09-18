package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var normalURLs = URLs{ //nolint:gochecknoglobals
	"https://www.example.hk/zh",
	"https://www.example.hk/en",
	"https://www.example.hk",
	"https://www.example.com.uk",
	"https://vip-example.com/de",
	"https://vip-example.com",
	"https://example.go-vip.net/fr",
	"https://example.go-vip.net",
}

func TestNewURLsFromJSONFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    URLs
		wantErr bool
	}{
		{
			name:    "not_exist",
			path:    "testdata/site/not/exist/urls.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			path:    "testdata/site/empty.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not_json",
			path:    "testdata/site/not-json.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no_urls",
			path:    "testdata/site/no-urls.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "some_urls",
			path:    "testdata/site/some-urls.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "normal",
			path:    "testdata/site/normal.json",
			want:    normalURLs,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewURLsFromJSONFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NewURLsFromJSONFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewURLsFromJSONFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
