package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// 词法解析
// 包含关键字、字面量、运算符
type Tag uint

const (
	ENV Tag = iota + 256
	RES
	REQ
	RAW
	BODY
	JSON
	HEADER

	// 全局函数
	IN

	// 操作符
	EQ
	ASSIGN_OPERATOR
	DOT

	// 字面量
	NUM        // 数字
	REAL       // 浮点数
	INDENTIFER // 变量

	// 其他标识符
	EOF
	ERROR
)

var tokenMap = map[Tag]string{
	ENV:             "$env",
	RES:             "$res",
	REQ:             "$req",
	RAW:             "$raw",
	BODY:            "$body",
	JSON:            "$json",
	HEADER:          "$header",
	IN:              "$in",
	EQ:              "==",
	ASSIGN_OPERATOR: "=",
	DOT:             ".",
	NUM:             "num",
	REAL:            "real",
	INDENTIFER:      "indentifer",
	EOF:             "EOF",
	ERROR:           "syntax error",
}

var keyWord = []KeyWord{
	NewKeyWord(ENV),
	NewKeyWord(RES),
	NewKeyWord(REQ),
	NewKeyWord(RAW),
	NewKeyWord(JSON),
	NewKeyWord(BODY),
	NewKeyWord(HEADER),
	NewKeyWord(IN),
}

// token字符分类
type Token struct {
	Tag Tag
}

func NewToken(tag Tag) Token {
	return Token{
		Tag: tag,
	}
}

func (t *Token) String() string {
	return tokenMap[t.Tag]
}

// 关键字
type KeyWord struct {
	lexeme string
	Tag    Token
}

func NewKeyWord(tag Tag) KeyWord {
	return KeyWord{
		lexeme: tokenMap[tag],
		Tag:    NewToken(tag),
	}
}

func (w *KeyWord) String() string {
	return w.lexeme
}

// 词法解析器
type Lexer struct {
	Lexeme      string
	lexemeStack []string
	peek        byte             // 读入的字符
	line        int              // 当前字符串处于第几行
	reader      *bufio.Reader    // 用于读取字节流
	keyWords    map[string]Token // 存储关键字
}

func NewLexer(source string) Lexer {
	str := strings.NewReader(source)
	sourceReader := bufio.NewReaderSize(str, len(source))
	lexer := Lexer{
		line:     1,
		reader:   sourceReader,
		keyWords: map[string]Token{},
	}
	lexer.reserve() // 保留所有关键字
	return lexer
}

func (l *Lexer) reserve() {
	for _, keyword := range keyWord {
		l.keyWords[keyword.String()] = keyword.Tag
	}
}

func (l *Lexer) ReverseScan() {
	backLen := len(l.Lexeme)
	for i := 0; i < backLen; i++ {
		l.reader.UnreadByte()
	}

	l.lexemeStack = l.lexemeStack[:len(l.lexemeStack)-1]
	l.Lexeme = l.lexemeStack[len(l.lexemeStack)-1]
}

func (l *Lexer) Readch() error {
	char, err := l.reader.ReadByte()
	l.peek = char
	return err
}

func (l *Lexer) UnRead() error {
	return l.reader.UnreadByte()
}

func (l *Lexer) ReadCharacter(c byte) (bool, error) {
	chars, err := l.reader.Peek(1) // 不会从缓冲区删除
	if err != nil {
		return false, err
	}

	peekChar := chars[0]
	if peekChar != c {
		return false, nil
	}

	l.Readch()
	return true, nil
}

// !! 扫描 获取token
// $env.a = 1
// $res.$body.$json 获取响应body的json格式
// $req.$header.auth = $res.a
func (l *Lexer) Scan() (Token, error) {
	for {
		err := l.Readch()
		if err == io.EOF {
			return NewToken(EOF), err
		}
		if err != nil {
			return NewToken(ERROR), err
		}

		if l.peek == ' ' || l.peek == '\t' {
			continue
		} else if l.peek == '\n' {
			l.line += 1
		} else {
			break
		}
	}

	l.Lexeme = ""

	switch l.peek {
	case '$':
		// 说明是关键字，一直读取后面的字符，直到非letter
		keyword, err := l.ScanKeyword()
		if err != nil {
			return NewToken(ERROR), err
		}

		return keyword.Tag, nil
	case '.':
		return NewToken(DOT), nil
	case '=':
		l.Lexeme = "="
		l.lexemeStack = append(l.lexemeStack, "=")
		return NewToken(ASSIGN_OPERATOR), nil
	}

	// 判断是否是数字
	if unicode.IsNumber(rune(l.peek)) {
		var v int
		var err error
		for {
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				l.UnRead()
				break
			} else {
				l.Lexeme += string(l.peek)
			}
			v = v*10 + num
			if l.Readch() == io.EOF {
				return NewToken(NUM), nil
			}
		}

		if l.peek != '.' {
			// 整型
			l.lexemeStack = append(l.lexemeStack, fmt.Sprint(v))
			return NewToken(NUM), err
		}
		l.Lexeme += string(l.peek)

		// 浮点型
		x := float64(v)
		d := float64(10)
		for {
			l.Readch()
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				l.UnRead()
				break
			}

			x = x + float64(num)/d
			d *= 10
			l.Lexeme += string(l.peek)
		}
		l.lexemeStack = append(l.lexemeStack, fmt.Sprint(x))
		return NewToken(REAL), err
	}

	// 读取变量字符串
	if unicode.IsLetter(rune(l.peek)) {
		var buffer []byte
		for {
			buffer = append(buffer, l.peek)
			l.Lexeme += string(l.peek)

			if err := l.Readch(); err == io.EOF {
				break
			}
			if !unicode.IsLetter(rune(l.peek)) {
				l.UnRead()
				break
			}
		}

		token, ok := l.keyWords[string(buffer)]
		if ok {
			return token, nil
		}
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		return NewToken(INDENTIFER), nil // 变量字符串
	}

	return NewToken(EOF), io.EOF
}

func (l *Lexer) ScanKeyword() (KeyWord, error) {
	var buffer []byte
	for {
		buffer = append(buffer, l.peek)
		l.Lexeme += string(l.peek)

		l.Readch()
		if !unicode.IsLetter(rune(l.peek)) {
			l.UnRead()
			break
		}
	}

	words := string(buffer)
	for _, item := range keyWord {
		if item.lexeme == words {
			return item, nil
		}
	}
	return KeyWord{}, errors.New("非关键字")
}
