package drc

import (
	"flag"

	"github.com/dantin/mysql-tools/pkg/logutil"
)

// Config is the configuration.
type Config struct {
	*flag.FlagSet `json:"-"`

	// Log related configuration.
	Log logutil.LogConfig `toml:"log" json:"log"`
	// Source is the data source.
	//Source *DBConfig `toml:"source" json:"source"`
}

// NewConfig creates a new config.
func NewConfig() *Config {
	cfg := &Config{}

	return cfg
}
