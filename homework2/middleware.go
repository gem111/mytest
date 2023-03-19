package web

//Middleware 函数式的责任
type Middleware func(next HandleFunc) HandleFunc
