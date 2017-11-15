package main

import (
	"flag"
	"os"

	"github.com/dantin/mysql-tools/drc"
	"github.com/dantin/mysql-tools/pkg/logutil"
	"github.com/dantin/mysql-tools/pkg/utils"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// main is the bootstrap.
func main() {
	cfg := drc.NewConfig()
	err := cfg.Parse(os.Args[1:])
	if cfg.Version {
		utils.PrintRawInfo("gravity")
		os.Exit(0)
	}
	switch errors.Cause(err) {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		log.Fatalf("parse cmd flags errors: %s", err)
	}

	err = logutil.InitLogger(&cfg.Log)
	if err != nil {
		log.Fatalf("initialize log error: %s", err)
	}

	utils.LogRawInfo("gravity")
}
