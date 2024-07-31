package config

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/flags"
	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/utils/template"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configDir  = "config"
	configFile = "app.toml"
)

// InitializeCommand initializes the command.
func InitializeCommand[
	ConfigT Config[ConfigT],
	NodeT Node[ConfigT],
](
	cmd *cobra.Command,
	appTemplate string,
	logger log.Logger[any],
) error {
	executableName, err := ExecutableName()
	if err != nil {
		return err
	}
	v := NewPrefixedViper(executableName)
	if err = BindFlags(executableName, cmd, v); err != nil {
		return err
	}

	// Build the config.
	config, err := readOrGenerateConfig[ConfigT](v, appTemplate)
	if err != nil {
		return err
	}
	return AttachConfigToCommand(cmd, v, config)
}

// ExecutableName returns the name of the executable.
func ExecutableName() (string, error) {
	executableName, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err, ErrFailedToFetchExecutable.Error())
	}
	return path.Base(executableName), nil
}

// AttachConfigToCommand attaches the config to the command.
func AttachConfigToCommand[ConfigT Config[ConfigT]](
	cmd *cobra.Command, v *viper.Viper, config ConfigT,
) error {
	cmd.SetContext(
		context.WithValue(cmd.Context(), ConfigKey{}, config),
	)

	// Merge the app.toml into the viper instance.
	file := strings.Split(configFile, ".")
	viper.SetConfigName(file[0])
	viper.SetConfigType(file[1])
	viper.AddConfigPath(configDir)
	return v.MergeInConfig()
}

// readOrGenerateConfig reads the config from the app.toml file or
// generates it if it doesn't exist.
func readOrGenerateConfig[ConfigT Config[ConfigT]](
	v *viper.Viper,
	appTemplate string,
) (config ConfigT, err error) {
	rootDir := v.GetString(flags.FlagHome)
	configDirPath := filepath.Join(rootDir, configDir)
	configFilePath := filepath.Join(configDirPath, configFile)
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		return generateConfig[ConfigT](
			v,
			configFilePath,
			appTemplate,
		)
	}

	return config, nil
}

func generateConfig[ConfigT Config[ConfigT]](
	v *viper.Viper,
	configFilePath string,
	appTemplate string,
) (config ConfigT, err error) {
	// set the config template
	if err = template.Set(appTemplate); err != nil {
		return config, err
	}
	// populate appConfig with the values from the viper instance
	if err = v.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, writeConfig(configFilePath, config)
}

func writeConfig[ConfigT Config[ConfigT]](
	configFilePath string,
	config ConfigT,
) (err error) {
	return template.WriteConfigFile(configFilePath, config)
}
