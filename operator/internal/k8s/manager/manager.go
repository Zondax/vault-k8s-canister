package manager

import (
	"github.com/zondax/golem/pkg/utils"
	"github.com/zondax/tororu-operator/operator/internal/conf"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// getClientConfig returns a kubernetes client config
func getClientConfig(config *conf.Config) (*rest.Config, error) {
	if config.Dev == nil {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			zap.S().Fatal(err.Error())
		}
		return cfg, err
	}

	zap.S().Info("[Manager] Detected development mode")
	kubeConfigPath, err := utils.ExpandPath(config.Dev.Kubeconfig)
	if err != nil {
		zap.S().Fatal(err.Error())
	}

	zap.S().Info("[Manager] kubeConfig: ", kubeConfigPath)

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		zap.S().Fatalf("[Manager] Failed to load Kubernetes configuration: %v", err)
	}

	return cfg, err
}

func NewManager(config *conf.Config, operatorName, metricsBindAddress string) (manager.Manager, error) {
	var cfg *rest.Config
	var err error

	cfg, err = getClientConfig(config)
	if err != nil {
		return nil, err
	}

	// Create a new manager to communicate with the Kubernetes API
	zap.S().Infof("[Manager] Creating for %s on %s", operatorName, metricsBindAddress)
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: metricsBindAddress})
	if err != nil {
		zap.S().Fatalf("[Manager] Failed to create new manager: %v", err)
	}

	return mgr, nil
}
