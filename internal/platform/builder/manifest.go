package builder

import (
	"encoding/json"
	"fmt"
)

// Manifest is the portable build description used by Web and Android builders.
type Manifest struct {
	Name        string            `json:"name"`
	Template    string            `json:"template"`
	Backend     string            `json:"backend,omitempty"`
	Database    string            `json:"database,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Features    []string          `json:"features,omitempty"`
}

func (m Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Template == "" {
		return fmt.Errorf("template is required")
	}
	if _, err := GetTemplate(m.Template); err != nil {
		return err
	}
	return nil
}

func EncodeManifest(m Manifest) ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(m, "", "  ")
}
