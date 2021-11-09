package cmd

import (
	"fmt"
	cosmcfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	"path/filepath"
)

func readCosmosConfig(homeDir string) (cosmcfg.Config, error) {
	v := viper.New()

	v.SetConfigType("toml")
	v.SetConfigName("app")
	v.AddConfigPath(filepath.Join(homeDir, "config"))

	if err := v.ReadInConfig(); err != nil {
		return cosmcfg.Config{}, fmt.Errorf("failed to read in app.toml: %w", err)
	}

	config := cosmcfg.GetConfig(v)

	return config, nil
}

func writeCosmosConfig(homeDir string, config *cosmcfg.Config) {
	tmConfigPath := filepath.Join(homeDir, "config", "app.toml")
	cosmcfg.WriteConfigFile(tmConfigPath, config)
}

func readTmConfig(homeDir string) (tmcfg.Config, error) {
	v := viper.New()

	v.SetConfigType("toml")
	v.SetConfigName("config")
	v.AddConfigPath(filepath.Join(homeDir, "config"))

	if err := v.ReadInConfig(); err != nil {
		return tmcfg.Config{}, fmt.Errorf("failed to read in config.toml: %w", err)
	}

	var config tmcfg.Config
	err := v.Unmarshal(&config)
	if err != nil {
		return tmcfg.Config{}, err
	}

	return config, nil
}

func writeTmConfig(homeDir string, config *tmcfg.Config) {
	tmConfigPath := filepath.Join(homeDir, "config", "config.toml")
	tmcfg.WriteConfigFile(tmConfigPath, config)
}