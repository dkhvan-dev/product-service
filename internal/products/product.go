package products

import "encoding/json"

type ProductUpsert struct {
	Category string          `json:"category"`
	Input    json.RawMessage `json:"input"`
}
