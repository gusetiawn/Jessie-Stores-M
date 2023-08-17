package repository

//go:generate mockgen -source=./repository.go -destination=./mock/mock.go -package=mock

import (
	"context"
	"database/sql"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

type Repository interface {
	// GetUserLastToken mengambil token versi terbaru (last token) dari input token
	// yang diberikan. Jika input token merupakan token yang tidak valid (not exist),
	// maka last token menjadi nil.
	GetUserLastToken(ctx context.Context, token string) (*model.UserToken, error)

	// GetStoreByRef mengambil informasi suatu store berdasarkan reference tertentu
	// yang berasal dari tabel deep_links
	GetStoreByRef(ctx context.Context, typeDeeplink string, referenceID string) (*model.Store, error)

	// UpdateTokenStore mengisi field store_id pada tabel user_tokens untuk tokenID
	// tertentu.
	UpdateTokenStore(ctx context.Context, tokenID, storeID int64) error

	// GetAddress mengambil informasi location_lat dan location_long suatu alamat yang
	// terdapat pada tabel address untuk id address tertentu.
	GetAddress(ctx context.Context, addressID int64) (*model.Address, error)

	// GetNearestStores mengambil informasi stores yang memenuhi kriteria
	// 1. store milik merchant tertentu (sesuai input merchantID)
	// 2. store menjual produk berstatus aktif, demikian juga store itu sendiri berstatus
	//    aktif
	// 3. kategori store selain 'Makanan dan Minuman'
	// 4. admin dari store bukanlah user yang melakukan get (sesuai input userID)
	// 5. store ber-category_id tertentu (sesuai input categoryID, jika ada)
	// 6. store memiliki nama, atau menjual produk, atau menjual produk yang kategorinya,
	//    atau store itu sendiri termasuk ke dalam kategori yang mengandung string tertentu
	//    (sesuai input name, jika ada)
	GetNearestStores(ctx context.Context, input GetNearestStoresInput) ([]model.StoreWithCategory, error)

	// GetStoreOpenConfig mengembalikan key-value pair yang menyatakan konfigurasi status
	// buka/tutup suatu store, yaitu dari key closed_<storeID> dan service_<storeID>.
	GetStoreOpenConfig(ctx context.Context, storeIDs []int64) (model.StoreConfig, error)

	// CheckUserIsStoreAdmin mengembalikan nilai true jika suatu user adalah admin dari
	// suatu store, dan sebaliknya.
	CheckUserIsStoreAdmin(ctx context.Context, userID int64, storeID int64) (bool, error)
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
