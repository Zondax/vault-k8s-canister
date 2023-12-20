package admctrl_injector

import (
	"context"
	"github.com/zondax/tororu-operator/operator/common"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// getPatchObject generates a PatchObject containing patches based on annotations in the Pod object.
// It takes a context and a Pod object as input and returns the generated PatchObject.
func getPatchObject(ctx context.Context, pod corev1.Pod) *PatchObject {
	// Log the Pod name and its annotations
	zap.S().Infof("[ADM_CTRL] Pod: %s", pod.Name)
	for k, v := range pod.ObjectMeta.Annotations {
		zap.S().Infof("[ADM CTRL] Annotation: %s = %s", k, v)
	}

	var patch PatchObject

	// If the Pod is not managed by Tororu, return an empty patch
	if pod.Annotations[common.TororuManagedAnnotation] == "false" {
		return &patch
	}

	// If a read-write (RW) resource annotation is present, generate sidecar patches
	if val, ok := pod.Annotations[common.TororuResourceReqRWAnnotation]; ok && len(val) > 0 {
		patch = append(patch, *genSidecarPatches(ctx, pod)...)
	}

	// If a read-only (RO) resource annotation is present, generate environment variable patches
	if val, ok := pod.Annotations[common.TororuResourceReqROAnnotation]; ok && len(val) > 0 {
		patch = append(patch, *genEnvPatches(ctx, pod)...)
	}

	// Return the generated PatchObject
	return &patch
}
