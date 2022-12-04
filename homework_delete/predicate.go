package homework_delete

// Expression 不能直接 predicate 里面调用predicate 加一层expression 封起来
type Expression interface {
	expr()
}

const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opIN  = "IN"
	opAND = "AND"
	opOR  = "OR"
	opNOT = "NOT"
)

type Predicate struct {
	left  Expression
	op    string
	right Expression
}

func (Predicate) expr() {

}

func (p Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAND,
		right: right,
	}
}
func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOR,
		right: right,
	}
}
