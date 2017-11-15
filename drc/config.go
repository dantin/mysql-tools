package drc

import (
	"flag"

	"github.com/dantin/mysql-tools/pkg/logutil"
	"github.com/dantin/mysql-tools/pkg/sqlutil"
)

// Config is the configuration.
type Config struct {
	*flag.FlagSet `json:"-"`

	// Log related configuration.
	Log logutil.LogConfig `toml:"log" json:"log"`
	// Source is the data source.
	Source *sqlutil.DBConfig `toml:"source" json:"source"`

	configFile string
	Version    bool
}

// NewConfig creates a new config.
func NewConfig() *Config {
	cfg := &Config{}
	cfg.FlagSet = flag.NewFlagSet("drc", flag.ContinueOnError)
	fs := cfg.FlagSet

	fs.BoolVar(&cfg.Version, "V", false, "print version and exit")
	fs.StringVar(&cfg.configFile, "config", "", "path to the config file")

	fs.StringVar(&cfg.Log.Level, "L", "info", "log level: debug, info, warn, error, fatal (default 'info')")
	fs.StringVar(&cfg.Log.File.Filename, "log-file", "", "log file path")

	return cfg
}
