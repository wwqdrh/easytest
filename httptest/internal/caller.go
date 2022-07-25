package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type IHTTPCtx interface {
	GetRequest() *http.Request
	GetResponse() *http.Response
	GetEnv(string) interface{}
	SetEnv(string, interface{})
}

type IInstance interface {
	GetAttr(string) interface{}
	ReadAttr() interface{}
}

type ISetInstance interface {
	SetValue(interface{}) interface{}
	ReadAttr() interface{}
}

type DynamicIInstance struct {
	fn     func(string) interface{}
	readFn func() interface{}
}

type DynamicISetInstance struct {
	fn     func(interface{}) interface{}
	readFn func() interface{}
}

func NewDynamicIInstance(fn func(string) interface{}, readFn func() interface{}) IInstance {
	return &DynamicIInstance{
		fn:     fn,
		readFn: readFn,
	}
}

func NewDynamicISetInstance(fn func(interface{}) interface{}, readFn func() interface{}) ISetInstance {
	return &DynamicISetInstance{fn: fn, readFn: readFn}
}

func (d *DynamicIInstance) GetAttr(name string) interface{} {
	return d.fn(name)
}

func (d *DynamicIInstance) ReadAttr() interface{} {
	return d.readFn()
}

func (d *DynamicISetInstance) SetValue(value interface{}) interface{} {
	return d.fn(value)
}

func (d *DynamicISetInstance) ReadAttr() interface{} {
	return d.readFn()
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

	v, err := doCall(ctx, node)
	if err != nil {
		return nil, err
	}

	switch v := v.(type) {
	case ISetInstance:
		return v.ReadAttr(), nil
	case IInstance:
		return v.ReadAttr(), nil
	default:
		return v, nil
	}
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
	} else if node.Type == "callable" && node.Name == "@contain" {
		return CallerFuntion(ctx, node)
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
			func() interface{} { return nil },
		), nil
	case "$body":
		return "$body", nil
	case "$env":
		return wrapEnv(c), nil
	}
	return nil, errors.New("TODO")
}

// 全局函数
func CallerFuntion(c IHTTPCtx, node *SyntaxNode) (interface{}, error) {
	switch node.Name {
	case "@contain":
		if len(node.Params) != 2 {
			return nil, errors.New("@contain必须有两个参数")
		}

		val1, err := doCall(c, node.Params[0])
		if err != nil {
			return nil, err
		}

		val2, err := doCall(c, node.Params[1])
		if err != nil {
			return nil, err
		}
		val2Str, ok := val2.(string)
		if !ok {
			return nil, errors.New("第二个值非字符串")
		}
		val2Str = fmt.Sprintf("%#v", val2Str)
		return strings.Contains(fmt.Sprint(val1), val2Str), nil
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
					return err
				}
				return wrapDict(res)
			case "$str":
				return fmt.Sprint(string(body))
			default:
				return nil
			}
		},
		func() interface{} { return nil },
	)
}

func wrapEnv(ctx IHTTPCtx) IInstance {
	return NewDynamicIInstance(
		func(s string) interface{} {
			return NewDynamicISetInstance(
				func(i interface{}) interface{} {
					ctx.SetEnv(s, i)
					return nil
				},
				func() interface{} {
					return ctx.GetEnv(s)
				},
			)
		},
		func() interface{} { return nil },
	)
}

func wrapDict(data map[string]interface{}) IInstance {
	return NewDynamicIInstance(
		func(s string) interface{} {
			return data[s]
		},
		func() interface{} { return data },
	)
}
