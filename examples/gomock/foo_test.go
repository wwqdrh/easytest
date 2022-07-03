package gomock

import (
	"log"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestGetName(t *testing.T) {
	//新建一个mockController
	ctrl := gomock.NewController(t)
	// 断言 DB.GetName() 方法是否被调用
	defer ctrl.Finish()

	//mock接口
	mock := NewMockOrderDBI(ctrl)
	//模拟传入值与预期的返回值
	mock.EXPECT().GetName(gomock.Eq(1225)).Return("xdcutecute")

	//前面定义了传入值与返回值
	//在这里
	if v := mock.GetName(1225); v != "xdcutecute" {
		t.Fatal("expected xdcute, but got", v)
	} else {
		log.Println("通过mock取到的name：", v)
	}
}
