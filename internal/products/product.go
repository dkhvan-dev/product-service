package products

import "encoding/json"

type ProductUpsert struct {
	Type  string          `json:"type"`
	Input json.RawMessage `json:"input"`
}
