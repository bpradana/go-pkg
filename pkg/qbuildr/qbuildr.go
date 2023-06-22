package qbuildr

type QueryBuilder struct {
	table    string
	columns  []string
	selects  []string
	wheres   []string
	orderBys []string
	limit    int
	offset   int

	insertData interface{}
	updateData interface{}
	deleteData interface{}
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:   table,
		columns: make([]string, 0),
	}
}
