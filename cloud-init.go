package cloudinitgo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// User represents a system user configuration
type user struct {
	Name       string `yaml:"name"`
	Sudo       string `yaml:"sudo,omitempty"`
	Shell      string `yaml:"shell,omitempty"`
	Password   string `yaml:"passwd,omitempty"`
	LockPasswd bool   `yaml:"lock_passwd,omitempty"`
}

// CloudConfig represents the minimal cloud-init configuration
type cloudConfig struct {
	Users     []user `yaml:"users,omitempty"`
	SSHPwauth bool   `yaml:"ssh_pwauth,omitempty"`
}

func (c *cloudInit) CreateConfigFiles() error {

	// check if dataDir exists, if not create it
	if _, err := os.Stat(c.Config.DataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.Config.DataDir, 0755); err != nil {
			return fmt.Errorf("error creating data directory: %v", err)
		}
	}

	// Create cloud config
	config := NewCloudConfig(c.Config)

	// Configurar marshaller para mantener indentación consistente
	yamlEncoder := yaml.NewEncoder(&bytes.Buffer{})
	yamlEncoder.SetIndent(2)

	// Write user-data file
	err := c.writeYAMLFile("user-data", config)
	if err != nil {
		return fmt.Errorf("error writing user-data: %v", err)
	}

	// Write meta-data file
	meta := map[string]string{
		"instance-id":    "i-local",
		"local-hostname": "cloudinit",
	}
	err = c.writeYAMLFile("meta-data", meta)
	if err != nil {
		return fmt.Errorf("error writing meta-data: %v", err)
	}

	return nil
}

func NewCloudConfig(c Config) cloudConfig {
	return cloudConfig{
		Users: []user{
			{
				Name:       c.UserName,
				Sudo:       "ALL=(ALL) NOPASSWD:ALL",
				Shell:      "/bin/bash",
				Password:   c.Password,
				LockPasswd: false,
			},
		},
		SSHPwauth: true,
	}
}

// marshalYAML marshals data to YAML with consistent indentation
func (s *cloudInit) MarshalYAML(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling YAML: %v", err)
	}
	return buf.Bytes(), nil
}

// writeYAMLFile marshals data to YAML and writes to file
func (s *cloudInit) writeYAMLFile(filename string, data interface{}) error {
	yamlData, err := s.MarshalYAML(data)
	if err != nil {
		return err
	}

	// Write to file
	filePath := filepath.Join(s.Config.DataDir, filename)
	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
