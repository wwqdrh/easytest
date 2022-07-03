package gomock

//go:generate mockgen -package gomock -source foo.go -destination=foo_mock.go

import (
	"log"
)

//定义了一个订单接口，有一个获取名称的方法
type OrderDBI interface {
	GetName(orderid int) string
}

//定义结构体
type OrderInfo struct {
	orderid int
}

//实现接口的方法
func (order OrderInfo) GetName(orderid int) string {
	log.Println("原本应该连接数据库去取名称")
	return "xdcute"
}
