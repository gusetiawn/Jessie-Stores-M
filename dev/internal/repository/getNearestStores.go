package repository

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

type GetNearestStoresInput struct {
	UserID     int64
	MerchantID int64
	CategoryID int
	Name       string
}

const GET_NEAREST_STORE_QUERY = `
WITH sa_cte AS (
	SELECT
		COALESCE(store_admin, 0) AS store_admin
	FROM
		users
	WHERE
		id = $1
)
SELECT 
  	DISTINCT s.id, 
  	s.name, 
  	s.address, 
 	s.city, 
  	s.country, 
  	s.max_radius, 
  	s.pic, 
  	s.location_lat, 
  	s.location_long, 
  	s.merchant_id, 
  	s.zip_code, 
  	s.category_id, 
  	sc.name AS category_name 
FROM 
  	stores s
	JOIN sa_cte ON sa_cte.store_admin != s.id 
  	JOIN items i ON i.store_id = s.id 
  	JOIN categories ic ON i.category_id = ic.id 
  	JOIN store_categories sc ON sc.id = s.category_id
WHERE 
  	s.merchant_id = $2 
  	AND i.is_active = true 
  	AND s.is_active = true 
  	AND ic.is_active = true 
  	AND i.deleted_at IS NULL 
  	AND s.deleted_at IS NULL
  	AND sc.name != 'Makanan dan Minuman'`

func (ths *repository) GetNearestStores(ctx context.Context, input GetNearestStoresInput) ([]model.StoreWithCategory, error) {
	query, args := buildGetStoresQuery(input)
	rows, err := ths.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed executing GetStoresQuery: %s", err.Error())
	}
	defer rows.Close()

	stores := make([]model.StoreWithCategory, 0, 10)
	for rows.Next() {
		var store model.StoreWithCategory
		if err = rows.Scan(
			&store.Store.ID,
			&store.Store.Name,
			&store.Store.Address,
			&store.Store.City,
			&store.Store.Country,
			&store.Store.MaxRadius,
			&store.Store.Pic,
			&store.Store.LocationLat,
			&store.Store.LocationLong,
			&store.Store.MerchantID,
			&store.Store.ZipCode,
			&store.Store.CategoryID,
			&store.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("failed scanning GetStoresQuery result: %s", err.Error())
		}

		stores = append(stores, store)
	}

	return stores, nil
}

var whitespaceNormalizer = regexp.MustCompile(`\s+`)

func buildGetStoresQuery(input GetNearestStoresInput) (string, []interface{}) {
	query := strings.Builder{}
	query.WriteString(GET_NEAREST_STORE_QUERY)

	args := make([]interface{}, 0, 4)
	args = append(args, input.UserID, input.MerchantID)

	if input.Name != "" {
		args = append(args, "%"+input.Name+"%")
		placeholder := strconv.Itoa(len(args))
		query.WriteString(`
			AND (
				i.name ILIKE $` + placeholder + `
				OR ic.name ILIKE $` + placeholder + `
				OR s.name ILIKE $` + placeholder + `
				OR sc.name ILIKE $` + placeholder + `
			)`)
	}

	if input.CategoryID != 0 {
		args = append(args, input.CategoryID)
		query.WriteString(fmt.Sprintf(`
			AND s.category_id = $%d`,
			len(args)))
	}

	return whitespaceNormalizer.ReplaceAllString(query.String(), " "), args
}
