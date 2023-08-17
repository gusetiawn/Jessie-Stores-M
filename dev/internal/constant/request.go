package constant

import "git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"

type StoreStateRequest struct {
	ReferenceID string `json:"reference_id"`
	Type        string `json:"type"`
}

type SelectStoreRequest struct {
	ID int64 `json:"id"`
}

type GetNearestStoresRequest struct {
	Token        *model.UserToken
	Name         string
	CategoryType string
	CategoryID   int
	Limit        int
}
