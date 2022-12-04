package homework_delete

import (
	"fmt"
	"reflect"
	"strings"
)

type Deleter[T any] struct {
	table string
	where []Predicate
	args  []any
	sql   strings.Builder
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sql.WriteString("DELETE FROM ")
	if d.table == "" {
		var t T
		d.sql.WriteByte('`')
		d.sql.WriteString(UnderscoreName(reflect.TypeOf(t).Name()))
		d.sql.WriteByte('`')
	} else {
		d.sql.WriteString(d.table)
	}

	if len(d.where) > 0 {
		d.sql.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}

	}
	d.sql.WriteByte(';')
	return &Query{
		SQL:  d.sql.String(),
		Args: d.args,
	}, nil
}
func (d *Deleter[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}

	switch exp := e.(type) {
	case Column:
		d.sql.WriteByte('`')
		d.sql.WriteString(exp.name)
		d.sql.WriteByte('`')
	case value:
		d.sql.WriteByte('?')
		d.args = append(d.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			d.sql.WriteByte('(')
		}
		if err := d.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			d.sql.WriteByte(')')
		}
		d.sql.WriteByte(' ')
		d.sql.WriteString(exp.op)
		d.sql.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			d.sql.WriteByte('(')
		}
		if err := d.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			d.sql.WriteByte(')')
		}
	default:
		return fmt.Errorf("orm : 不支持的表达式 %v", exp)
	}
	return nil
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = predicates
	return d
}
