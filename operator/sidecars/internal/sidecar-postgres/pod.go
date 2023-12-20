package sidecarPostgres

// TODO: clean if not found any use
// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/zondax/tororu-operator/common"
// 	"go.uber.org/zap"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// const (
// 	// Specify the path to the namespace file
// 	namespaceFilePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
// )

// func getNamespaceFromInsidePod() (string, error) {
// 	// Read the content of the file
// 	content, err := os.ReadFile(namespaceFilePath)
// 	if err != nil {
// 		fmt.Printf("Error reading file: %v\n", err)
// 		return "", err
// 	}

// 	// Convert the content to a string
// 	namespace := string(content)

// 	return namespace, nil
// }

// func getOwnPod() (*corev1.Pod, error) {
// 	namespace, err := getNamespaceFromInsidePod()
// 	if err != nil {
// 		zap.S().Errorf("Error getting namespace: %v", err)
// 		return nil, err
// 	}

// 	podName, err := os.Hostname()
// 	if err != nil {
// 		zap.S().Errorf("Error getting hostname: %v", err)
// 		return nil, err
// 	}

// 	_, kClient := common.GetKubernetesClients()

// 	// Get the pod by name in the specified namespace.
// 	pod, err := kClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
// 	return pod, err
// }
