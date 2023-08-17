package repository

import (
	"context"
	"fmt"
)

const UPDATE_STORE_ID_TOKEN_QUERY = `
UPDATE 
    user_tokens 
SET 
    store_id = $1 
WHERE
    id = $2`

func (ths *repository) UpdateTokenStore(ctx context.Context, tokenID int64, storeID int64) error {
	a, err := ths.db.ExecContext(ctx, UPDATE_STORE_ID_TOKEN_QUERY, storeID, tokenID)
	if err != nil {
		return err
	}
	fmt.Println("a.RowsAffected()")
	fmt.Println(a.RowsAffected())
	fmt.Println(storeID)
	fmt.Println(tokenID)

	return nil
}
