package sqlutil

import (
	"database/sql"
	"fmt"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// DBConfig is the DB configuration.
type DBConfig struct {
	Host     string `toml:"host" json:"host"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
	Port     int    `toml:"port" json:"port"`
}

// CreateDB creates a connection to MySQL database.
func CreateDB(cfg *DBConfig, timeout string) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("db config is nil")
	}
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&interpolateParams=true&readTimeout=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, timeout)
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return db, nil
}

// CloseDB closes the connection to MySQL database.
func CloseDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return errors.Trace(db.Close())
}

// CloseDBS closes connections to MySQL databases.
func CloseDBs(dbs ...*sql.DB) {
	for _, db := range dbs {
		if err := CloseDB(db); err != nil {
			log.Errorf("close db failed: %v", err)
		}
	}
}
