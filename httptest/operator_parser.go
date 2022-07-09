package httptest

// parser版的operaotr

import (
	"net/http"

	"github.com/wwqdrh/easytest/httptest/internal"
)

type HTTPCtx struct {
	ctx *HttpContext
}

func NewIHTTPCtx(ctx *HttpContext) internal.IHTTPCtx {
	return &HTTPCtx{
		ctx: ctx,
	}
}

func (c *HTTPCtx) GetRequest() *http.Request {
	return c.ctx.request
}
func (c *HTTPCtx) GetResponse() *http.Response {
	return c.ctx.CopyResponse(c.ctx.response)
}
func (c *HTTPCtx) GetEnv(key string) interface{} {
	return c.ctx.enviroment[key]
}
func (c *HTTPCtx) SetEnv(key string, val interface{}) {
	c.ctx.enviroment[key] = val
}

// 判断c响应是否满足expect
func ParserHandleExpect(c *HttpContext, expect []string) bool {
	curCtx := NewIHTTPCtx(c)

	for _, item := range expect {
		val, err := internal.DoCaller(curCtx, item)
		if err != nil {
			return false
		}

		switch val := val.(type) {
		case bool:
			return val
		}
	}
	return true
}

func ParserHandleEvent(c *HttpContext, event []string) bool {
	curCtx := NewIHTTPCtx(c)

	for _, item := range event {
		val, err := internal.DoCaller(curCtx, item)
		if err != nil {
			return false
		}

		switch val := val.(type) {
		case bool:
			return val
		}
	}
	return true
}
