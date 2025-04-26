package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/constants"
	customErrors "github.com/dkhvan-dev/web-commons/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

const CmsURL = "http://localhost:8080/v1/graphql"

var cmsCache = make(map[string]map[string]bool)

func ValidateCmsValue(data interface{}) *customErrors.CustomError {
	v := reflect.ValueOf(data)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("cms")
		if tag == "" {
			continue
		}

		fieldValue := value.String()

		if fieldValue == "" {
			continue
		}

		if validValues, ok := cmsCache[tag]; ok {
			if _, exists := validValues[fieldValue]; exists {
				continue
			}
		} else {
			validValues, err := fetchCmsValues(tag, value.String())
			if err != nil {
				config.Logger.Error("Failed to fetch cms value", zap.Error(err))
				return customErrors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}

			cmsCache[tag] = validValues
		}

		if _, exists := cmsCache[tag][fieldValue]; !exists {
			config.Logger.Error("Field value not found in CMS",
				zap.String("table", tag), zap.String("field_value", fieldValue))

			jsonTag := field.Tag.Get("json")
			errMsg := fmt.Sprintf("Field %s value not found in CMS", jsonTag)
			return &customErrors.CustomError{Key: "VALIDATION_ERROR", Code: 512, Msg: &errMsg}
		}
	}

	return nil
}

func fetchCmsValues(table string, code string) (map[string]bool, error) {
	query := fmt.Sprintf(`{"query":"query {%s(where: {code: {_in: %s}}){code}}"}`, table, code)
	req, err := http.NewRequest(http.MethodPost, CmsURL, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}

	req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	req.Header.Set("X-Hasura-Admin-Secret", "mysecret")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response struct {
		Data map[string][]struct {
			Code string `json:"code"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	validValues := make(map[string]bool)
	for _, entry := range response.Data[table] {
		validValues[entry.Code] = true
	}

	return validValues, nil
}

func ToString(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case primitive.ObjectID:
		return v.Hex()
	case primitive.Decimal128:
		return v.String()
	case *primitive.Decimal128:
		if v != nil {
			return v.String()
		}
	}

	return ""
}
