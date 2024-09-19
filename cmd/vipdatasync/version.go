package main

import (
	_ "embed"
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	iversion "github.com/typisttech/vipdatasync/cmd/vipdatasync/internal/version"
)

//go:embed art.txt
var art string

const (
	description = `CLI utility for WordPress VIP data sync management`
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
	i := iversion.NewInfo(
		`vipdatasync`,
		description,
		`https://github.com/typisttech/vipdatasync`,
		art,
		version,
		fullCommit,
		commitDate,
		gitTreeState,
		builtBy,
	)
	fmt.Fprintln(writer, i.String())
}
