package http_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	easyhttptest "github.com/wwqdrh/easytest/httptest"
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

	ctx := easyhttptest.NewHttpContext()
	ctx.Do(t, "user login", &easyhttptest.HandleOption{
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

	ctx.Do(t, "user info", &easyhttptest.HandleOption{
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

func TestHTTPFromJson(t *testing.T) {
	// mock 实现
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	postmanJsonFile, err := os.Open("./testdata/gomall.postman_collection.json")
	require.Nil(t, err)

	postmanJsonData, err := ioutil.ReadAll(postmanJsonFile)
	require.Nil(t, err)

	specInfo, err := easyhttptest.NewPostmanSpecInfo(postmanJsonData, func(item *easyhttptest.PostmanItem) {
		item.Request.Url.Host = strings.Split(ts.URL, ".")
	})
	require.Nil(t, err)

	specInfo.StartHandle(t)
}
