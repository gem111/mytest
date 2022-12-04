package homework_delete

type Query struct {
	SQL  string
	Args []any
}

type QueryBuilder interface {
	Build() (*Query, error)
}
