package additional

type MyStruct struct {
	Id      uint64
	Name    string
	Address string
}

type MyStructOption func(myStruct *MyStruct)

func NewMyStruct(id uint64, name string, opts ...MyStructOption) *MyStruct {

	res := &MyStruct{
		Id:   id,
		Name: name,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

func WithAddress(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.Address = address
	}
}
