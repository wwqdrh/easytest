package httptest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

func TestHTTPFromPostmanJson(t *testing.T) {
	// mock 实现
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	postmanJsonFile, err := os.Open("./testdata/gomall.postman_collection.json")
	require.Nil(t, err)

	postmanJsonData, err := ioutil.ReadAll(postmanJsonFile)
	require.Nil(t, err)

	specInfo, err := NewPostmanSpecInfo(postmanJsonData, func(item *PostmanItem) {
		item.Request.Url.Host = strings.Split(ts.URL, ".")
	})
	require.Nil(t, err)

	specInfo.StartHandle(t)
}

func TestHTTPFromBasicJson(t *testing.T) {
	// mock 实现
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/userinfo" {
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(500)
				return
			}
		}

		body, _ := json.Marshal(map[string]interface{}{
			"msg":         "ok",
			"accessToken": "132",
		})
		w.Write(body)
	}))
	defer ts.Close()

	postmanJsonFile, err := os.Open("./testdata/gomall.basic_collection.json")
	require.Nil(t, err)

	postmanJsonData, err := ioutil.ReadAll(postmanJsonFile)
	require.Nil(t, err)

	specInfo, err := NewBasicSpecInfo(postmanJsonData, func(item *BasicItem) {
		item.Url = ts.URL + getPath(item.Url)
	})
	require.Nil(t, err)

	specInfo.StartHandle(t)
}

func TestHTTPFromBasicParserJson(t *testing.T) {
	// mock 实现
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/userinfo" {
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(500)
				return
			}
		}

		body, _ := json.Marshal(map[string]interface{}{
			"msg":         "ok",
			"accessToken": "132",
		})
		w.Write(body)
	}))
	defer ts.Close()

	postmanJsonFile, err := os.Open("./testdata/gomall.parser_collection.json")
	require.Nil(t, err)

	postmanJsonData, err := ioutil.ReadAll(postmanJsonFile)
	require.Nil(t, err)

	specInfo, err := NewBasicParserSpecInfo(postmanJsonData, func(item *BasicItem) {
		item.Url = ts.URL + getPath(item.Url)
	})
	require.Nil(t, err)

	specInfo.StartHandle(t)
}

func getPath(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return "404"
	}
	return u.Path
}
