package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

func (ths *repository) GetStoreOpenConfig(ctx context.Context, storeIDs []int64) (model.StoreConfig, error) {
	if len(storeIDs) == 0 {
		return nil, nil
	}

	query, args := buildGetStoreOpenConfigQuery(storeIDs)
	var result sql.NullString
	err := ths.db.QueryRowContext(ctx, query, args...).Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("failed executing GET_STORE_CONFIG_QUERY: %s", err.Error())
	}

	if !result.Valid {
		return nil, nil
	}

	var config model.StoreConfig
	err = json.Unmarshal([]byte(result.String), &config)
	if err != nil {
		return nil, fmt.Errorf("failed parsing GET_STORE_CONFIG_QUERY result: %s", err.Error())
	}

	return config, nil
}

func buildGetStoreOpenConfigQuery(storeIDs []int64) (string, []interface{}) {
	query := strings.Builder{}
	query.WriteString(`
	WITH cte AS (
		SELECT
			SPLIT_PART(key, '_', 2) AS store_id,
			JSON_OBJECT_AGG(
				SPLIT_PART(key, '_', 1), value
			) AS configs
		FROM
			config
		WHERE
			key IN (`)

	args := make([]interface{}, 0, len(storeIDs)*2)
	for i := range storeIDs {
		args = append(args, fmt.Sprintf("closed_%d", storeIDs[i]))
		query.WriteString("$" + strconv.Itoa(len(args)) + ", ")
		args = append(args, fmt.Sprintf("service_%d", storeIDs[i]))
		query.WriteString("$" + strconv.Itoa(len(args)))

		if i < len(storeIDs)-1 {
			query.WriteString(", ")
		} else {
			query.WriteString(")")
		}
	}

	query.WriteString(`
			AND deleted_at IS NULL
		GROUP BY
			store_id
		)
		SELECT
		JSON_OBJECT_AGG(
			store_id, configs
		) AS result
		FROM
		cte`)

	return query.String(), args
}
