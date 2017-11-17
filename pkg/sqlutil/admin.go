package sqlutil

import (
	"database/sql"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// CheckBinlogFormat checks binlog format
func CheckBinlogFormat(db *sql.DB) error {
	rows, err := db.Query(`SHOW GLOBAL VARIABLES LIKE "binlog_format";`)
	if err != nil {
		return errors.Trace(err)
	}
	defer rows.Close()

	// Show an example.
	/*
			   mysql> SHOW GLOBAL VARIABLES LIKE "binlog_format";
		       +---------------+-------+
		       | Variable_name | Value |
		       +---------------+-------+
		       | binlog_format | ROW   |
		       +---------------+-------+
	*/
	for rows.Next() {
		var (
			variable string
			value    string
		)

		err = rows.Scan(&variable, &value)

		if err != nil {
			return errors.Trace(err)

		}

		if variable == "binlog_format" && value != "ROW" {
			log.Fatalf("binlog_format is not 'ROW': %v", value)
		}
	}

	if rows.Err() != nil {
		return errors.Trace(rows.Err())
	}

	return nil
}
