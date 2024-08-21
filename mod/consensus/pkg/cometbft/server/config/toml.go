// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

//go:embed config.toml.tpl
var DefaultConfigTemplate string

var configTemplate *template.Template

func init() {
	var err error

	tmpl := template.New("appConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(DefaultConfigTemplate); err != nil {
		panic(err)
	}
}

// ParseConfig retrieves the default environment configuration for the application.
func ParseConfig(v *viper.Viper) (*Config, error) {
	conf := DefaultConfig()
	err := v.Unmarshal(conf)

	return conf, err
}

// SetConfigTemplate sets the custom app config template for the application.
func SetConfigTemplate(customTemplate string) error {
	var err error

	tmpl := template.New("appConfigFileTemplate")

	if configTemplate, err = tmpl.Parse(customTemplate); err != nil {
		return err
	}

	return nil
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config interface{}) error {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		return err
	}

	if err := os.WriteFile(configFilePath, buffer.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
