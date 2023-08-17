package constant

type NearestStore struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Picture      *string  `json:"picture"`
	Address      *string  `json:"address"`
	Distance     string   `json:"distance,omitempty"`
	LocationLat  *float64 `json:"location_lat"`
	LocationLong *float64 `json:"location_long"`
	IsOpen       bool     `json:"is_open"`
	Dist         float64  `json:"-"`
}
