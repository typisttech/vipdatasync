package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type DomainMap []DomainMapItem

type DomainMapItem struct {
	Source      string
	Destination string
}

func NewDomainMapFromConfigFile(path string) (DomainMap, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var y struct {
		DataSync struct {
			DomainMap DomainMap `yaml:"domain_map"`
		} `yaml:"data_sync"`
	}

	if err := yaml.Unmarshal(content, &y); err != nil {
		return nil, err
	}

	if len(y.DataSync.DomainMap) == 0 {
		return nil, errors.New("domain_map is empty")
	}

	return y.DataSync.DomainMap, nil
}

func (dm *DomainMap) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("domain_map must contain YAML mapping, has %v", value.Kind)
	}

	*dm = make([]DomainMapItem, len(value.Content)/2) //nolint:mnd
	for i := 0; i < len(value.Content); i += 2 {
		res := &(*dm)[i/2]
		if err := value.Content[i].Decode(&res.Source); err != nil {
			return err
		}

		if err := value.Content[i+1].Decode(&res.Destination); err != nil {
			return err
		}
	}

	return nil
}

type replacement struct {
	// from is the URL before the replacement. Expects to be a production URL.
	from string

	// to is the URL after the replacement. Same as from if no replacement was made.
	to string

	// culprit is the pointer to DomainMapItem that caused the replacement. nil if no replacement was made.
	culprit *DomainMapItem
}

func (dm *DomainMap) replace(sites URLs) []replacement {
	rs := make([]replacement, len(sites))

	for i, site := range sites {
		rs[i] = replacement{
			from:    site,
			to:      site,
			culprit: nil,
		}

		for _, dmi := range *dm {
			if strings.Contains(site, dmi.Source) {
				rs[i].to = strings.Replace(site, dmi.Source, dmi.Destination, 1)
				rs[i].culprit = &dmi

				break
			}
		}
	}

	return rs
}
