package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/wwqdrh/easytest/httptest"
	"github.com/wwqdrh/logger"
)

var json = flag.String("json", "", "用于http检查的json文件")

func main() {
	flag.Parse()
	if *json == "" {
		flag.Usage()
		os.Exit(0)
	}

	jsonFile, err := os.Open(*json)
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return
	}
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return
	}

	specInfo, err := httptest.NewBasicSpecInfo(jsonData, nil)
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		return
	}
	if err := specInfo.StartHandle(&testing.T{}); err != nil {
		logger.DefaultLogger.Error(err.Error())
		return
	}
}
