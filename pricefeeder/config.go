package pricefeeder

import (
	"fmt"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

const (
	DefaultConfigTemplate = `
[pricefeeder]
# Path to price feeder config file.
config_path = ""

# Log level of price feeder process.
log_level = "info"

# Enable the price feeder.
enable = false
`
)

const (
	FlagConfigPath        = "pricefeeder.config_path"
	FlagLogLevel          = "pricefeeder.log_level"
	FlagEnablePriceFeeder = "pricefeeder.enable"
)

// AppConfig defines the app configuration for the price feeder that must be set in the app.toml file.
type AppConfig struct {
	ConfigPath string `mapstructure:"config_path"`
	LogLevel   string `mapstructure:"log_level"`
	Enable     bool   `mapstructure:"enable"`
}

// ValidateBasic performs basic validation of the price feeder app config.
func (c *AppConfig) ValidateBasic() error {
	if c.ConfigPath == "" {
		return fmt.Errorf("path to price feeder config must be set")
	}

	return nil
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (AppConfig, error) {
	var (
		cfg AppConfig
		err error
	)

	if v := opts.Get(FlagConfigPath); v != nil {
		if cfg.ConfigPath, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	if v := opts.Get(FlagLogLevel); v != nil {
		if cfg.LogLevel, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	if v := opts.Get(FlagEnablePriceFeeder); v != nil {
		if cfg.Enable, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return cfg, err
	}

	return cfg, err
}
