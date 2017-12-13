package replicator

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/dantin/mysql-tools/pkg/logutil"
	"github.com/dantin/mysql-tools/pkg/sqlutil"
	"github.com/juju/errors"
)

// Config is the configuration.
type Config struct {
	*flag.FlagSet `json:"-"`

	// DB related settings.
	// Source is the data source.
	Source *sqlutil.DBConfig `toml:"source" json:"source"`
	// ServerID is the slave server ID.
	ServerID int `toml:"server-id" json:"srver-id"`
	// Meta is the meta file.
	Meta string `toml:"meta" json:"meta"`
	// EnableGTID is enabled when MySQL turns GTID mod on.
	EnableGTID bool `toml:"enable-gtid" json:"enable-gtid"`
	// AutoFixGTID is enabled when it is need to fix GTID set.
	AutoFixGTID bool `toml:"auto-fix-gtid" json:"auto-fix-gtid"`

	// Log related configuration.
	Log logutil.LogConfig `toml:"log" json:"log"`

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
	fs.IntVar(&cfg.ServerID, "server-id", 101, "MySQL slave server ID")
	fs.StringVar(&cfg.Meta, "meta", "gravity.meta", "meta info")

	fs.StringVar(&cfg.Log.Level, "L", "info", "log level: debug, info, warn, error, fatal (default 'info')")
	fs.StringVar(&cfg.Log.File.Filename, "log-file", "", "log file path")
	fs.BoolVar(&cfg.Log.File.LogRotate, "log-rotate", true, "rotate log")

	return cfg
}

// Parse parses flag definitions from argument list
func (c *Config) Parse(arguments []string) error {
	// Parse first to get config file.
	err := c.FlagSet.Parse(arguments)
	if err != nil {
		return errors.Trace(err)
	}

	if c.configFile != "" {
		err = c.configFromFile(c.configFile)
		if err != nil {
			return errors.Trace(err)
		}
	}

	// Parse again to replace with command line options.
	err = c.FlagSet.Parse(arguments)
	if err != nil {
		return errors.Trace(err)
	}

	if len(c.FlagSet.Args()) != 0 {
		return errors.Errorf("'%s' is an invalid flag", c.FlagSet.Arg(0))
	}

	return nil
}

// configFromFile loads config from file.
func (c *Config) configFromFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	return errors.Trace(err)
}
