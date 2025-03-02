package utils

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dkhvan-dev/product-service/internal/common"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

var structCache sync.Map
var structFieldsCache sync.Map

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if prefix == "content" {
		return name
	}

	if len(prefix) > 0 {
		return prefix + "." + name
	}

	return name
}

func SelectedFields(ctx context.Context) []string {
	return GetPreloads(ctx)
}

func StructFields(tType reflect.Type) []string {
	if cached, ok := structFieldsCache.Load(tType); ok {
		return cached.([]string)
	}

	var fields []string

	for i := 0; i < tType.NumField(); i++ {
		fields = append(fields, tType.Field(i).Name)
	}

	return fields
}

func StructFieldsMap(tType reflect.Type) map[string]string {
	if cached, ok := structCache.Load(tType); ok {
		return cached.(map[string]string)
	}

	fieldsMap := make(map[string]string)

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		dbTag := field.Tag.Get("db")
		isParentTag := field.Tag.Get("isParent")

		if dbTag != "" && isParentTag == "" {
			fieldsMap[firstLetterToLower(field.Name)] = dbTag
		}
	}

	structCache.Store(tType, fieldsMap)
	return fieldsMap
}

func GenerateSQL[T common.TableNameInterface](fields []string) string {
	var t T
	tableName := t.TableName()
	tType := reflect.TypeOf(t)
	fieldsMap := StructFieldsMap(tType)
	joinMap := sync.Map{}

	var columns []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, field := range fields {
		parts := strings.Split(field, ".")

		if len(parts) == 1 {
			if dbTag, exists := fieldsMap[field]; exists {
				columns = append(columns, fmt.Sprintf("%s.%s as \"%s\"", tableName, dbTag, dbTag))
			}
		} else {
			parent := parts[0]
			child := parts[1]
			wg.Add(1)

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
								joinColumn := structField.Tag.Get("db")
								val := fmt.Sprintf("inner join %s on %s.id = %s.%s",
									nestedTable, nestedTable, tableName, joinColumn)

								joinMap.LoadOrStore(parent, val)
								nestedStructFieldMap := StructFieldsMap(nestedType)

								if nestedTag, exists := nestedStructFieldMap[child]; exists {
									mu.Lock()
									columns = append(columns, fmt.Sprintf("%s.%s as \"%s.%s\"",
										nestedTable, nestedTag, joinColumn, nestedTag))

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

	if len(columns) == 0 {
		columns = append(columns, fmt.Sprintf("%s.id", tableName))
	}

	joinList := []string{}
	joinMap.Range(func(_, join interface{}) bool {
		joinList = append(joinList, join.(string))
		return true
	})

	var query strings.Builder
	query.WriteString("select ")
	query.WriteString(strings.Join(columns, ", "))
	query.WriteString(" from ")
	query.WriteString(tableName + " ")
	query.WriteString(strings.Join(joinList, " "))

	return query.String()
}

func firstLetterToLower(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}
