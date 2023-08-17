package repository

import (
	"context"
	"database/sql"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

const GET_USER_LAST_TOKEN_QUERY = `
WITH cte AS (
	SELECT 
		user_id,
		merchant_id
	FROM 
		user_tokens 
	WHERE 
		token = $1
  )
SELECT 
	ut.id, 
	ut.token, 
	ut.merchant_id, 
	ut.user_id, 
	ut.address_id, 
	ut.state, 
	ut.store_id, 
	ut.expire_at, 
	ut.created_at, 
	ut.updated_at, 
	ut.deleted_at 
FROM 
	user_tokens ut 
	RIGHT JOIN cte ON cte.user_id = ut.user_id 
WHERE
  	ut.merchant_id = cte.merchant_id
	AND ut.token IS NOT NULL
	AND ut.deleted_at IS NULL
ORDER BY 
	ut.id desc 
LIMIT 
	1`

func (ths *repository) GetUserLastToken(ctx context.Context, requestToken string) (*model.UserToken, error) {
	var token model.UserToken
	err := ths.db.
		QueryRowContext(ctx, GET_USER_LAST_TOKEN_QUERY, requestToken).
		Scan(&token.ID, &token.Token, &token.MerchantID, &token.UserID, &token.AddressID, &token.State, &token.StoreID, &token.ExpireAt, &token.CreatedAt, &token.UpdatedAt, &token.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}
