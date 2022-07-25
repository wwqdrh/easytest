//go:generate mockgen -package internal -source caller.go -destination=caller_mock.go

package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CallerTestSuite struct {
	suite.Suite
	mockServer *httptest.Server
}

func TestCaller(t *testing.T) {
	suite.Run(t, new(CallerTestSuite))
}

// 变成字符串后 => "{\"accessToken\":\"12345\",\"msg\":\"ok\",\"msgwithline\":\"\\\"ok\\\"\",\"msgzh\""
func (suite *CallerTestSuite) SetupTest() {
	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(map[string]interface{}{
			"msg":         "ok",
			"accessToken": "12345",
			"msgzh":       "请求成功",
			"msgwithline": `"ok"`,
		})
		w.Write(body)
	}))

}

func (suite *CallerTestSuite) TearDownSuite() {
	suite.mockServer.Close()
}

func (suite *CallerTestSuite) TestCaller1() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockReuqest, err := http.NewRequest("get", suite.mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

	// newResponse := func() *http.Response {
	// 	// read the response body to a variable
	// 	bodyBytes, _ := ioutil.ReadAll(mockResponse.Body)
	// 	newResponse := *mockResponse
	// 	//reset the response body to the original unread state
	// 	newResponse.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// 	mockResponse.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// 	return &newResponse
	// }

	mock := NewMockIHTTPCtx(ctrl)
	mock.EXPECT().GetRequest().AnyTimes().Return(mockReuqest)
	mock.EXPECT().GetResponse().AnyTimes().Return(mockResponse)
	mock.EXPECT().SetEnv(gomock.Eq("a"), gomock.Eq(1)).AnyTimes()
	mock.EXPECT().GetEnv(gomock.Eq("a")).AnyTimes().Return(1)

	val, err := DoCaller(mock, "$res.$body.$json")
	require.Nil(t, err)
	fmt.Println(val)

	val, err = DoCaller(mock, "$env.a = 1")
	require.Nil(t, err)
	fmt.Println(val)
	require.Equal(t, mock.GetEnv("a"), 1)
}

func (suite *CallerTestSuite) TestCaller2() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReuqest, err := http.NewRequest("get", suite.mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

	mock := NewMockIHTTPCtx(ctrl)
	mock.EXPECT().GetRequest().AnyTimes().Return(mockReuqest)
	mock.EXPECT().GetResponse().AnyTimes().Return(mockResponse)
	mock.EXPECT().SetEnv(gomock.Eq("token"), gomock.Eq("12345")).AnyTimes()
	mock.EXPECT().GetEnv(gomock.Eq("token")).AnyTimes().Return("12345")

	val, err := DoCaller(mock, "$env.token = $res.$body.$json.accessToken")
	require.Nil(t, err)
	fmt.Println(val)
}

func (suite *CallerTestSuite) TestCaller3() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReuqest, err := http.NewRequest("get", suite.mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

	mock := NewMockIHTTPCtx(ctrl)
	mock.EXPECT().GetRequest().AnyTimes().Return(mockReuqest)
	mock.EXPECT().GetResponse().AnyTimes().Return(mockResponse)
	mock.EXPECT().SetEnv(gomock.Eq("token"), gomock.Eq("12345")).AnyTimes()
	mock.EXPECT().GetEnv(gomock.Eq("token")).AnyTimes().Return("12345")

	val, err := DoCaller(mock, `@contain($res.$body.$str, "ok")`)
	require.Nil(t, err)
	fmt.Println(val)
}

func (suite *CallerTestSuite) TestCaller4() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReuqest, err := http.NewRequest("get", suite.mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

	mock := NewMockIHTTPCtx(ctrl)
	mock.EXPECT().GetRequest().AnyTimes().Return(mockReuqest)
	mock.EXPECT().GetResponse().AnyTimes().Return(mockResponse)
	mock.EXPECT().SetEnv(gomock.Eq("token"), gomock.Eq("12345")).AnyTimes()
	mock.EXPECT().GetEnv(gomock.Eq("token")).AnyTimes().Return("12345")

	val, err := DoCaller(mock, `$env.token`)
	require.Nil(t, err)
	fmt.Println(val)
}

func (suite *CallerTestSuite) TestContainerWithUtf8() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReuqest, err := http.NewRequest("get", suite.mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

	mock := NewMockIHTTPCtx(ctrl)
	mock.EXPECT().GetRequest().AnyTimes().Return(mockReuqest)
	mock.EXPECT().GetResponse().AnyTimes().Return(mockResponse)
	mock.EXPECT().SetEnv(gomock.Eq("token"), gomock.Eq("12345")).AnyTimes()
	mock.EXPECT().GetEnv(gomock.Eq("token")).AnyTimes().Return("12345")

	val, err := DoCaller(mock, `@contain($res.$body.$str, "请求成功")`)
	require.Nil(t, err)
	fmt.Println(val)
}

func (suite *CallerTestSuite) TestContainerWithline() {
	t := suite.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReuqest, err := http.NewRequest("get", suite.mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

	mock := NewMockIHTTPCtx(ctrl)
	mock.EXPECT().GetRequest().AnyTimes().Return(mockReuqest)
	mock.EXPECT().GetResponse().AnyTimes().Return(mockResponse)
	mock.EXPECT().SetEnv(gomock.Eq("token"), gomock.Eq("12345")).AnyTimes()
	mock.EXPECT().GetEnv(gomock.Eq("token")).AnyTimes().Return("12345")

	val, err := DoCaller(mock, `@contain($res.$body.$str, "\"ok\"")`)
	require.Nil(t, err)
	fmt.Println(val)
}
