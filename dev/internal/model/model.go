package model

import "time"

type Address struct {
	ID           int64      `db:"id"`
	UserID       int64      `db:"user_id"`
	Name         string     `db:"name"`
	BillingName  string     `db:"billing_name"`
	Address      string     `db:"address"`
	Notes        string     `db:"notes"`
	LocationLat  float64    `db:"location_lat"`
	LocationLong float64    `db:"location_long"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	IsDefault    *bool      `db:"is_default"`
}

type UserToken struct {
	ID         int64      `db:"id"`
	Token      *string    `db:"token"`
	UserID     int64      `db:"user_id"`
	ExpireAt   *time.Time `db:"expire_at"`
	CreatedAt  *time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	State      *string    `db:"state"`
	StoreID    *int64     `db:"store_id"`
	MerchantID *int64     `db:"merchant_id"`
	AddressID  *int64     `db:"address_id"`
}

type Store struct {
	ID           int64      `db:"id"`
	Name         string     `db:"name"`
	Address      *string    `db:"address"`
	City         *string    `db:"city"`
	Country      *string    `db:"country"`
	MaxRadius    *int64     `db:"max_radius"`
	Pic          *string    `db:"pic"`
	IsActive     *bool      `db:"is_active"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	LocationLat  *float64   `db:"location_lat"`
	LocationLong *float64   `db:"location_long"`
	MerchantID   *int64     `db:"merchant_id"`
	ZipCode      *string    `db:"zip_code"`
	CategoryID   *int64     `db:"category_id"`
	MinStock     *int       `db:"min_stock"`
	MaxSalesDays *int       `db:"max_sales_days"`
}

type StoreWithCategory struct {
	Store
	CategoryName string `db:"category_name"`
}

type StoreConfig map[string]map[string]string
