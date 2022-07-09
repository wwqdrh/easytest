package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	easyhttp "github.com/wwqdrh/easytest/httptest"

	"net/http/httptest"

	"github.com/wwqdrh/logger"
)

var (
	jsonfile = flag.String("json", "api.json", "用于http检查的json文件")
	check    = flag.Bool("check", false, "检查当前版本功能是否正常")
)

var (
	checkUrl string
)

//go:embed api.json
var testapi []byte

func main() {
	flag.Parse()
	if *jsonfile == "" {
		flag.Usage()
		os.Exit(0)
	}

	if *check {
		checkRun()
	} else {
		basicRun()
	}
}

func checkRun() {
	// mock 实现
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/user/userinfo") {
			if !strings.Contains(r.Header.Get("Authorization"), "123456") {
				w.WriteHeader(500)
				return
			}
		}
		body, _ := json.Marshal(map[string]interface{}{
			"msg":         "ok",
			"accessToken": "123456",
		})
		w.Write(body)
	}))
	defer ts.Close()

	checkUrl = ts.URL

	basicRun()
}

func getJsonStr() ([]byte, error) {
	if *check {
		return testapi, nil
	}

	jsonFile, err := os.Open(*jsonfile)
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return nil, err
	}
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return nil, err
	}
	return jsonData, nil
}

func basicRun() {
	jsonData, err := getJsonStr()
	if err != nil {
		panic(err)
	}

	var specInfo *easyhttp.BasicParserSpecInfo
	if checkUrl != "" {
		specInfo, err = easyhttp.NewBasicParserSpecInfo(jsonData, func(item *easyhttp.BasicItem) {
			item.Url = checkUrl + getPath(item.Url)
		})
	} else {
		specInfo, err = easyhttp.NewBasicParserSpecInfo(jsonData, nil)
	}

	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return
	}

	if err := specInfo.StartHandle(&testing.T{}); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}

func getPath(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return "404"
	}
	return u.Path
}
