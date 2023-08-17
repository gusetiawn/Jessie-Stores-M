package repository

import (
	"context"
	"database/sql"
	"fmt"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

const GET_ADDRESS_QUERY = `
SELECT
    location_lat,
    location_long
FROM
    address
WHERE
    id = $1
    AND deleted_at IS NULL`

func (ths *repository) GetAddress(ctx context.Context, addressID int64) (*model.Address, error) {
	var address model.Address
	err := ths.db.
		QueryRowContext(ctx, GET_ADDRESS_QUERY, addressID).
		Scan(&address.LocationLat, &address.LocationLong)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed executing GET_ADDRESS_QUERY: %s", err.Error())
	}

	return &address, nil
}
