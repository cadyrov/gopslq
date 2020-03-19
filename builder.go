package gopsql

type Builder struct {
	sqlString string
	params    []interface{}
}

func NewB() *Builder {
	return &Builder{}
}

func (b *Builder) Add(statement string, args ...interface{}) {
	b.sqlString += " " + statement
	b.params = append(b.params, args...)
}

func (b *Builder) RawSql() string {
	return b.sqlString
}

func (b *Builder) Values() []interface{} {
	return b.params
}
