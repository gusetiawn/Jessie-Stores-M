package service

//go:generate mockgen -source=./service.go -destination=./mock/mock.go -package=mock

import (
	"context"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository"
)

type Service interface {
	// StoreState memperbaharui token user setelah user memilih store tertentu
	// berdasarkan referenceID (dari tabel deep_links)
	StoreState(ctx context.Context, token *model.UserToken, req constant.StoreStateRequest) (*model.Store, error)

	// GetNearestStores mengembalikan daftar stores terdekat berdasarkan kriteria
	// tertentu.
	GetNearestStores(ctx context.Context, req constant.GetNearestStoresRequest) ([]constant.NearestStore, error)

	// SelectStore memperbaharui token user setelah user memilih store tertentu
	// berdasarkan storeID yang mana user tersebut bukan sebagai adminnya.
	SelectStore(ctx context.Context, token *model.UserToken, req constant.SelectStoreRequest) error
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return &service{
		repository: repository,
	}
}
