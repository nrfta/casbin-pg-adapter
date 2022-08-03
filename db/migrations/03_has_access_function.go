package migrations

import (
	"database/sql"
	"fmt"
)

const hasAccessFunction = `
CREATE OR REPLACE FUNCTION %[1]s.has_access (
	IN user_id text,
	IN domain text,
	IN resources text[],
	IN action text
) RETURNS boolean AS $$
	BEGIN
        RETURN EXISTS(
            WITH roles AS (
                SELECT v1 as role
                FROM %[1]s.casbin_rules
                WHERE p_type = 'g'
                  AND v0 = user_id
                  AND (v2 = domain OR v2 = '*')
            )
            SELECT *
            FROM %[1]s.casbin_rules r
            LEFT JOIN roles ON r.v0 = roles.role
            WHERE (r.p_type = 'p'
                AND (r.v0 = user_id OR r.v0 = roles.role OR (r.v0 = 'USER' AND roles.role = 'ADMIN'))
                AND (r.v1 = domain OR (position('*' in r.v1) > 0 AND starts_with(domain, rtrim(r.v1, '*'))))
                AND r.v2 = ANY(resources)
                AND (r.v3 = '*' OR r.v3 = action OR (r.v3 = 'update' AND action = 'read')))
                 OR roles.role = 'SUPERADMIN'
        );
    END;
$$ LANGUAGE plpgsql
STABLE
PARALLEL SAFE;
`

func init() {
	migrationFuncs[3] = func(schemaName, tableName string, tx *sql.Tx) error {
		if _, err := tx.Exec(fmt.Sprintf(hasAccessFunction, schemaName)); err != nil {
			return err
		}
		return nil
	}
}
