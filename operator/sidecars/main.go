package main

import (
	"strings"

	"github.com/zondax/golem/pkg/cli"
	"github.com/zondax/sidecars/internal/commands"
	"github.com/zondax/sidecars/internal/conf"
	"github.com/zondax/sidecars/internal/version"
)

func main() {
	appName := "sidecars"
	envPrefix := strings.ReplaceAll(appName, "-", "_")

	appSettings := cli.AppSettings{
		Name:        appName,
		Description: "",
		ConfigPath:  "$HOME/.sidecars/",
		EnvPrefix:   envPrefix,
		GitVersion:  version.GitVersion,
		GitRevision: version.GitRevision,
	}

	// Define application level features
	cli := cli.New[conf.Config](appSettings)
	defer cli.Close()

	cli.GetRoot().AddCommand(commands.GetStartCommand(cli))

	cli.Run()
}
