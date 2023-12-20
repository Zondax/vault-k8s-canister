package conf

type Config struct {
	Dev *DevConfig `json:"dev"`
}

type DevConfig struct {
	Kubeconfig string `json:"kubeconfig"`
}

func (c Config) SetDefaults() {
}

func (c Config) Validate() error {
	// TODO implement me
	return nil
}
