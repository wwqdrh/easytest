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
	STR
	BODY
	JSON
	HEADER

	// 全局函数
	CONTAIN

	// 操作符
	EQ
	ASSIGN_OPERATOR
	DOT

	// 括号
	LEFT_PATREN
	RIGHT_PATERN
	COMMA

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
	STR:             "$str",
	BODY:            "$body",
	JSON:            "$json",
	HEADER:          "$header",
	CONTAIN:         "@contain",
	EQ:              "==",
	ASSIGN_OPERATOR: "=",
	DOT:             ".",
	LEFT_PATREN:     "(",
	RIGHT_PATERN:    ")",
	COMMA:           ",",
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
	NewKeyWord(STR),
	NewKeyWord(JSON),
	NewKeyWord(BODY),
	NewKeyWord(HEADER),
	NewKeyWord(CONTAIN),
}

// token字符分类
type Token struct {
	Tag Tag
	Raw interface{}
}

func NewToken(tag Tag) Token {
	token := Token{
		Tag: tag,
	}
	token.Raw = token.String()
	return token
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
	peek        rune             // 读入的字符
	line        int              // 当前字符串处于第几行
	reader      *bufio.Reader    // 用于读取字节流
	keyWords    map[string]Token // 存储关键字
}

func NewLexer(source string) Lexer {
	str := strings.NewReader(source)
	sourceReader := bufio.NewReaderSize(str, len([]rune(source)))
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
		if err := l.reader.UnreadRune(); err != nil {
			return
		}
	}

	l.lexemeStack = l.lexemeStack[:len(l.lexemeStack)-1]
	l.Lexeme = l.lexemeStack[len(l.lexemeStack)-1]
}

func (l *Lexer) Readch() error {
	r, _, err := l.reader.ReadRune()
	l.peek = r
	return err
}

func (l *Lexer) UnRead() error {
	return l.reader.UnreadRune()
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

	if err := l.Readch(); err != nil {
		return false, err
	}
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
	case '@':
		// 说明是函数
		keyword, err := l.ScanKeyword()
		if err != nil {
			return NewToken(ERROR), err
		}
		return keyword.Tag, nil
	case '(':
		return NewToken(LEFT_PATREN), nil
	case ')':
		return NewToken(RIGHT_PATERN), nil
	case ',':
		return NewToken(COMMA), nil
	case '.':
		return NewToken(DOT), nil
	case '"':
		// 字符串
		return l.ScanString()
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
				if err := l.UnRead(); err != nil {
					break
				}
				break
			} else {
				l.Lexeme += string(l.peek)
			}
			v = v*10 + num
			if l.Readch() == io.EOF {
				token := NewToken(NUM)
				token.Raw = v
				return token, nil
			}
		}

		if l.peek != '.' {
			// 整型
			l.lexemeStack = append(l.lexemeStack, fmt.Sprint(v))
			token := NewToken(NUM)
			token.Raw = v
			return token, err
		}
		l.Lexeme += string(l.peek)

		// 浮点型
		x := float64(v)
		d := float64(10)
		for {
			if err := l.Readch(); err != nil {
				break
			}
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				if err := l.UnRead(); err != nil {
					break
				}
				break
			}

			x = x + float64(num)/d
			d *= 10
			l.Lexeme += string(l.peek)
		}
		l.lexemeStack = append(l.lexemeStack, fmt.Sprint(x))
		token := NewToken(REAL)
		token.Raw = x
		return token, err
	}

	// 读取变量字符串
	if unicode.IsLetter(rune(l.peek)) {
		var buffer []rune
		for {
			buffer = append(buffer, l.peek)
			l.Lexeme += string(l.peek)

			if err := l.Readch(); err == io.EOF {
				break
			}
			if !unicode.IsLetter(rune(l.peek)) {
				if err := l.UnRead(); err != nil {
					break
				}
				break
			}
		}

		token, ok := l.keyWords[string(buffer)]
		if ok {
			return token, nil
		}
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)

		token = NewToken(INDENTIFER)
		token.Raw = string(buffer)
		return token, nil // 变量字符串
	}

	return NewToken(EOF), io.EOF
}

func (l *Lexer) ScanKeyword() (KeyWord, error) {
	var buffer []rune
	for {
		buffer = append(buffer, l.peek)
		l.Lexeme += string(l.peek)

		if err := l.Readch(); err == io.EOF {
			break
		}
		if !unicode.IsLetter(rune(l.peek)) {
			if err := l.UnRead(); err != nil {
				break
			}
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

func (l *Lexer) ScanString() (Token, error) {
	var buffer []rune
	for {
		if err := l.Readch(); err == io.EOF {
			break
		}

		if l.peek == '\\' {
			// 转义符
			if err := l.Readch(); err == io.EOF {
				break
			}

			buffer = append(buffer, l.peek)
			l.Lexeme += string(l.peek)
			continue
		}

		if l.peek == '"' {
			break
		}

		buffer = append(buffer, l.peek)
		l.Lexeme += string(l.peek)
	}

	return Token{
		INDENTIFER,
		string(buffer),
	}, nil
}
