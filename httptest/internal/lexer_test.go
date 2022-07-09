package internal

import (
	"fmt"
	"testing"
)

// $env.a = 1
// $env.token = $res.$body.$json.token
// $req.$header.auth = $env.token
// @contain($res.$body, "ok")
func TestLexerParse(t *testing.T) {
	source := "$env.a = 1"
	lexerParse := NewLexer(source)
	for {
		token, _ := lexerParse.Scan()
		fmt.Println(token.String())
		if token.Tag == EOF {
			break
		}
	}

	source = "$env.token = $res.$body.$json.token"
	lexerParse = NewLexer(source)
	for {
		token, _ := lexerParse.Scan()
		fmt.Println(token.String())
		if token.Tag == EOF {
			break
		}
	}

	source = "$req.$header.auth = $env.token"
	lexerParse = NewLexer(source)
	for {
		token, _ := lexerParse.Scan()
		fmt.Println(token.String())
		if token.Tag == EOF {
			break
		}
	}

	source = `@contain($res.$body.$str, "ok")`
	lexerParse = NewLexer(source)
	for {
		token, _ := lexerParse.Scan()
		fmt.Println(token.String(), token.Raw)
		if token.Tag == EOF {
			break
		}
	}
}
