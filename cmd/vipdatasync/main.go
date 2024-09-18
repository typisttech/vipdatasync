package main

import (
	"github.com/alecthomas/kong"
)

const (
	description = `CLI utility for WordPress VIP data sync management`
)

type cli struct {
	ValidateCmd validateCmd `cmd:"" help:"Validate the environment-specific YAML config file against production URLs" name:"validate"` //nolint:lll

	VersionCmd  versionCmd  `cmd:""  help:"Show version information about vipdatasync" name:"version"`
	VersionFlag versionFlag `env:"-" help:"Show version information about vipdatasync" name:"version" short:"v"`
}

func run() int {
	ctx := kong.Parse(
		&cli{}, //nolint:exhaustruct
		kong.Description(description),
		kong.DefaultEnvars("VIPDATASYNC"),
		kong.UsageOnError(),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)

	return 0
}

func main() {
	run()
}
