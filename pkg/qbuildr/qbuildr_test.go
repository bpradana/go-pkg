package qbuildr_test

import (
	"go-pkg/pkg/qbuildr"
	"testing"
)

func TestNewQueryBuilder_Where(t *testing.T) {
	query := qbuildr.NewQueryBuilder("transactions").Select("*").Where("id = 1").Where("user_name = 'John'").Limit(1).Build()
	expected := "SELECT * FROM transactions WHERE id = 1 AND user_name = 'John' LIMIT 1"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_EntityWhereStruct(t *testing.T) {
	type User struct {
		ID       string `json:"id,omitempty"`
		UserName string `json:"user_name,omitempty"`
		Age      int    `json:"age,omitempty"`
	}

	var entity User

	structParams := struct {
		ID        string `json:"id,omitempty"`
		UserName  string `json:"user_name,omitempty"`
		Age       int    `json:"age,omitempty"`
		StartDate string `json:"start_date,omitempty"`
		EndDate   string `json:"end_date,omitempty"`
		Q         string `json:"q,omitempty"`
	}{
		ID:        "1",
		UserName:  "John",
		Age:       20,
		StartDate: "2020-01-01",
		EndDate:   "2020-01-31",
		Q:         "test",
	}

	query := qbuildr.NewQueryBuilder("transactions").Entity(entity).Select("*").WhereStruct(structParams).Build()
	expected := "SELECT * FROM transactions WHERE user_name = 'John' AND age = 20 AND created_at >= '2020-01-01' AND created_at <= '2020-01-31' AND id = '1'"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_WhereStruct(t *testing.T) {
	structParams := struct {
		ID        string `json:"id,omitempty"`
		UserName  string `json:"user_name,omitempty"`
		Age       int    `json:"age,omitempty"`
		StartDate string `json:"start_date,omitempty"`
		EndDate   string `json:"end_date,omitempty"`
	}{
		ID:        "1",
		UserName:  "John",
		Age:       20,
		StartDate: "2020-01-01",
		EndDate:   "2020-01-31",
	}

	query := qbuildr.NewQueryBuilder("transactions").Select("*").WhereStruct(structParams).Build()
	expected := "SELECT * FROM transactions WHERE user_name = 'John' AND age = 20 AND created_at >= '2020-01-01' AND created_at <= '2020-01-31' AND id = '1'"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_WhereStruct_Recursive(t *testing.T) {
	type User struct {
		ID       string `json:"id,omitempty"`
		UserName string `json:"user_name,omitempty"`
		Age      int    `json:"age,omitempty"`
		Foo      int    `json:"foo,omitempty"`
		Bar      string `json:"bar,omitempty"`
	}

	structData := User{
		ID:       "1",
		UserName: "John",
		Age:      20,
	}

	structMeta := struct {
		StartDate string `json:"start_date,omitempty"`
		EndDate   string `json:"end_date,omitempty"`
	}{
		StartDate: "2020-01-01",
		EndDate:   "2020-01-31",
	}

	structParams := struct {
		Data interface{} `json:"data,omitempty"`
		Meta interface{} `json:"meta,omitempty"`
		Foo  int         `json:"foo,omitempty"`
		Bar  string      `json:"bar,omitempty"`
	}{
		Data: structData,
		Meta: structMeta,
		Foo:  1,
		Bar:  "test",
	}

	query := qbuildr.NewQueryBuilder("transactions").Entity(User{}).Select("*").WhereStruct(structParams).Build()
	expected := "SELECT * FROM transactions WHERE id = '1' AND user_name = 'John' AND age = 20 AND created_at >= '2020-01-01' AND created_at <= '2020-01-31' AND foo = 1 AND bar = 'test'"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_Insert(t *testing.T) {
	query := qbuildr.NewQueryBuilder("transactions").Insert(struct {
		ID       string `json:"id,omitempty"`
		UserName string `json:"user_name,omitempty"`
		Age      int    `json:"age,omitempty"`
	}{
		ID:       "1",
		UserName: "John",
		Age:      20,
	}).Build()
	expected := "INSERT INTO transactions(id, user_name, age) VALUES ('1', 'John', 20)"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_Update(t *testing.T) {
	query := qbuildr.NewQueryBuilder("transactions").Update(struct {
		ID       string `json:"id,omitempty"`
		UserName string `json:"user_name,omitempty"`
		Age      int    `json:"age,omitempty"`
	}{
		ID:       "1",
		UserName: "John",
		Age:      20,
	}).Where("id = 1").Build()
	expected := "UPDATE transactions SET id = '1', user_name = 'John', age = 20 WHERE id = 1"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_Delete(t *testing.T) {
	query := qbuildr.NewQueryBuilder("transactions").Delete().Where("id = 1").Build()
	expected := "DELETE FROM transactions WHERE id = 1"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}

func TestNewQueryBuilder_In(t *testing.T) {
	type User struct {
		ID       string `json:"id,omitempty"`
		UserName string `json:"user_name,omitempty"`
	}

	structData := User{
		ID:       "1",
		UserName: "John",
	}

	structMeta := struct {
		Groups []string `json:"groups,omitempty"`
	}{
		Groups: []string{"A", "B", "C"},
	}

	structParams := struct {
		Data interface{} `json:"data,omitempty"`
		Meta interface{} `json:"meta,omitempty"`
	}{
		Data: structData,
		Meta: structMeta,
	}

	query := qbuildr.NewQueryBuilder("transactions").Entity(User{}).Select("*").WhereStruct(structParams).In("groups", structMeta.Groups).Build()
	expected := "SELECT * FROM transactions WHERE id = '1' AND user_name = 'John' AND groups IN ('A', 'B', 'C')"
	if query != expected {
		t.Errorf("expected %s, got %s", expected, query)
	}
}
