package migrations

import (
	"database/sql"
)

const createMigrationsTable = `
CREATE TABLE migrations (
  version integer NOT NULL,
  applied_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  failed boolean
)
`

func init() {
	migrationFuncs[1] = func(schemaName, tableName string, tx *sql.Tx) error {
		_, err := tx.Exec(createMigrationsTable)
		return err
	}
}
