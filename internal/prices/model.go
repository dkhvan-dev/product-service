package prices

import "time"

type ActualPriceEntity struct {
	Id        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Value     int       `json:"value" db:"value"`
	ProductId int       `db:"product_id"`
}
