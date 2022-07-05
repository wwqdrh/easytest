package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -package internal -source caller.go -destination=caller_mock.go

func TestCaller(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(map[string]interface{}{
			"msg": "ok",
		})
		w.Write(body)
	}))
	defer mockServer.Close()

	mockReuqest, err := http.NewRequest("get", mockServer.URL, nil)
	require.Nil(t, err)
	mockResponse, err := http.DefaultClient.Do(mockReuqest)
	require.Nil(t, err)

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
