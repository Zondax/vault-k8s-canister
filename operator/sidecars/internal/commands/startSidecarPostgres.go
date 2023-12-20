package commands

import (
	"github.com/spf13/cobra"
	"github.com/zondax/golem/pkg/cli"
	"github.com/zondax/sidecars/internal/conf"
	"github.com/zondax/sidecars/internal/service"
	"go.uber.org/zap"
)

func GetStartCommand(c *cli.CLI) *cobra.Command {
	return &cobra.Command{
		Use:   "start-sidecar-postgres",
		Short: "Start",
		Run: func(cmd *cobra.Command, args []string) {
			start(c, cmd, args)
		},
	}
}

func start(c *cli.CLI, _ *cobra.Command, _ []string) {
	zap.S().Infof(c.GetVersionString())

	config, err := cli.LoadConfig[conf.Config]()
	if err != nil {
		zap.S().Errorf("Error loading config: %s", err)
		return
	}

	service.StartSidecarPostgres(config)
}
