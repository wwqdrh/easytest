caller执行的时候需要httpctx的上下文

## 关键字

- $env: 上下文的环境变量
- $res: 响应数据
- $req: 请求报文
- $raw: 原始数据
- $json: json格式数据
- $header: 报文头
- $body: 报文体

- $in: 全局函数，判断字符串是否包含指定的字符串

## 运算符

- .: 取对象的值, 后面接将数据格式如何转换, 存在关键字或者普通变量, 普通变量时默认将前面的数据转为json后处理
- =: 赋值语句

## TODO

3、占位符的实现

例如为req header添加环境变量时

`Authorization: bearer {{ $env.token }}`