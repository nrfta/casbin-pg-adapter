package migrations

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/neighborly/go-errors"
	"github.com/neighborly/go-pghelpers"
)

var migrationFuncs = make(map[int]func(schemaName, tableName string, tx *sql.Tx) error)

func Migrate(schemaName, tableName string, db *sql.DB) error {
	latest, err := getLatest(schemaName, db)
	if err != nil {
		return err
	}

	ver := latest + 1
	fn, ok := migrationFuncs[ver]
	if !ok {
		return nil
	}
	var txErr error
	for ok {
		err = pghelpers.ExecInTx(db, func(tx *sql.Tx) bool {
			if txErr = fn(schemaName, tableName, tx); txErr != nil {
				return false
			}
			return true
		})
		if err != nil {
			return errors.Wrap(err, "failed to start migration transaction")
		}
		if txErr != nil {
			if ver > 1 {
				err2 := insertVersion(ver, true, db)
				if err2 != nil {
					txErr = errors.Wrap(txErr, err2.Error())
				}
			}
			return txErr
		}

		if txErr = insertVersion(ver, false, db); txErr != nil {
			return txErr
		}

		ver++
		fn, ok = migrationFuncs[ver]
	}
	return nil
}

func getLatest(schemaName string, db *sql.DB) (int, error) {
	const query = `SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = $1 AND tablename  = 'migrations')`
	var exists bool
	err := db.QueryRow(query, schemaName).Scan(&exists)
	if err != nil {
		return -1, err
	}
	if !exists {
		return 0, nil
	}

	const query2 = `SELECT MAX(version) FROM %s.migrations`
	var latest int
	err = db.QueryRow(fmt.Sprintf(query2, schemaName)).Scan(&latest)
	if err != nil {
		if err.Error() == `sql: Scan error on column index 0, name "max": converting NULL to int is unsupported` {
			return 1, nil
		}
		return -1, err
	}
	return latest, nil
}

func insertVersion(ver int, fail bool, db *sql.DB) error {
	const query = `INSERT INTO migrations (version, applied_at, failed) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, ver, time.Now(), fail)
	return err
}
