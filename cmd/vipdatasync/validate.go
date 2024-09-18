package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/typisttech/vipdatasync/internal"
)

type validateCmd struct {
	URLsJSONPath   string `help:"Path to JSON file containing a list of production URLs generated via '$ vip @my-app.production -- wp site list --fields=url --format=json'" name:"urls"   placeholder:"<path>" required:"" type:"existingfile"` //nolint:lll
	ConfigYAMLPath string `help:"Path to environment-specific YAML config file"                                                                                              name:"config" placeholder:"<path>" required:"" type:"existingfile"` //nolint:lll
}

func (c *validateCmd) Run(app *kong.Kong) error {
	urls, err := internal.NewURLsFromJSONFile(c.URLsJSONPath)
	if err != nil {
		return err
	}

	dm, err := internal.NewDomainMapFromConfigFile(c.ConfigYAMLPath)
	if err != nil {
		return err
	}

	v := internal.NewValidator()
	rs := v.Validate(dm, urls)

	count := 0
	for i := range rs {
		count += rs[i].Len()
		fmt.Fprint(app.Stdout, rs[i].Text())
		fmt.Fprint(app.Stdout, "\n")
	}

	if count > 0 {
		return fmt.Errorf("%d problem(s) found", count)
	}

	return nil
}
