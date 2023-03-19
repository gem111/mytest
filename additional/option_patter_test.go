package additional

import (
	"fmt"
	"testing"
)

func TestNewMyStruct(t *testing.T) {

	testMyStruct := NewMyStruct(1, "234", func(myStruct *MyStruct) {
		myStruct.Name = "222"
		myStruct.Address = "3333"
	})
	fmt.Println(testMyStruct)
	testMyStructv1 := NewMyStruct(1, "234", WithAddress("测试"))
	fmt.Println(testMyStructv1)
}
