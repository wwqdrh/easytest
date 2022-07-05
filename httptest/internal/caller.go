package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type IHTTPCtx interface {
	GetRequest() *http.Request
	GetResponse() *http.Response
	GetEnv(string) interface{}
	SetEnv(string, interface{})
}

type IInstance interface {
	GetAttr(string) interface{}
}

type ISetInstance interface {
	SetValue(interface{}) interface{}
}

type DynamicIInstance struct {
	fn func(string) interface{}
}

type DynamicISetInstance struct {
	fn func(interface{}) interface{}
}

func NewDynamicIInstance(fn func(string) interface{}) IInstance {
	return &DynamicIInstance{
		fn: fn,
	}
}

func NewDynamicISetInstance(fn func(interface{}) interface{}) ISetInstance {
	return &DynamicISetInstance{fn: fn}
}

func (d *DynamicIInstance) GetAttr(name string) interface{} {
	return d.fn(name)
}

func (d *DynamicISetInstance) SetValue(value interface{}) interface{} {
	return d.fn(value)
}

// 参数也是*SyntaxNode结构
// 存储符号表以及执行函数定义
type Caller func(IHTTPCtx, []interface{}) (interface{}, error)

func DoCaller(ctx IHTTPCtx, source string) (interface{}, error) {
	p := NewSimpleParser(NewLexer(source))
	node, err := p.Parse()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return doCall(ctx, node)
}

func doCall(ctx IHTTPCtx, node *SyntaxNode) (interface{}, error) {
	if node.Type == "expression" && node.Name == "." {
		return CallerDot(ctx, node.Params)
	} else if node.Type == "expression" && node.Name == "=" {
		return CallAssign(ctx, node.Params)
	} else if node.Type == "global" {
		return CallerGlobal(ctx, node)
	} else if node.Type == "attr" {
		return node.Name, nil
	} else if node.Type == "literial" || node.Type == "variable" {
		return node.Value, nil
	}
	return nil, errors.New("解析失败")
}

// 实现了.取值符的必须实现了IInstance接口或者本身是map数据类型
func CallerDot(c IHTTPCtx, params []*SyntaxNode) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("ast error, 参数只能为两个")
	}

	ins, attr := params[0], params[1]

	insValI, err := doCall(c, ins)
	if err != nil {
		return nil, err
	}
	insVal, ok := insValI.(IInstance)
	if !ok {
		return nil, errors.New(". 左边不是IInstance类型，没有GetAttr方法")
	}

	attrVal, err := doCall(c, attr)
	if err != nil {
		return nil, err
	}

	return insVal.GetAttr(fmt.Sprint(attrVal)), nil
}

func CallAssign(c IHTTPCtx, params []*SyntaxNode) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("ast error, 参数只能为两个")
	}

	ins, attr := params[0], params[1]

	insValI, err := doCall(c, ins)
	if err != nil {
		return nil, err
	}
	insVal, ok := insValI.(ISetInstance)
	if !ok {
		return nil, errors.New(". 左边不是IInstance类型，没有GetAttr方法")
	}

	attrVal, err := doCall(c, attr)
	if err != nil {
		return nil, err
	}

	return insVal.SetValue(attrVal), nil
}

// 获取当前http请求上下文中的response
func CallerGlobal(c IHTTPCtx, node *SyntaxNode) (interface{}, error) {
	switch node.Name {
	case "$res":
		return NewDynamicIInstance(
			func(s string) interface{} {
				switch s {
				case "$body":
					response := c.GetResponse()
					body, err := ioutil.ReadAll(response.Body)
					if err != nil {
						return nil
					}
					return wrapResBody(body)
				default:
					return nil
				}
			},
		), nil
	case "$body":
		return "$body", nil
	case "$env":
		return wrapEnv(c), nil
	}
	return nil, errors.New("TODO")
}

func wrapResBody(body []byte) IInstance {
	return NewDynamicIInstance(
		func(s string) interface{} {
			switch s {
			case "$json":
				res := map[string]interface{}{}
				if err := json.Unmarshal(body, &res); err != nil {
					return nil
				}
				return res
			default:
				return nil
			}
		},
	)
}

func wrapEnv(ctx IHTTPCtx) IInstance {
	return NewDynamicIInstance(
		func(s string) interface{} {
			return NewDynamicISetInstance(func(i interface{}) interface{} {
				ctx.SetEnv(s, i)
				return nil
			})
		},
	)
}
