package main

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	goversion "github.com/caarlos0/go-version"
)

const (
	art = `
###########################################################################
#                                                                         #
#  #    # # #####  #####    ##   #####   ##    ####  #   # #    #  ####   #
#  #    # # #    # #    #  #  #    #    #  #  #       # #  ##   # #    #  #
#  #    # # #    # #    # #    #   #   #    #  ####    #   # #  # #       #
#  #    # # #####  #    # ######   #   ######      #   #   #  # # #       #
#   #  #  # #      #    # #    #   #   #    # #    #   #   #   ## #    #  #
#    ##   # #      #####  #    #   #   #    #  ####    #   #    #  ####   #
#                                                                         #
###########################################################################

`
)

// To be set by ldflags.
var (
	builtBy string //nolint:gochecknoglobals
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
	if builtBy == "" {
		builtBy = "unknown"
	}

	vi := goversion.GetVersionInfo(
		goversion.WithASCIIName(art),
		goversion.WithAppDetails("vipdatasync", description, "https://github.com/typisttech/vipdatasync"),
		goversion.WithBuiltBy(builtBy),
	)

	fmt.Fprintln(writer, vi.String())
}
