package gopsql

import "fmt"

type Builder struct {
	sqlString string
	params    []interface{}
	limit     int
	offset    int
}

func NewB() *Builder {
	return &Builder{}
}

func (b *Builder) Add(statement string, args ...interface{}) {
	b.sqlString += " " + statement
	b.params = append(b.params, args...)
}

func (b *Builder) RawSql() string {
	sql := b.sqlString
	if b.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return sql
}

func (b *Builder) Values() []interface{} {
	return b.params
}

func (b *Builder) Pagination(limit int, offset int) {
	if limit >= 0 {
		b.limit = limit
	}
	if offset >= 0 {
		b.offset = offset
	}
}
