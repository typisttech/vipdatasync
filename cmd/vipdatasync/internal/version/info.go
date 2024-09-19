package version

import (
	"fmt"
	"strings"
	"text/tabwriter"

	goversion "github.com/caarlos0/go-version"
)

type Info struct {
	goversion.Info
}

func NewInfo(name, description, website, art, version, fullCommit, commitDate, gitTreeState, builtBy string) Info {
	return Info{
		goversion.GetVersionInfo(
			goversion.WithASCIIName(art),
			goversion.WithAppDetails(name, description, website),
			func(i *goversion.Info) {
				if version != "" {
					i.GitVersion = version
				}

				if fullCommit != "" {
					i.GitCommit = fullCommit
				}

				if commitDate != "" {
					i.BuildDate = commitDate
				}

				if gitTreeState != "" {
					i.GitTreeState = gitTreeState
				}

				if builtBy == "" {
					builtBy = "unknown"
				}
				i.BuiltBy = builtBy
			},
		),
	}
}

func (i Info) String() string {
	b := strings.Builder{}
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0) //nolint:mnd

	// name and description are optional.
	if i.Name != "" {
		if i.ASCIIName != "" {
			fmt.Fprint(w, i.ASCIIName)
			fmt.Fprint(w, "\n")
		}

		fmt.Fprint(w, i.Name)

		if i.Description != "" {
			fmt.Fprintf(w, ": %s", i.Description)
		}

		if i.URL != "" {
			fmt.Fprintf(w, "\n%s", i.URL)
		}

		fmt.Fprint(w, "\n\n")
	}

	fmt.Fprintf(w, "Version:\t%s\n", i.GitVersion)
	fmt.Fprintf(w, "Commit:\t%s\n", i.GitCommit)
	fmt.Fprintf(w, "Commit Date:\t%s\n", i.BuildDate)
	fmt.Fprintf(w, "Git Tree State:\t%s\n", i.GitTreeState)
	fmt.Fprintf(w, "Built By:\t%s\n", i.BuiltBy)
	fmt.Fprintf(w, "Go Version:\t%s\n", i.GoVersion)
	fmt.Fprintf(w, "Compiler:\t%s\n", i.Compiler)
	fmt.Fprintf(w, "Module Sum:\t%s\n", i.ModuleSum)
	fmt.Fprintf(w, "Platform:\t%s\n", i.Platform)

	w.Flush()

	return b.String()
}
