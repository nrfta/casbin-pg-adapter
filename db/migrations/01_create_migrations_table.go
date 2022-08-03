package migrations

import (
	"database/sql"
	"fmt"
)

const createMigrationsTable = `
CREATE TABLE IF NOT EXISTS %s.migrations (
  version integer NOT NULL,
  applied_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  failed boolean
)
`

func init() {
	migrationFuncs[1] = func(schemaName, tableName string, tx *sql.Tx) error {
		_, err := tx.Exec(fmt.Sprintf(createMigrationsTable, schemaName))
		return err
	}
}
