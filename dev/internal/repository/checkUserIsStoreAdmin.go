package repository

import (
	"context"
	"database/sql"
)

const CHECK_USER_IS_STORE_ADMIN_QUERY = `
SELECT 
    true AS is_admin 
FROM 
    users 
WHERE
	id = $1
    AND store_admin = $2 
	AND deleted_at IS NULL`

func (ths *repository) CheckUserIsStoreAdmin(ctx context.Context, userID int64, storeID int64) (bool, error) {
	var isAdmin sql.NullBool
	err := ths.db.
		QueryRowContext(ctx, CHECK_USER_IS_STORE_ADMIN_QUERY, userID, storeID).
		Scan(&isAdmin)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if !isAdmin.Valid {
		return false, nil
	}

	return true, nil
}
