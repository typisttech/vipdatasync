package internal

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
)

type URLs []string

func NewURLsFromJSONFile(path string) (URLs, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var j []struct {
		URL string `json:"url"`
	}

	if err := json.Unmarshal(content, &j); err != nil {
		return nil, err
	}

	if len(j) == 0 {
		return nil, errors.New("no URLs found")
	}

	ss := make(URLs, len(j))
	for i, s := range j {
		ss[i] = s.URL
	}

	if slices.Contains(ss, "") {
		return nil, errors.New("empty URLs found")
	}

	return ss, nil
}
