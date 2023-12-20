package service

import (
	"github.com/zondax/golem/pkg/metrics"
	"github.com/zondax/golem/pkg/runner"
	"github.com/zondax/vault-k8s-canister/operator/internal/conf"
	"github.com/zondax/vault-k8s-canister/operator/internal/k8s/admctrl-injector"
	operatorCRD "github.com/zondax/vault-k8s-canister/operator/internal/k8s/op-crd"
	operatorSidecar "github.com/zondax/vault-k8s-canister/operator/internal/k8s/op-sidecar"
)

func Start(config *conf.Config) {
	r := runner.NewRunner()

	r.AddTask(metrics.NewTaskMetrics("/metrics", "9090"))

	r.AddTask(operatorSidecar.NewSidecarOperator(config))
	r.AddTask(operatorCRD.NewCRDOperator(config))
	r.AddTask(admctrl_injector.NewAdmCtrlInjector(config))

	r.StartAndWait()
}
