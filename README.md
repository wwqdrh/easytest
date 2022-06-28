<p align='center'>
  <pre style="float:left;">
   ('-.      ('-.       .-')                  .-') _       ('-.     .-')     .-') _    
 _(  OO)    ( OO ).-.  ( OO ).               (  OO) )    _(  OO)   ( OO ).  (  OO) )   
(,------.   / . --. / (_)---\_)   ,--.   ,--./     '._  (,------. (_)---\_) /     '._  
 |  .---'   | \-.  \  /    _ |     \  `.'  / |'--...__)  |  .---' /    _ |  |'--...__) 
 |  |     .-'-'  |  | \  :` `.   .-')     /  '--.  .--'  |  |     \  :` `.  '--.  .--' 
(|  '--.   \| |_.'  |  '..`''.) (OO  \   /      |  |    (|  '--.   '..`''.)    |  |    
 |  .--'    |  .-.  | .-._)   \  |   /  /\_     |  |     |  .--'  .-._)   \    |  |    
 |  `---.   |  | |  | \       /  `-./  /.__)    |  |     |  `---. \       /    |  |    
 `------'   `--' `--'  `-----'     `--'         `--'     `------'  `-----'     `--'    
  </pre>
</p>

<p align='center'>
方便地<sup><em>EasyTest</em></sup>测试模板
<br> 
</p>

<br>

## 背景


封装常用的测试代码

再也不用担心测试代码编写麻烦了


## 特性

### 集成测试

- 🗂 http接口的集成测试，添加环境变量机制
- 📦 postman接口文件支持

### 单元测试

- soon

## 使用手册

```go

import (
	"encoding/json"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	easyhttptest "github.com/wwqdrh/easytest/httptest"
)

func main() {
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
```