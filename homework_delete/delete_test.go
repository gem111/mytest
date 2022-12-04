package homework_delete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleter_Build(t *testing.T) {

	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no where",
			builder: (&Deleter[TestModel]{}).From("`test_model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "no form",
			builder: &Deleter[TestModel]{},
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: (&Deleter[TestModel]{}).Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `Id` = ?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `Id` = ?;",
				Args: []any{16},
			},
		},
		{
			name:    "where and",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").EQ(16).And(C("Name").EQ("nihao"))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`Id` = ?) AND (`Name` = ?);",
				Args: []any{16, "nihao"},
			},
		},
		{
			name:    "where or",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").EQ(16).Or(C("Name").EQ("nihao"))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`Id` = ?) OR (`Name` = ?);",
				Args: []any{16, "nihao"},
			},
		},
		{
			name:    "where in",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").In([]any{2, 3, 4, 5, 7})),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `Id` IN ?;",
				Args: []any{[]any{2, 3, 4, 5, 7}},
			},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
