package qbuildr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func (qb *QueryBuilder) buildInsert(args []interface{}) string {
	query := strings.Builder{}
	query.WriteString("INSERT INTO ")
	query.WriteString(qb.table)

	columns := make([]string, 0)
	values := make([]string, 0)

	// convert struct to map string interface
	var dataMap map[string]interface{}
	dataJson, _ := json.Marshal(qb.insertData)
	json.Unmarshal(dataJson, &dataMap)

	// loop through map
	for key, value := range dataMap {
		columns = append(columns, StringPascalToSnake(key))

		// if type is string, wrap with single quote
		if reflect.TypeOf(value).Kind() == reflect.String {
			values = append(values, fmt.Sprintf("'%v'", value))
			continue
		}

		// if type is int, float, or bool, no need to wrap with single quote
		values = append(values, fmt.Sprintf("%v", value))
	}

	query.WriteString(fmt.Sprintf("(%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(values, ", ")))

	return query.String()
}

func (qb *QueryBuilder) buildUpdate(args []interface{}) string {
	query := strings.Builder{}
	query.WriteString("UPDATE ")
	query.WriteString(qb.table)
	query.WriteString(" SET ")

	// convert struct to map string interface
	var dataMap map[string]interface{}
	dataJson, _ := json.Marshal(qb.updateData)
	json.Unmarshal(dataJson, &dataMap)

	updates := make([]string, 0)

	// loop through map
	for key, value := range dataMap {
		// if type is string, wrap with single quote
		if reflect.TypeOf(value).Kind() == reflect.String {
			updates = append(updates, fmt.Sprintf("%s = '%v'", key, value))
			continue
		}

		// if type is int, float, or bool, no need to wrap with single quote
		updates = append(updates, fmt.Sprintf("%s = %v", key, value))
	}

	query.WriteString(strings.Join(updates, ", "))

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.wheres, " AND "))
	}

	return query.String()
}

func (qb *QueryBuilder) buildDelete(args []interface{}) string {
	query := strings.Builder{}
	query.WriteString("DELETE FROM ")
	query.WriteString(qb.table)

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.wheres, " AND "))
	}

	return query.String()
}

func (qb *QueryBuilder) Build() string {
	query := strings.Builder{}
	var args []interface{}

	if qb.insertData != nil {
		return qb.buildInsert(args)
	} else if qb.updateData != nil {
		return qb.buildUpdate(args)
	} else if qb.deleteData != nil {
		return qb.buildDelete(args)
	}

	query.WriteString("SELECT ")

	if len(qb.selects) > 0 {
		query.WriteString(strings.Join(qb.selects, ", "))
	} else {
		query.WriteString("*")
	}

	query.WriteString(" FROM ")
	query.WriteString(qb.table)

	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.wheres, " AND "))
	}

	if len(qb.orderBys) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(qb.orderBys, ", "))
	}

	if qb.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	if qb.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return query.String()
}
