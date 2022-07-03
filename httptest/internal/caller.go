package internal

import "errors"

type IHTTPCtx interface {
	GetReqHeader() interface{} // 获取请求体的header
	GetResHeader() interface{} // 获取响应的header
	GetResBody() interface{}   // 获取响应的raw
	GetRes() interface{}       // 获取http.response
}

type IInstance interface {
	GetAttr(string) interface{}
}

// 参数也是*SyntaxNode结构
// 存储符号表以及执行函数定义
type Caller func(IHTTPCtx, []interface{}) (interface{}, error)

// 实现了.取值符的必须实现了IInstance接口或者本身是map数据类型
func CallerDot(c IHTTPCtx, params []*SyntaxNode) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("ast error, 参数只能为两个")
	}

	ins, attr := params[0], params[1]
	var insVal IInstance
	switch ins.Type {
	case "global":
		v, err := CallerGlobal(c, []*SyntaxNode{ins})
		if err != nil {
			return nil, err
		}
		if v, ok := v.(IInstance); !ok {
			return nil, errors.New("不是能获取属性的类型")
		} else {
			insVal = v
		}
	}

	switch attr.Type {
	case "global":
		return nil, errors.New("ast errror, dot的第二个参数只能为identifer")
	}

	return insVal.GetAttr(attr.Name), nil
}

// 获取当前http请求上下文中的response
func CallerGlobal(c IHTTPCtx, params []*SyntaxNode) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("ast error, 参数只能为一个")
	}

	switch params[0].Name {
	case "$res":
		return c.GetRes(), nil
	case "$reqheader":
		return c.GetReqHeader(), nil
	}
	return nil, errors.New("TODO")
}
