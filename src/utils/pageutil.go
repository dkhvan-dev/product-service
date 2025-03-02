package utils

import (
	"fmt"
	"github.com/dkhvan-dev/product-service/internal/common"
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"reflect"
	"strings"
	"sync"
)

func ToSqlSort[T any](sorts []*model.SortInput) string {
	var t T
	tType := reflect.TypeOf(t)
	sortClauses := make([]string, 0)

	if len(sorts) == 0 {
		sortClauses = append(sortClauses, "id ASC")
	}

	fieldsMap := StructFieldsMap(tType)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, s := range sorts {
		parts := strings.Split(s.Field, ".")

		if len(parts) == 1 {
			if dbAlias, exists := fieldsMap[s.Field]; exists {
				sortClauses = append(sortClauses, fmt.Sprintf("%s %s", dbAlias, s.Direction))
			}
		} else {
			wg.Add(1)
			parent := parts[0]
			child := parts[1]

			go func(parent, child string) {
				defer wg.Done()
				for i := 0; i < tType.NumField(); i++ {
					structField := tType.Field(i)

					if firstLetterToLower(structField.Name) == parent {
						nestedType := structField.Type
						if nestedType.Kind() == reflect.Ptr {
							nestedType = nestedType.Elem()
						}

						if nestedType.Kind() == reflect.Struct {
							nestedInstance := reflect.New(nestedType).Interface()
							if nestedStruct, ok := nestedInstance.(common.TableNameInterface); ok {
								nestedTable := nestedStruct.TableName()
								nestedFieldsMap := StructFieldsMap(nestedType)

								if nestedAlias, exists := nestedFieldsMap[child]; exists {
									mu.Lock()
									sortClauses = append(sortClauses, fmt.Sprintf("%s.%s %s", nestedTable, nestedAlias, s.Direction))
									mu.Unlock()
								}
							}
						}
					}
				}
			}(parent, child)
		}
	}

	wg.Wait()
	return strings.Join(sortClauses, ", ")
}

func SetDefaults(pageable *model.PageInput) {
	if pageable.Size == 0 {
		pageable.Size = 20
	}

	if pageable.Sort == nil {
		pageable.Sort = []*model.SortInput{{Field: "id", Direction: "asc"}}
	}
}
