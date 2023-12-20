package admctrl_injector

import (
	"context"
	"fmt"
	common2 "github.com/zondax/tororu-operator/operator/common"
	"regexp"

	"github.com/samber/lo"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func extractResourceName(withBraces string) string {
	// Define the regex pattern to capture the name between < and >
	pattern := `<(.+?)>`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Find the first match in the input string
	match := regex.FindStringSubmatch(withBraces)

	if len(match) == 2 {
		// The captured name is in match[1]
		name := match[1]
		return name
	} else {
		return ""
	}
}

type rMap struct {
	EnvVar         string
	Value          string
	ValueFound     bool
	ContainerIndex int
	EnvVarIndex    int
}

// genEnvPatches generates environment variable patches for the given Pod based on allowed TororuResources.
func genEnvPatches(ctx context.Context, pod corev1.Pod) *PatchObject {
	logger := zap.S().With(pod.Namespace, pod.Name)

	// Get the names of allowed TororuResources for Read-Only (RO) access.
	tororuAllowedResources := getAllowedTResourceNames(ctx, common2.TororuResourceReqROAnnotation, pod)
	tororuAllowedResourcesMap := lo.Associate(tororuAllowedResources, func(str string) (string, int) {
		return str, len(str)
	})

	patch := []PatchEle{}
	envVarReverseMap := map[string]rMap{}

	// Iterate through the containers and their environment variables to build a reverse mapping.
	for idx1, c := range pod.Spec.Containers {
		for idx2, e := range c.Env {
			resName := extractResourceName(e.Value)
			if resName != "" {
				envVarReverseMap[resName] = rMap{EnvVar: e.Name, ContainerIndex: idx1, EnvVarIndex: idx2}
			}
		}
	}

	logger.Infof("[ADM CTRL] Env reverse map: %v", envVarReverseMap)
	logger.Infof("[ADM CTRL] Tororu Allowed RO resource: %v", tororuAllowedResources)

	// Iterate through the TororuResources and update the environment variables with secret values.
	for tResNsName := range envVarReverseMap {
		if _, ok := tororuAllowedResourcesMap[tResNsName]; !ok {
			zap.S().Debugf("[ADM CTRL] %s not allowed to read %s", pod.Name, tResNsName)
			continue
		}

		nsName, err := common2.GetNamespacedNameFromNameString(tResNsName)
		if err != nil {
			zap.S().Error(err)
			return &patch
		}

		_, k8Client := common2.GetKubernetesClients()

		k8Secret, err := k8Client.CoreV1().Secrets(nsName.Namespace).Get(context.Background(), nsName.Name+"-secret", metav1.GetOptions{})
		if err != nil {
			zap.S().Error(err)
			return &patch
		}

		temp := envVarReverseMap[tResNsName]
		temp.Value = string(k8Secret.Data["value"])
		if temp.Value != "" {
			temp.ValueFound = true
		}
		envVarReverseMap[tResNsName] = temp
	}

	// Create patches to update environment variables with secret values.
	for _, v := range envVarReverseMap {
		if !v.ValueFound {
			continue
		}

		patch = append(
			patch,
			PatchEle{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%d/env/%d", v.ContainerIndex, v.EnvVarIndex),
				Value: createEnvVar(v.EnvVar, v.Value),
			},
		)
	}

	// Add annotation and patch if environment variables were updated.
	if len(patch) > 0 {
		pod.ObjectMeta.Annotations[common2.TororuEnvVarsAppliedAnnotation] = common2.TororuAnnotationsTrueString
		patch = append(
			patch,
			PatchEle{
				Op:    "add",
				Path:  fmt.Sprintf("/metadata/annotations/%s", common2.TororuEnvVarsAppliedAnnotation),
				Value: "true",
			},
		)
	}

	return &patch
}

// createEnvVar creates an environment variable with the given name and value.
func createEnvVar(name, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  name,
		Value: value,
	}
}
