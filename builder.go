package gopsql

import "fmt"

type Builder struct {
	selectString string
	selectParams []interface{}
	sqlString    string
	params       []interface{}
	limit        int
	offset       int
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Add(statement string, args ...interface{}) {
	b.sqlString += " " + statement
	b.params = append(b.params, args...)
}

func (b *Builder) Select(statement string, args ...interface{}) {
	b.selectString = statement
	b.selectParams = append(make([]interface{}, 0), args...)
}

func (b *Builder) RawSql() string {
	
	var sql string 
	if b.selectString != "" {
		sql +=  "SELECT " + b.selectString
	}
	sql += " " + b.sqlString
	if b.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return sql
}

func (b *Builder) Values() []interface{} {
	res := make([]interface{}, 0)
	res = append(res, b.selectParams...)
	res = append(res, b.params...)
	return res
}

func (b *Builder) Pagination(limit int, offset int) {
	if limit >= 0 {
		b.limit = limit
	}
	if offset >= 0 {
		b.offset = offset
	}
}
