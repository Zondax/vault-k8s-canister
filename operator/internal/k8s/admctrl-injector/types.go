package admctrl_injector

import (
	"github.com/zondax/tororu-operator/operator/common/v1"
)

type PatchObject = []PatchEle

type PatchEle struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type SidecarConfig struct {
	// PermissionType  string

	Image           string
	RotationSeconds int
	Config          string

	SecretRefs []v1.TororuResource
}
