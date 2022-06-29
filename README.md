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

### http接口测试

- 🗂 http接口的集成测试，添加环境变量机制
- 📦 postman接口文件支持

### 单元测试

- soon

## 使用手册

使用json文件声明接口测试规则, 省去postman中创建、修改参数、切换选项卡...(想着就繁琐)

*节省时间，时间就是金钱*


```json
[
    {
        "name": "用户注册",
        "url": "http://127.0.0.1:8000/api/user/register",
        "method": "post",
        "body": "{\r\n    \"name\":\"ving\",\r\n    \"gender\":1,\r\n    \"mobile\": \"15212230311\",\r\n    \"password\": \"123456\"\r\n}\r\n",
        "content-type": "application/json",
        "expect": [
            "$contains($res, ok)"
        ]
    },
    {
        "name": "用户登录",
        "url": "http://127.0.0.1:8000/api/user/login",
        "method": "post",
        "body": "{\r\n    \"name\":\"ving\",\r\n    \"gender\":1,\r\n    \"mobile\": \"15212230311\",\r\n    \"password\": \"123456\"\r\n}\r\n",
        "content-type": "application/json",
        "expect": [
            "$contains($res, ok)"
        ],
        "event": [
            "$env.token=$json.accessToken"
        ]
    },
    {
        "name": "用户信息",
        "url": "http://127.0.0.1:8000/api/user/userinfo",
        "method": "post",
        "body": "{\r\n    \"name\":\"ving\",\r\n    \"gender\":1,\r\n    \"mobile\": \"15212230311\",\r\n    \"password\": \"123456\"\r\n}\r\n",
        "content-type": "application/json",
        "header": [
            "Authorization: bearer {{token}}"
        ],
        "expect": [
            "$status(200)",
            "$contains($res, ok)"
        ]
    }
]
```

```go

import (
	"encoding/json"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	easyhttptest "github.com/wwqdrh/easytest/httptest"
)

func TestHTTPByJSON() {
    postmanJsonFile, err := os.Open(".json")
	require.Nil(t, err)

	postmanJsonData, err := ioutil.ReadAll(postmanJsonFile)
	require.Nil(t, err)

	specInfo, err := NewBasicSpecInfo(postmanJsonData, func(item *BasicItem) {
		item.Url = ts.URL
	})
	require.Nil(t, err)

	specInfo.StartHandle(t)
}
```