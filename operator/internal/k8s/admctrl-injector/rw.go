package admctrl_injector

import (
	"context"
	"fmt"
	common2 "github.com/zondax/vault-k8s-canister/operator/common"
	"github.com/zondax/vault-k8s-canister/operator/common/v1"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// generateSidecarConfigs generates configurations for sidecar containers based on allowed Tororu resources.
// It takes a context and a Pod object as input and returns a slice of SidecarConfig objects.
func generateSidecarConfigs(ctx context.Context, pod corev1.Pod) []*SidecarConfig {
	// Get a list of Tororu resources allowed for read-write (RW) access in the Pod
	tororuResources := getAllowedTResourceNames(ctx, common2.TororuResourceReqRWAnnotation, pod)

	// Initialize a slice to store sidecar configurations
	sidecarConfigs := []*SidecarConfig{}

	// Iterate through the list of allowed Tororu resources
	for _, tResName := range tororuResources {
		// Parse the NamespacedName from the resource name string
		nsName, err := common2.GetNamespacedNameFromNameString(tResName)
		if err != nil {
			zap.S().Error(err)
			continue
		}

		// Retrieve the TororuResource object based on its name and namespace
		tRes, err := common2.GetTResFromName(ctx, nsName.Namespace, nsName.Name)
		if err != nil {
			zap.S().Error(err)
			continue
		}

		// Get the image associated with the TororuResource's kind
		image, err := common2.GetImageFromTResourceKind(tRes.Spec.Kind)
		if err != nil {
			zap.S().Error(err)
			continue
		}

		// Create a SidecarConfig object and add it to the slice
		sidecarConfigs = append(sidecarConfigs,
			&SidecarConfig{
				Image:           image,
				RotationSeconds: tRes.Spec.Rotate,
				Config:          tRes.Spec.Config,
				SecretRefs:      []v1.TororuResource{*tRes},
			},
		)
	}

	// Merge duplicate sidecar configurations based on image
	return mergeSidecarConfigs(sidecarConfigs)
}

// genSidecarPatches generates patches for adding sidecar containers and volumes to a Pod's spec.
// It takes a context and a Pod object as input and returns a pointer to a PatchObject.
func genSidecarPatches(ctx context.Context, pod corev1.Pod) *PatchObject {
	// Initialize an empty patch slice
	patch := []PatchEle{}

	// Generate sidecar configurations for the given Pod
	scConfs := generateSidecarConfigs(ctx, pod)

	// If there are no sidecar configurations, return an empty patch
	if len(scConfs) == 0 {
		return &patch
	}

	// Prepare the service account for the Pod and add a service account patch
	serviceAccount, err := prepServiceAccountForPod(ctx, pod, scConfs)
	if err != nil {
		zap.S().Error(err)
		return &patch
	}

	patch = append(patch, PatchEle{
		Op:    "replace",
		Path:  "/spec/serviceAccountName",
		Value: serviceAccount.Name,
	})

	// Iterate through the sidecar configurations and add patches for volumes, containers, and volume mounts
	for _, c := range scConfs {
		patch = append(
			patch,
			PatchEle{
				Op:    "add",
				Path:  "/spec/volumes/-",
				Value: createPodVolume(),
			},
			PatchEle{
				Op:    "add",
				Path:  "/spec/containers/-",
				Value: c.createSidecarContainer(),
			},
		)

		// Add volume mounts for the sidecar containers to the existing containers
		for idx := range pod.Spec.Containers {
			patch = append(patch,
				PatchEle{
					Op:    "add",
					Path:  fmt.Sprintf("/spec/containers/%d/volumeMounts/-", idx),
					Value: createContainerVolumeMount(),
				},
			)
		}
	}

	return &patch
}

// GetTororuResourceNsNames returns a slice of TororuResource names in the format "namespace/name".
// It takes no additional input parameters and extracts the namespace and name from each TororuResource
// in the SecretRefs field of the SidecarConfig, combining them in the specified format.
func (c *SidecarConfig) GetTororuResourceNsNames() []string {
	// Initialize a slice to store TororuResource names in the "namespace/name" format
	tororuResourceNames := lo.Map(
		c.SecretRefs,
		func(item v1.TororuResource, idx int) string { return common2.GetPodOrCRDId(item.Name, item.Namespace) },
	)

	return tororuResourceNames
}

// GetTororuResourceNames returns a slice of TororuResource names.
// It takes no additional input parameters and extracts the names of each TororuResource
// in the SecretRefs field of the SidecarConfig.
func (c *SidecarConfig) GetTororuResourceNames() []string {
	// Initialize a slice to store TororuResource names
	tororuResourceNames := lo.Map(
		c.SecretRefs,
		func(item v1.TororuResource, idx int) string { return item.Name },
	)

	return tororuResourceNames
}

// createSidecarContainer creates a new CoreV1 container for the sidecar.
// It takes no input parameters and returns a CoreV1 Container configuration.
func (c *SidecarConfig) createSidecarContainer() corev1.Container {
	// Create a volume mount for the container
	volume := createContainerVolumeMount()

	return corev1.Container{
		Name:  common2.SidecarName, // Set the container name to a common value
		Image: c.Image,             // Set the container image from the SidecarConfig
		// Command:         []string{"/bin/sh", "-c", "while true; do echo 'Hello tororu!'; sleep 30; done"},
		ImagePullPolicy: corev1.PullNever,             // Set image pull policy to 'PullAlways'
		VolumeMounts:    []corev1.VolumeMount{volume}, // Mount the volume created earlier
		Env: []corev1.EnvVar{
			{Name: "CONFIG", Value: c.Config},                                                       // Set the CONFIG environment variable
			{Name: "ROTATION_SECONDS", Value: strconv.Itoa(c.RotationSeconds)},                      // Set the ROTATION_SECONDS environment variable
			{Name: "TORORU_RESOURCE_NAMES", Value: strings.Join(c.GetTororuResourceNsNames(), ",")}, // Set TORORU_RESOURCE_NAMES
		},
	}
}

func createContainerVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      common2.PodVolumeName,
		MountPath: common2.PodVolumeMountPath,
	}

}

func createPodVolume() corev1.Volume {
	return corev1.Volume{
		Name: common2.PodVolumeName,
	}
}

// mergeSidecarConfigs merges multiple sidecar configurations based on their images.
// It takes a slice of SidecarConfig objects as input and returns a slice of merged SidecarConfig objects.
func mergeSidecarConfigs(configs []*SidecarConfig) []*SidecarConfig {
	// Create a map to store merged sidecar configurations with image as the key
	configMap := map[string]*SidecarConfig{}

	// Iterate through the provided sidecar configurations
	for _, c := range configs {
		// Check if a configuration with the same image already exists in the map
		if _, ok := configMap[c.Image]; ok {
			// If the configuration exists, append its SecretRefs to the existing one
			configMap[c.Image].SecretRefs = append(configMap[c.Image].SecretRefs, c.SecretRefs...)
		} else {
			// If the configuration does not exist, add it to the map
			configMap[c.Image] = c
		}
	}

	// Create a result slice to hold the merged configurations
	res := []*SidecarConfig{}

	// Iterate through the merged configurations in the map and append them to the result slice
	for _, c := range configMap {
		res = append(res, c)
	}

	// Return the merged sidecar configurations
	return res
}
