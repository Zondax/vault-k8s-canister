package main

import (
	"github.com/zondax/golem/pkg/cli"
	"github.com/zondax/vault-k8s-canister/operator/internal/commands"
	"github.com/zondax/vault-k8s-canister/operator/internal/conf"
	"github.com/zondax/vault-k8s-canister/operator/internal/version"
)

func main() {
	appSettings := cli.AppSettings{
		Name:        "tororu-operator",
		Description: "Please override",
		ConfigPath:  "$HOME/.tororu-operator/",
		EnvPrefix:   "tororu-operator",
		GitVersion:  version.GitVersion,
		GitRevision: version.GitRevision,
	}

	// Define application level features
	cli := cli.New[conf.Config](appSettings)
	defer cli.Close()

	cli.GetRoot().AddCommand(commands.GetStartCommand(cli))

	cli.Run()
}
