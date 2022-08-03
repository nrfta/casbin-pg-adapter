package migrations

import (
	"database/sql"
	"fmt"
)

const casbinRulesTable = `
CREATE TABLE IF NOT EXISTS "%s"."%s" (
		p_type varchar(256) not null default '',
		v0 		varchar(256) not null default '',
		v1 		varchar(256) not null default '',
		v2 		varchar(256) not null default '',
		v3 		varchar(256) not null default '',
		v4 		varchar(256) not null default '',
		v5 		varchar(256) not null default ''
)
`

const CasbinRulesTableIndex = `
CREATE INDEX IF NOT EXISTS idx_%[2]s_%[3]s ON "%[1]s"."%[2]s" (%[3]s)
`

func init() {
	migrationFuncs[2] = func(schemaName, tableName string, tx *sql.Tx) error {
		if _, err := tx.Exec(fmt.Sprintf(casbinRulesTable, schemaName, tableName)); err != nil {
			return err
		}
		columns := [7]string{
			"p_type",
			"v0",
			"v1",
			"v2",
			"v3",
			"v4",
			"v5",
		}
		for _, column := range columns {
			if _, err := tx.Exec(fmt.Sprintf(CasbinRulesTableIndex, schemaName, tableName, column)); err != nil {
				return err
			}
		}
		return nil
	}
}
