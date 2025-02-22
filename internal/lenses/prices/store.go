package prices

import (
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
)

func CreateActualPrice(tx *sqlx.Tx, price float64, lensId int) (*int, *errors.CustomError) {
	query := `
			insert into lenses_price_history(price, lens_id) 
			values ($1, $2)
			returning id
		`

	var actualPriceId int
	err := tx.QueryRowx(query, price, lensId).Scan(&actualPriceId)

	if err != nil {
		config.Logger.Error("Failed create actual lens price", zap.String("db", err.Error()))
		return nil, errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return &actualPriceId, nil
}
