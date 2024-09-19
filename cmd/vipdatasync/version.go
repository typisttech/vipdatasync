package main

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/kong"
	goversion "github.com/caarlos0/go-version"
)

//go:embed art.txt
var art string

const (
	description = `CLI utility for WordPress VIP data sync management`
	website     = `https://github.com/typisttech/vipdatasync`
)

// To be set by ldflags.
//
//nolint:gochecknoglobals
var (
	version      string
	fullCommit   string
	commitDate   string
	gitTreeState string
	builtBy      string
)

type versionFlag bool

func (vf versionFlag) BeforeReset(app *kong.Kong) error { //nolint:unparam
	printVersionInfo(app.Stdout)
	app.Exit(0)

	return nil
}

type versionCmd struct{}

func (c *versionCmd) Run(app *kong.Kong) error { //nolint:unparam
	printVersionInfo(app.Stdout)

	return nil
}

func printVersionInfo(writer io.Writer) {
	i := info{goversion.GetVersionInfo(
		goversion.WithASCIIName(art),
		goversion.WithAppDetails("vipdatasync", description, website),
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
	)}

	fmt.Fprintln(writer, i.String())
}

type info struct {
	goversion.Info
}

func (i info) String() string {
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
