package analyzer

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration structure for the analyzer.
type Config struct {
	Files  []string `yaml:"files"`
	Checks struct {
		Security struct {
			OpenSSH          bool `yaml:"open_ssh"`
			PublicS3         bool `yaml:"public_s3"`
			HardcodedSecrets bool `yaml:"hardcoded_secrets"`
			UnencryptedData  bool `yaml:"unencrypted_data"`
			Compliance       struct {
				PCI  bool `yaml:"pci"`
				HIPAA bool `yaml:"hipaa"`
				GDPR  bool `yaml:"gdpr"`
			} `yaml:"compliance"`
		} `yaml:"security"`
		Cost struct {
			OversizedInstances map[string]string `yaml:"oversized_instances"` // Instance type mappings
			UnderutilizedEBS   bool              `yaml:"underutilized_ebs"`
			UnattachedEIP      bool              `yaml:"unattached_eip"`
		} `yaml:"cost"`
	} `yaml:"checks"`
	Output struct {
		Verbosity string `yaml:"verbosity"` // "low", "medium", "high"
	} `yaml:"output"`
}

// LoadConfig loads the configuration from the specified YAML file.
func LoadConfig(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set default values if not specified
	if cfg.Output.Verbosity == "" {
		cfg.Output.Verbosity = "medium"
	}

	return cfg, nil
}
