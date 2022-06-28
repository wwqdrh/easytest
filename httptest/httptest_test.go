package httptest

import (
	"fmt"
	"testing"
)

func TestReqHeader(t *testing.T) {
	ctx := NewHttpContext()
	ctx.Setenv("token", "123456")
	ctx.Setenv("token2", "654321")
	fmt.Println(ctx.ReqHeader(map[string]string{"withtoken": "token: {{token}}; token2: {{token2}}"}))
}
