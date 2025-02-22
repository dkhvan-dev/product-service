package prices

import "time"

type ActualLensPriceEntity struct {
	Id        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Price     float64   `json:"price" db:"price"`
}
