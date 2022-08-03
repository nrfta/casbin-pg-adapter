package migrations

import (
	"database/sql"
	"fmt"
)

const createSchemaQuery = "CREATE SCHEMA IF NOT EXISTS %s"

func init() {
	migrationFuncs[2] = func(schemaName, tableName string, tx *sql.Tx) error {
		_, err := tx.Exec(fmt.Sprintf(createSchemaQuery, schemaName))
		return err
	}
}
