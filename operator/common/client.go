package common

import (
	"fmt"
	"github.com/zondax/tororu-operator/operator/common/v1"

	"github.com/zondax/golem/pkg/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var dynamicClient *dynamic.DynamicClient
var kubernetesClient *kubernetes.Clientset

func CreateCommonKubernetesClient(kubeConfigPath string) {
	var config *rest.Config

	// Get the in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("[Common k8s client] Error creating in-cluster config: %v\n", err)
		zap.S().Info("[Common k8s client] will try use kubeconfig to create common client")
		kubeConfigPath, err := utils.ExpandPath(kubeConfigPath)
		if err != nil {
			zap.S().Fatalf("utils.ExpandPath, err: %s", err.Error())
		}

		kConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			zap.S().Fatalf("clientcmd.BuildConfigFromFlags, err: %s", err)
		}

		config = kConfig
	}

	err = v1.AddToScheme(scheme.Scheme)
	if err != nil {
		zap.S().Fatalf("v1.AddToScheme, err: %s", err)
	}
	// Create a Kubernetes client object using the configuration object.
	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		zap.S().Fatalf("dynamic.NewForConfig, err: %s", err)
	}

	kubernetesClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		zap.S().Fatalf("kubernetes.NewForConfig, err: %s", err)
	}

}

func GetKubernetesClients() (*dynamic.DynamicClient, *kubernetes.Clientset) {
	if dynamicClient == nil {
		zap.S().Fatal("Common kubernetes client is nil")
	}

	return dynamicClient, kubernetesClient
}
