package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/spf13/viper"
)

// TODO: annoying from sdk
var configTemplate *template.Template

// ParseConfig retrieves the default configuration for the application.
func ParseConfig[ConfigT interface{ Default() ConfigT }](
	v *viper.Viper,
) (ConfigT, error) {
	var cfg ConfigT
	cfg = cfg.Default()
	err := v.Unmarshal(cfg)

	return cfg, err
}

// Set sets the custom app config template for the application.
func Set(customTemplate string) error {
	tmpl := template.New("appConfigFileTemplate")

	var err error
	if configTemplate, err = tmpl.Parse(customTemplate); err != nil {
		return err
	}

	return nil
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config any) error {
	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, config); err != nil {
		return err
	}

	if err := os.WriteFile(configFilePath, buffer.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
