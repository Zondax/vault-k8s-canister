package admctrl_injector

import (
	"context"
	"fmt"
	common2 "github.com/zondax/vault-k8s-canister/operator/common"
	"strings"

	"github.com/samber/lo"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// getPermissionTypeFromAnnotation maps an annotation to its corresponding permission type.
// It returns the permission type or an error if the annotation is not recognized.
func getPermissionTypeFromAnnotation(annotation string) (string, error) {
	switch annotation {
	case common2.TororuResourceReqROAnnotation:
		return common2.ResourcePermissionTypeRO, nil
	case common2.TororuResourceReqRWAnnotation:
		return common2.ResourcePermissionTypeRW, nil
	default:
		return "", fmt.Errorf("not found")
	}
}

// getAllowedTResourceNames retrieves the list of allowed TororuResources based on the provided annotation.
// It parses the annotation in the Pod's annotations, retrieves requested resources, and validates them against
// allowed resources. It also creates SecretAccessRequests for resources that are not allowed.
// The resulting list contains the intersection of requested and allowed resources.
func getAllowedTResourceNames(ctx context.Context, annotation string, pod corev1.Pod) []string {
	secretsString := pod.Annotations[annotation]
	secretsToRequest := strings.Split(secretsString, ",")
	for idx, s := range secretsToRequest {
		secretsToRequest[idx] = strings.TrimSpace(s)
	}

	pType, err := getPermissionTypeFromAnnotation(annotation)
	if err != nil {
		zap.S().Error(err)
		return []string{}
	}

	podId := common2.GetPodOrCRDId(pod.Name, pod.Namespace)
	secretsAllowed, secretsToBeAllowed := common2.ListTororuResourcesForPod(ctx, podId, pType)
	intersect := lo.Intersect[string](secretsToRequest, secretsAllowed)

	secretsNotAllowed, _ := lo.Difference(secretsToRequest, secretsAllowed)
	secretsNotAllowed, _ = lo.Difference(secretsNotAllowed, secretsToBeAllowed)

	zap.S().Infof("[getAllowedTResourceNames] secretsToRequest: %+v", secretsToRequest)
	zap.S().Infof("[getAllowedTResourceNames] secretsToBeAllowed: %+v", secretsToBeAllowed)
	zap.S().Infof("[getAllowedTResourceNames] secretsAllowed: %+v", secretsAllowed)
	zap.S().Infof("[getAllowedTResourceNames] secretsNotAllowed: %v", secretsNotAllowed)
	zap.S().Infof("[getAllowedTResourceNames] intersect: %v", intersect)

	common2.CreateSecretAccessRequests(ctx, secretsNotAllowed, podId, pType)

	if len(intersect) == len(secretsToRequest) {
		return intersect
	}

	return []string{}
}
