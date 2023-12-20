package common

import (
	"context"
	"fmt"
	v12 "github.com/zondax/vault-k8s-canister/operator/common/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetImageFromTResourceKind(kind string) (string, error) {
	fmt.Println("Kind: ", kind)
	// switch kind {
	// case "tororu.zondax.io/postgres":
	// 	return "zondax/sidecar-postgres:latest", nil
	// }

	if kind == "tororu.zondax.io/postgres" {
		return "zondax/sidecar-postgres:latest", nil
	}

	return "", fmt.Errorf("can't find image for the kind")
}

func GetTResFromName(ctx context.Context, namespace, name string) (*v12.TororuResource, error) {
	// resAbsPath := GetCrdAbsPath(namespace, name)
	tResRaw, err := dynamicClient.
		Resource(v12.SchemeGroupVersion.WithResource(TororuResourceNamePlural)).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get tororu resource %s: %s", name, err)
	}

	return v12.FromUnstructured(tResRaw)
}
