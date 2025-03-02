package prices

import "time"

type LensPriceHistoryView struct {
	Id        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Price     float64   `json:"price" db:"price"`
}

func (l *LensPriceHistoryView) TableName() string {
	return "lenses_price_history"
}
