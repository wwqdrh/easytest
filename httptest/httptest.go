package httptest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wwqdrh/logger"
)

var envReg *regexp.Regexp

func init() {
	var err error
	envReg, err = regexp.Compile("{{(.*?)}}")
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
}

type HttpContext struct {
	enviroment map[string]interface{}

	responseStatus int
	responseData   string
	responseJson   map[string]interface{}
}

type HandleOption struct {
	Url         string
	Method      string
	ContentType string
	Header      map[string]string
	Body        io.Reader
	Handle      func(resp *http.Response) error

	Expect []string
	Event  []string
}

func NewHttpContext() *HttpContext {
	return &HttpContext{
		enviroment: map[string]interface{}{},
	}
}

func (c *HttpContext) Do(t *testing.T, title string, option *HandleOption) {
	req, err := http.NewRequest(option.Method, option.Url, option.Body)
	require.Nil(t, err, title)
	for key, value := range c.ReqHeader(option.Header) {
		req.Header.Add(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err, title)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.DefaultLogger.Warn(err.Error())
		return
	}
	bodyData := string(body)
	c.responseData = bodyData
	c.responseStatus = resp.StatusCode

	// 获取json
	jsonData := map[string]interface{}{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		// logger.DefaultLogger.Warn(err.Error())
		return
	}
	c.responseJson = jsonData

	if option.Handle != nil {
		option.Handle(resp)
		// require.Nil(t, err, title)
	}
	// 处理response expect
	assert.True(t, HandleExpect(c, option.Expect))
	// 处理event
	assert.True(t, HandleEvent(c, option.Event))
}

func (c *HttpContext) Setenv(key string, value interface{}) {
	c.enviroment[key] = value
}

func (c *HttpContext) ReqHeader(header map[string]string) map[string]string {
	res := map[string]string{}
	for key, value := range header {
		res[key] = value
		for _, v := range envReg.FindAllStringSubmatch(value, -1) {
			if len(v) == 2 {
				envVal := c.enviroment[v[1]]
				if envVal != "" {
					res[key] = strings.Replace(res[key], v[0], fmt.Sprint(envVal), 1)
				}
			} else {
				res[key] = value
			}
		}
	}
	return res
}

func (c *HttpContext) Json(resp *http.Response) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}
