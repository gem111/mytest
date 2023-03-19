package orm

import "database/sql"

type DBOption func(*DB)

type DB struct {
	r  *registry
	db *sql.DB
}

func Open(driver, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)

	if err != nil {
		return nil, err
	}
	return OpenDb(db, opts...)
}

// OpenDb 方便进行测试  和 用户传入其他的db
func OpenDb(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:  &registry{},
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

// MustOpenDB 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustOpenDB(driver, dataSourceName string, opts ...DBOption) *DB {
	db, err := Open(driver, dataSourceName, opts...)
	if err != nil {
		panic(err)
	}
	return db
}
