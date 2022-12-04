package homework_delete

type Column struct {
	name string
}

//标记成 predicate
func (c Column) expr() {}
func (v value) expr()  {}

type value struct {
	val any
}

func valueOf(val any) Expression {
	return value{val: val}
}

func C(name string) Column {
	return Column{name: name}
}

func (c Column) Not(ary any) Predicate {
	return Predicate{
		left:  c,
		op:    opNOT,
		right: valueOf(ary),
	}
}
func (c Column) EQ(ary any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: valueOf(ary),
	}
}
func (c Column) In(ary []any) Predicate {
	return Predicate{
		left:  c,
		op:    opIN,
		right: valueOf(ary),
	}
}
