package httptest

import (
	"fmt"
	"strings"
)

////////////////////
// 处理json转req请求时的各种操作
// 1、$env，设置环境变量
// 2、{{...}}. 读取环境变量
// 3、$contains, 响应中包含某个字符串
// 4、$json, 获取相应并以json格式解析
//
// 词法分析
// 语法分析
////////////////////

// 判断c响应是否满足expect
func HandleExpect(c *HttpContext, expect []string) bool {
	for _, item := range expect {
		if strings.Index(item, "$contains") == 0 {
			value := item[len("$contains")+1 : len(item)-1]
			v := strings.Split(value, ",")
			parts := v[:0]
			for _, item := range v {
				parts = append(parts, strings.TrimSpace(item))
			}

			if !HandleContains(c, parts) {
				return true
			}
		} else if strings.Index(item, "$status") == 0 {
			statuscode := item[len("$status")+1 : len(item)-1]
			return fmt.Sprint(c.responseStatus) == statuscode
		}
	}
	return true
}

func HandleEvent(c *HttpContext, event []string) bool {
	for _, item := range event {
		if strings.Index(item, "$env") == 0 {
			pairs := strings.Split(item, "=")
			if len(pairs) != 2 {
				return false
			}

			left := strings.Split(pairs[0], ".")[1]
			if strings.Index(pairs[1], "$json") == 0 {
				right := c.responseJson[strings.Split(pairs[1], ".")[1]]
				if right == "" {
					return false
				}
				c.Setenv(left, right)
			} else {
				return false
			}
		}
	}
	return true
}

// $contains
func HandleContains(c *HttpContext, args []string) bool {
	if len(args) != 2 {
		return false
	}

	data := ""
	if args[0] == "$res" {
		data = c.responseData
	}
	return strings.Contains(data, args[1])
}
