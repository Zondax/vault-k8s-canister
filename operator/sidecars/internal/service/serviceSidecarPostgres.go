package service

import (
	"github.com/zondax/golem/pkg/metrics"
	"github.com/zondax/golem/pkg/runner"
	"github.com/zondax/sidecars/internal/conf"
	sidecarPostgres "github.com/zondax/sidecars/internal/sidecar-postgres"
)

func StartSidecarPostgres(config *conf.Config) {
	r := runner.NewRunner()

	r.AddTask(metrics.NewTaskMetrics("/metrics", "9090"))

	r.AddTask(sidecarPostgres.NewSidecarPostgres(config))

	r.StartAndWait()
}
