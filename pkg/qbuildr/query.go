package qbuildr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func (qb *QueryBuilder) Entity(entity interface{}) *QueryBuilder {
	// get fields from struct type if qb.entity is exist
	if entity != nil {
		entityTyp := reflect.TypeOf(entity)
		for i := 0; i < entityTyp.NumField(); i++ {
			qb.columns = append(qb.columns, StringPascalToSnake(entityTyp.Field(i).Name))
		}
	}

	return qb
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selects = columns
	return qb
}

func (qb *QueryBuilder) Where(conditions ...string) *QueryBuilder {
	qb.wheres = append(qb.wheres, conditions...)
	return qb
}

func (qb *QueryBuilder) WhereStruct(data interface{}) *QueryBuilder {
	return qb.processStruct(data)
}

func (qb *QueryBuilder) processStruct(data interface{}) *QueryBuilder {
	// convert struct to map string interface
	var dataMap map[string]interface{}
	dataJson, _ := json.Marshal(data)
	json.Unmarshal(dataJson, &dataMap)

	// loop through map
	for key, value := range dataMap {
		// if type is map, recursively process it with updated prefix
		if reflect.TypeOf(value).Kind() == reflect.Map {
			qb.processStruct(value)
			continue
		}

		// if theres a start_date
		if key == "start_date" {
			condition := fmt.Sprintf("created_at >= '%v'", value)
			qb.wheres = append(qb.wheres, condition)
			continue
		}

		// if theres a end_date
		if key == "end_date" {
			condition := fmt.Sprintf("created_at <= '%v'", value)
			qb.wheres = append(qb.wheres, condition)
			continue
		}

		// if columns is not empty, check if column is in columns
		if len(qb.columns) > 0 {
			// if column is not in columns, skip
			if !StringInSlice(qb.columns, key) {
				continue
			}
		}

		// if type is string, wrap with single quote
		if reflect.TypeOf(value).Kind() == reflect.String {
			condition := fmt.Sprintf("%s = '%v'", key, value)
			qb.wheres = append(qb.wheres, condition)
			continue
		}

		// if type is int, float, or bool, no need to wrap with single quote
		condition := fmt.Sprintf("%s = %v", key, value)
		qb.wheres = append(qb.wheres, condition)
	}
	return qb
}

func (qb *QueryBuilder) In(column string, data []string) *QueryBuilder {
	values := make([]string, 0)

	for i := range data {
		// if type is string, wrap with single quote
		if reflect.TypeOf(data[i]).Kind() == reflect.String {
			values = append(values, fmt.Sprintf("'%v'", data[i]))
			continue
		}

		// if type is int, float, or bool, no need to wrap with single quote
		values = append(values, fmt.Sprintf("%v", data[i]))
	}

	condition := fmt.Sprintf("%s IN (%s)", column, strings.Join(values, ", "))
	qb.wheres = append(qb.wheres, condition)
	return qb
}

func (qb *QueryBuilder) OrderBy(columns ...string) *QueryBuilder {
	qb.orderBys = columns
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) Insert(data interface{}) *QueryBuilder {
	qb.insertData = data
	return qb
}

func (qb *QueryBuilder) Update(data interface{}) *QueryBuilder {
	qb.updateData = data
	return qb
}

func (qb *QueryBuilder) Delete() *QueryBuilder {
	qb.deleteData = true
	return qb
}
