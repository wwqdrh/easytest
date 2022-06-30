package internal

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// "$res.$body.$json" => "$res $body . $json ."
func TestSimpleParse(t *testing.T) {
	var pairs = []struct{ source, target string }{
		{
			source: "$res.$body.$json",
			target: `
			{
				"type": "expression",
				"name": ".",
				"params": [
					{
						"type": "expression",
						"name": ".",
						"params": [
							{
								"type": "variable",
								"name": "$res"
							}, {
								"type": "variable",
								"name": "$body"
							}
						]
					}, {
						"type": "variable",
						"name": "$json"
					}
				]
			}
			`,
		},
		{
			source: "$env.a = 1",
			target: `
{
	"type": "expression",
	"name": "=",
	"params": [
		{
			"type": "expression",
			"name": ".",
			"params": [
				{
					"type": "variable",
					"params": [
						"type": "global",
						"name": "$env"
					]
				},
				{
					"type": "variable",
					"params": [
						"type": "attr",
						"name": "a"
					]
				}
			]
		},
		{
			"type": "literal",
			"name": "1"
		}
	]
}`,
		},
	}

	for _, item := range pairs {
		p := NewSimpleParser(NewLexer(item.source))
		node, err := p.Parse()
		assert.Equal(t, io.EOF, err)
		assert.Equal(t, item.target, node.Attribute())
	}
}
