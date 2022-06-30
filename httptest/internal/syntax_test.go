package internal

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// "$res.$body.$json" => "$res $body . $json ."
func TestSimpleParse(t *testing.T) {
	source := "$res.$body.$json"
	p := NewSimpleParser(NewLexer(source))
	node, err := p.Parse()
	require.Equal(t, io.EOF, err)
	fmt.Println(node.Attribute())
}
