package sidecarPostgres

import (
	"github.com/zondax/sidecars/internal/conf"
	"github.com/zondax/tororu-operator/operator/common"
	"go.uber.org/zap"
)

type SidecarPostgres struct {
	name   string
	config *conf.Config
}

func NewSidecarPostgres(config *conf.Config) *SidecarPostgres {
	// kubeconfigPath := ""
	// if config.Dev != nil {
	// 	kubeconfigPath = config.Dev.Kubeconfig
	// }
	common.CreateCommonKubernetesClient("~/.kube/config")
	return &SidecarPostgres{
		name:   "sidecars",
		config: config,
	}
}

func (a SidecarPostgres) Name() string {
	return a.name
}

func (a SidecarPostgres) Start() error {
	// Get the in-cluster configuration

	resourceNames := getTResNames()
	commsChan := make(chan string)

	dynamicClient, kubeClient := common.GetKubernetesClients()

	for _, res := range resourceNames {
		zap.S().Infof("Resource found for sidecar: %s", res)
		info := &tResInfo{
			commChan:         commsChan,
			Name:             res,
			RotationDuration: getRotationDuration(),
			dynamicClient:    dynamicClient,
			kubeClient:       kubeClient,
		}
		go info.rotateAndUpdateForever()
	}

	for {
		msg := <-commsChan
		zap.S().Info(msg)
	}
}

func (a SidecarPostgres) Stop() error {
	// TODO implement me
	return nil
}
