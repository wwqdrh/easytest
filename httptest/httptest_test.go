package httptest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoHandle(t *testing.T) {
	// mock 实现
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			resp := map[string]interface{}{
				"accessToken": "123456",
			}
			data, _ := json.Marshal(resp)
			w.Write(data)
		case "/info":
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(500)
			} else {
				w.WriteHeader(200)
			}
		default:
			w.WriteHeader(500)
		}
	}))
	defer ts.Close()

	ctx := NewHttpContext()
	ctx.Do(t, "user login", &HandleOption{
		Method: "POST",
		Url:    ts.URL + "/login",
		Handle: func(resp *http.Response) error {
			jsonData, err := ctx.Json(resp)
			if err != nil {
				return err
			}
			ctx.Setenv("token", jsonData["accessToken"])
			return nil
		},
	})

	ctx.Do(t, "user info", &HandleOption{
		Method: "GET",
		Url:    ts.URL + "/info",
		Header: map[string]string{
			"Authorization": "bearer {{token}}",
		},
		Handle: func(resp *http.Response) error {
			assert.Equal(t, resp.StatusCode, 200)
			return nil
		},
	})
}

func TestReqHeader(t *testing.T) {
	ctx := NewHttpContext()
	ctx.Setenv("token", "123456")
	ctx.Setenv("token2", "654321")
	fmt.Println(ctx.ReqHeader(map[string]string{"withtoken": "token: {{token}}; token2: {{token2}}"}))
}
