package orm

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mytest/orm/internal/errs"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db := memoryDB(t)

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// From 都不调用
			name: "no from",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM
			name: "with from",
			q:    NewSelector[TestModel](db).From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "empty from",
			q:    NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name: "with db",
			q:    NewSelector[TestModel](db).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewSelector[TestModel](db).From("`test_model_t`").
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			q:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},
		{
			// 非法列
			name:    "invalid column",
			q:       NewSelector[TestModel](db).Where(Not(C("Invalid").GT(18))),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_Get(t *testing.T) {

	mockDb, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := OpenDb(mockDb)
	require.NoError(t, err)

	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	_ = sqlmock.NewRows([]string{"id", "first_name", "last_name"})
	mock.ExpectQuery("SELECT .* WHERE `id` <.*").WillReturnError(ErrNoRows)

	//data
	row := sqlmock.NewRows([]string{"id", "first_name", "age"})
	_ = row.AddRow([]byte("1"), []byte("ceshi"), []byte("18"))
	mock.ExpectQuery("SELECT .* WHERE `id` =.*").WillReturnRows(row)

	//scan error
	row = sqlmock.NewRows([]string{"id", "first_name", "age"})
	_ = row.AddRow([]byte("abc"), []byte("ceshi"), []byte("18"))
	mock.ExpectQuery("SELECT .* WHERE `id` =.*").WillReturnRows(row)
	testCases := []struct {
		name string
		s    *Selector[TestModel]

		wantErr error
		wanRes  *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("xxx").EQ(1)),
			wantErr: errs.NewErrUnknownField("xxx"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").EQ(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").LT(1000)),
			wantErr: ErrNoRows,
		},
		{
			name: "rows data",
			s:    NewSelector[TestModel](db).Where(C("Id").EQ(1)),
			wanRes: &TestModel{
				Id:        1,
				FirstName: "ceshi",
				Age:       18,
				LastName:  nil,
			},
		},
		//{
		//	name:    "scan error",
		//	s:       NewSelector[TestModel](db).Where(C("Id").EQ(1)),
		//	wantErr: fmt.Errorf(`sql: Scan error on column index %d, name %q: %w`, 0, "id", err),
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wanRes, res)
		})
	}
}

// memoryDB 返回一个基于内存的 ORM，它使用的是 sqlite3 内存模式。
func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}
