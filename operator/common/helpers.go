package common

import (
	"context"
	"fmt"
	icp2 "github.com/zondax/tororu-operator/operator/common/icp"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
)

// Check using canister should pod be allowed access
// Works the same way for RW resource and RO resource

// ListTororuResourcesForPod(ctx, podId, ro/rw) return secretsAllowed
func ListTororuResourcesForPod(ctx context.Context, podId string, permissionType string) ([]string, []string) {
	// Get the icp canister state
	data, err := icp2.GetCanisterStatus()
	if err != nil {
		zap.S().Error(err)
		return []string{}, []string{}
	}

	serviceResources := make([]string, 0)
	for _, consumer := range data.Consumers {
		if permissionType == ResourcePermissionTypeRW && consumer.PermissionType == icp2.PermissionTypeEnumReadAndWrite && consumer.KubeId == podId {
			serviceResources = append(serviceResources, consumer.SecretKubeId)
		}
		if permissionType == ResourcePermissionTypeRO && consumer.PermissionType == icp2.PermissionTypeEnumReadOnly && consumer.KubeId == podId {
			serviceResources = append(serviceResources, consumer.SecretKubeId)
		}

	}

	pendingServiceResources := make([]string, 0)
	for _, consumer := range data.PendingConsumerReqs {
		if permissionType == ResourcePermissionTypeRW && consumer.PermissionType == icp2.PermissionTypeEnumReadAndWrite && consumer.KubeId == podId {
			pendingServiceResources = append(pendingServiceResources, consumer.SecretKubeId)
		}
		if permissionType == ResourcePermissionTypeRO && consumer.PermissionType == icp2.PermissionTypeEnumReadOnly && consumer.KubeId == podId {
			pendingServiceResources = append(pendingServiceResources, consumer.SecretKubeId)
		}
	}

	return serviceResources, pendingServiceResources
}

type ConsumersType struct {
	RO []string
	RW string
}

type CrdStatus struct {
	Status    uint8
	Consumers ConsumersType
	TTL       uint32
}

func GetCRDStatus(crdId string) (*CrdStatus, []*icp2.Secret, error) {
	// use state to build the CRDStatus
	data, err := icp2.GetCanisterStatus()
	if err != nil {
		return nil, nil, err
	}

	var secret *icp2.Secret
	var pendingReqs []*icp2.Secret
	for i, s := range data.Secrets {
		if s.KubeId == crdId {
			secret = &data.Secrets[i]
			break
		}
	}

	for i, s := range data.PendingSecretReqs {
		if s.KubeId == crdId {
			pendingReqs = append(pendingReqs, &data.PendingSecretReqs[i])
		}
	}

	roConsumers := make([]string, 0)
	rwConsumer := ""

	if secret == nil {
		return nil, pendingReqs, nil
	}

	for _, consumer := range data.Consumers {
		if consumer.PermissionType == icp2.PermissionTypeEnumReadAndWrite && consumer.SecretKubeId == secret.KubeId {
			rwConsumer = consumer.KubeId
		}
		if consumer.PermissionType == icp2.PermissionTypeEnumReadOnly && consumer.SecretKubeId == secret.KubeId {
			roConsumers = append(roConsumers, consumer.KubeId)
		}
	}

	return &CrdStatus{
		Status: secret.PendingType,
		Consumers: ConsumersType{
			RO: roConsumers,
			RW: rwConsumer,
		},
		TTL: secret.Ttl,
	}, pendingReqs, nil
}

// Maybe we want more info on crdID on canisterState, TBD
func CreateCRDOnboardRequest(crdId string, ttl uint32) error {
	client, err := icp2.GetBackendClient()
	if err != nil {
		return err
	}

	result, err := client.AddSecret(crdId, ttl)
	if err != nil {
		return err
	}

	if result.Err != nil {
		return fmt.Errorf(*result.Err)
	}

	return nil
}
func createROAccessRequest(_ context.Context, crdIds []string, podId string) error {
	client, err := icp2.GetBackendClient()
	if err != nil {
		return err
	}

	for _, crdId := range crdIds {
		zap.S().Infof("[ADM CTRL] createROAccessRequest - %s %s %d \n", podId, crdId, icp2.PermissionTypeEnumReadOnly)
		result, err := client.AddConsumer(podId, crdId, icp2.PermissionTypeEnumReadOnly)
		if err != nil {
			return err
		}

		if result.Err != nil {
			return fmt.Errorf(*result.Err)
		}
	}

	return nil
}

func UpdateSecret(ctx context.Context, crdId string, ttl uint32) error {
	client, err := icp2.GetBackendClient()
	if err != nil {
		return err
	}

	result, err := client.UpdateSecret(crdId, ttl)
	if err != nil {
		return err
	}

	if result.Err != nil {
		return fmt.Errorf(*result.Err)
	}

	return nil
}

func createRWAccessRequest(_ context.Context, crdIds []string, podId string) error {
	client, err := icp2.GetBackendClient()
	if err != nil {
		return err
	}

	for _, crdId := range crdIds {
		zap.S().Infof("[ADM CTRL] createRWAccessRequest - %s %s %d \n", podId, crdId, icp2.PermissionTypeEnumReadAndWrite)
		result, err := client.AddConsumer(podId, crdId, icp2.PermissionTypeEnumReadAndWrite)
		if err != nil {
			return err
		}

		if result.Err != nil {
			return fmt.Errorf(*result.Err)
		}
	}

	return nil
}

func CreateSecretAccessRequests(ctx context.Context, tResIds []string, podId, pType string) {
	if pType == ResourcePermissionTypeRO {
		err := createROAccessRequest(ctx, tResIds, podId)
		if err != nil {
			zap.S().Errorf("createROAccessRequest, err: %s", err)
		}
		return
	}

	if pType == ResourcePermissionTypeRW {
		err := createRWAccessRequest(ctx, tResIds, podId)
		if err != nil {
			zap.S().Errorf("createRWAccessRequest, err: %s", err)
		}
		return
	}
}

func GetNamespacedNameFromNameString(input string) (*types.NamespacedName, error) {
	// Split the input string by "/"
	parts := strings.Split(input, "/")

	// Ensure that the split resulted in two parts
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid input format")
	}

	// Extract the namespace and name
	namespace := parts[0]
	name := parts[1]

	// Create a metav1.NamespacedName
	nsName := &types.NamespacedName{
		Namespace: strings.TrimSpace(namespace),
		Name:      strings.TrimSpace(name),
	}

	return nsName, nil
}

func GetPodOrCRDId(name, namespace string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}
