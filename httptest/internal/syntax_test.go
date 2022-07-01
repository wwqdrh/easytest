package internal

import (
	"encoding/json"
	"io"
	"reflect"
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
								"type": "global",
								"name": "$res"
							}, {
								"type": "global",
								"name": "$body"
							}
						]
					}, {
						"type": "attr",
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
					"type": "global",
					"name": "$env"
				}, {
					"type": "variable",
					"name": "indentifer",
					"value": "a"
				}
			]
		},
		{
			"type": "literial",
			"name": "num",
			"value": 1
		}
	]
}`,
		},
	}

	for _, item := range pairs {
		p := NewSimpleParser(NewLexer(item.source))
		node, err := p.Parse()
		assert.Equal(t, io.EOF, err)

		var target SyntaxNode
		var source SyntaxNode
		assert.Nil(t, json.Unmarshal([]byte(item.target), &target))
		sourceData, err := json.Marshal(node)
		assert.Nil(t, err)
		assert.Nil(t, json.Unmarshal(sourceData, &source))
		assert.True(t, reflect.DeepEqual(target, source))
	}
}
