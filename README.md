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

### grpc接口测试

>cinspired by `https://github.com/bojand/ghz`

soon

## 使用手册

使用json文件声明接口测试规则, 省去postman中创建、修改参数、切换选项卡...(想着就繁琐)

*节省时间，时间就是金钱*


### 现在支持的语法列表

- `$res.$body.$json`: 返回响应body的json格式
- `$env.a = 1`: 设置环境变量
- `$env.token = $res.$body.$json.accessToken`: 获取响应的accessToken并作为环境变量
- `@contain($res.$body.$str, "ok")`: 判断响应体中是否包含ok字符串
- `$env.token`: 返回环境变量中的token

```json
[
    {
        "name": "用户注册",
        "url": "http://127.0.0.1:8000/api/user/register",
        "method": "post",
        "body": "{\r\n    \"name\":\"ving\",\r\n    \"gender\":1,\r\n    \"mobile\": \"15212230311\",\r\n    \"password\": \"123456\"\r\n}\r\n",
        "content-type": "application/json",
        "expect": [
            "@contain($res.$body.$str, \"ok\")"
        ]
    },
    {
        "name": "用户登录",
        "url": "http://127.0.0.1:8000/api/user/login",
        "method": "post",
        "body": "{\"name\":\"ving\",\"gender\":1,\"mobile\": \"15212230311\",\"password\": \"123456\"}",
        "content-type": "application/json",
        "expect": [
            "@contain($res.$body.$str, \"ok\")"
        ],
        "event": [
            "$env.token = $res.$body.$json.accessToken"
        ]
    },
    {
        "name": "用户信息",
        "url": "http://127.0.0.1:8000/api/user/userinfo",
        "method": "post",
        "body": "{\r\n    \"name\":\"ving\",\r\n    \"gender\":1,\r\n    \"mobile\": \"15212230311\",\r\n    \"password\": \"123456\"\r\n}\r\n",
        "content-type": "application/json",
        "header": [
            "Authorization: bearer {{ token }}"
        ],
        "expect": [
            "@contain($res.$body.$str, \"ok\")"
        ]
    }
]
```

### 集成在单元测试

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

	specInfo, err := NewBasicParserSpecInfo(postmanJsonData, func(item *BasicItem) {
		item.Url = ts.URL
	})
	require.Nil(t, err)

	specInfo.StartHandle(t)
}
```

### 二进制工具


```shell
go install github.com/wwqdrh/easytest/cmd/etcli@latest
```

- json: 指定需要检查的文件(格式与上面的一样)
- check: 测试当前版本功能是否正常