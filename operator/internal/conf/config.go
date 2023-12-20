package conf

import (
	"fmt"
)

type Config struct {
	Dev *DevConfig `json:"dev"`
}

type AdmControllerConfig struct {
	Port string `json:"port"`
}

type SidecarOp struct {
	Port string `json:"port"`
}

type CRDOp struct {
	Port string `json:"port"`
}

type DevConfig struct {
	Kubeconfig    string              `json:"kubeconfig"`
	AdmController AdmControllerConfig `json:"admController"`
	SidecarOp     SidecarOp           `json:"sidecarOperator"`
	CRDOp         CRDOp               `json:"crdOperator"`
}

func (c Config) SetDefaults() {
}

func (c Config) Validate() error {
	if c.Dev != nil {
		if c.Dev.Kubeconfig == "" {
			return fmt.Errorf("kubeconfig must be set")
		}

		if c.Dev.AdmController.Port == "" {
			return fmt.Errorf("adm controller port must be set")
		}
	}

	// TODO implement me
	return nil
}
