package examples

import (
	"fmt"
	"testing"

	"github.com/prashantv/gostub"
)

var GlobalVar = "xdcute"

func TestStubGlobal(t *testing.T) {
	stubs := gostub.Stub(&GlobalVar, "stubvar")
	defer stubs.Reset()
	fmt.Println(GlobalVar)

	stubs.Stub(&GlobalVar, "xdcute222")
	fmt.Println(GlobalVar)
}

func TestStubMethod(t *testing.T) {
	var printStr = func(val string) string {
		return val
	}

	// 针对有参数有返回值的
	stubs := gostub.Stub(&printStr, func(val string) string {
		return "hello," + val
	})
	defer stubs.Reset()
	fmt.Println("After stub: ", printStr("hhhhh"))

	var printStr2 = func(val string) string {
		return val
	}
	// StubFunc 第一个参数必须是一个函数变量的指针，该指针指向的必须是一个函数变量，第二个参数为函数 mock 的返回值
	stubs2 := gostub.StubFunc(&printStr2, "ddddd,万生世代")
	defer stubs2.Reset()
	fmt.Println("After stub:", printStr2("lalala"))
}
