package internal

import (
	"fmt"
	"io"
)

// 语法解析

// 抽象语法树的节点
type INode interface {
	Attribute() string
	Left() INode
	AddLeft(INode)
	Right() INode
	AddRight(INode)
}

type SyntaxNode struct {
	Type   string
	Name   string
	Params []*SyntaxNode

	Token     Token
	leftNode  INode
	rightNode INode
}

func NewSyntaxNode(token Token) *SyntaxNode {
	return &SyntaxNode{
		Token: token,
	}
}

func (s *SyntaxNode) Left() INode {
	return s.leftNode
}

func (s *SyntaxNode) Right() INode {
	return s.rightNode
}

func (s *SyntaxNode) AddLeft(node INode) {
	s.leftNode = node
}

func (s *SyntaxNode) AddRight(node INode) {
	s.rightNode = node
}

func (s *SyntaxNode) Attribute() string {
	if s == nil {
		return ""
	}

	var (
		left  string
		right string
	)
	if s.leftNode != nil {
		left = s.leftNode.Attribute()
	}
	if s.rightNode != nil {
		right = s.rightNode.Attribute()
	}

	return fmt.Sprintf("(%s, %s, %s)", left, s.Token.String(), right)
}

type SimpleParser struct {
	Lexer
}

// 遍历词法解析并构造语法抽象数
func NewSimpleParser(l Lexer) *SimpleParser {
	return &SimpleParser{
		Lexer: l,
	}
}

func (s *SimpleParser) Parse() (*SyntaxNode, error) {
	return s.list()
}

// 定义语义规则集，不同的符号有不同的规则
// 1、. 取值符号，一个表达式中可以存在多个，将 a . b作为新的左参数
// 2、= 赋值符号，左边的为左参数，右边的为右参数
func (s *SimpleParser) list() (*SyntaxNode, error) {
	var currNode *SyntaxNode

	stack := []*SyntaxNode{}
	for {
		token, err := s.Scan()
		if token.Tag == EOF {
			return currNode, io.EOF
		}
		if err != nil {
			return nil, err
		}

		switch token.Tag {
		case DOT:
			// 把下一个取出来
		// 	token, err := s.Scan()
		// if token.Tag == EOF {
		// 	return currNode, io.EOF
		// }
		// &SyntaxNode{
		// 	Type: "Expression",
		// 	Name: ".",
		// 	Params: [],
		// }
		case EQ:
			currNode = NewSyntaxNode(token)
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			currNode.leftNode = left
		case ENV, BODY, REQ, RES, JSON, RAW, INDENTIFER:
			if currNode != nil {
				currNode.rightNode = NewSyntaxNode(token)
			}
			stack = append(stack, currNode)
		}
	}
}
