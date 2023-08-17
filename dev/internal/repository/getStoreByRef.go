package repository

import (
	"context"
	"database/sql"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

const SELECT_STORE_BY_REF_ITEM_QUERY = `
SELECT
	s.id, 
	s.merchant_id, 
	s.name, 
	s.address, 
	s.city, 
	s.country, 
	s.max_radius, 
	s.location_lat, 
	s.location_long, 
	s.pic, 
	s.is_active, 
	s.zip_code, 
	s.category_id, 
	s.created_at, 
	s.updated_at, 
	s.deleted_at
FROM 
    stores s
    LEFT JOIN items i on i.store_id = s.id
    LEFT JOIN deep_links d on d.table_id = i.id 
WHERE
	d.reference_id = $1 
 	AND d.table_name = 'items' 
 	AND s.is_active = true 
 	AND s.deleted_at IS NULL`

const SELECT_STORE_BY_REF_CAT_QUERY = `
SELECT
	s.id, 
	s.merchant_id, 
	s.name, 
	s.address, 
	s.city, 
	s.country, 
	s.max_radius, 
	s.location_lat, 
	s.location_long, 
	s.pic, 
	s.is_active, 
	s.zip_code, 
	s.category_id, 
	s.created_at, 
	s.updated_at, 
	s.deleted_at
FROM 
	stores s
	LEFT JOIN categories c on c.store_id = s.id
	LEFT JOIN deep_links d on d.table_id = c.id 
WHERE
	d.reference_id = $1
	AND d.table_name = 'categories'
	AND s.is_active = true
	AND s.deleted_at IS NULL`

const SELECT_STORE_BY_REF_STORE_QUERY = `
SELECT
	s.id, 
	s.merchant_id, 
	s.name, 
	s.address, 
	s.city, 
	s.country, 
	s.max_radius, 
	s.location_lat, 
	s.location_long, 
	s.pic, 
	s.is_active, 
	s.zip_code, 
	s.category_id, 
	s.created_at, 
	s.updated_at, 
	s.deleted_at
FROM 
	stores s
	LEFT JOIN deep_links d on d.table_id = s.id 
WHERE
	d.reference_id = $1
	AND d.table_name = 'stores'
	AND s.is_active = true
	AND s.deleted_at IS NULL`

func (ths *repository) GetStoreByRef(ctx context.Context, typeDeeplink string, referenceID string) (*model.Store, error) {
	var query string
	switch typeDeeplink {
	case constant.ITEMS_TYPE_DEEPLINK:
		query = SELECT_STORE_BY_REF_ITEM_QUERY
	case constant.CATEGORIES_TYPE_DEEPLINK:
		query = SELECT_STORE_BY_REF_CAT_QUERY
	default:
		query = SELECT_STORE_BY_REF_STORE_QUERY
	}

	var store model.Store
	err := ths.db.
		QueryRowContext(ctx, query, referenceID).
		Scan(&store.ID, &store.MerchantID, &store.Name, &store.Address, &store.City, &store.Country, &store.MaxRadius, &store.LocationLat, &store.LocationLong, &store.Pic, &store.IsActive, &store.ZipCode, &store.CategoryID, &store.CreatedAt, &store.UpdatedAt, &store.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &store, nil
}
