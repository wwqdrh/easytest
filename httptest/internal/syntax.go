package internal

import (
	"errors"
	"fmt"
	"io"
)

// 语法解析

// 抽象语法树的节点

var ErrAst = errors.New("AST ERROR")

type SimpleParser struct {
	Lexer

	nodeStack []*SyntaxNode
}

type SyntaxNode struct {
	Type   string
	Name   string
	Value  interface{} `json:"-"` // TODO float64与int类型不同问题 不好测试
	Token  Token
	Params []*SyntaxNode
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
	for {
		token, err := s.Scan()
		if token.Tag == EOF {
			node := s.nodeStack[len(s.nodeStack)-1]
			s.nodeStack = s.nodeStack[:len(s.nodeStack)-1]
			return node, io.EOF
		}
		if err != nil {
			return nil, err
		}

		switch token.Tag {
		case CONTAIN:
			if err := s.parseContain(); err != nil {
				return nil, err
			}
		case DOT:
			if err := s.parseDOT(); err != nil {
				return nil, err
			}
		case COMMA, RIGHT_PATERN:
			// 解析到逗号，当前解析的部分返回
			node := s.nodeStack[len(s.nodeStack)-1]
			s.nodeStack = s.nodeStack[:len(s.nodeStack)-1]
			return node, nil
		case ASSIGN_OPERATOR:
			if err := s.parseASSIGN(); err != nil {
				return nil, err
			}
		case ENV, BODY, REQ, RES, JSON, RAW, INDENTIFER, NUM:
			s.nodeStack = append(s.nodeStack, s.builderNode(token))
		}
	}
}

// 下一个必须为left_pate然后有一个right_pate, 并且两者之间相差两个元素
func (s *SimpleParser) parseContain() error {
	nextToken, err := s.Scan()
	if nextToken.Tag != LEFT_PATREN || err != nil {
		return fmt.Errorf("%w: contain下一个token必须为LEFT_PATREN", ErrAst)
	}

	param1, err := s.list()
	if err != nil {
		return fmt.Errorf("%w: 解析失败", ErrAst)
	}

	param2, err := s.list()
	if err != nil {
		return fmt.Errorf("%w: 解析失败", ErrAst)
	}

	call := s.builderNode(NewToken(CONTAIN))
	call.Params = append(call.Params, param1, param2)
	s.nodeStack = append(s.nodeStack, call)
	return nil
}

func (s *SimpleParser) parseDOT() error {
	// 把下一个token取出来作为当前树节点的右节点
	// stack前面的作为当前树节点的左节点
	left := s.nodeStack[len(s.nodeStack)-1]
	s.nodeStack = s.nodeStack[:len(s.nodeStack)-1]

	nextToken, err := s.Scan()
	if err != nil {
		return err
	}

	node := s.builderNode(NewToken(DOT))

	node.Params = []*SyntaxNode{
		left, s.builderNode(nextToken),
	}
	s.nodeStack = append(s.nodeStack, node)
	return nil
}

func (s *SimpleParser) parseASSIGN() error {
	left := s.nodeStack[len(s.nodeStack)-1]
	s.nodeStack = s.nodeStack[:len(s.nodeStack)-1]

	right, err := s.list()
	if err != nil && err != io.EOF {
		return err
	}

	eqNode := s.builderNode(NewToken(ASSIGN_OPERATOR))
	eqNode.Params = []*SyntaxNode{left, right}
	s.nodeStack = append(s.nodeStack, eqNode)
	return nil
}

func (s *SimpleParser) builderNode(token Token) *SyntaxNode {
	switch token.Tag {
	case ENV, BODY, REQ, RES:
		return &SyntaxNode{
			Type: "global",
			Name: token.String(),
		}
	case CONTAIN:
		return &SyntaxNode{
			Type: "callable",
			Name: token.String(),
		}
	case JSON, RAW, STR:
		return &SyntaxNode{
			Type: "attr",
			Name: token.String(),
		}
	case DOT, EQ, ASSIGN_OPERATOR:
		return &SyntaxNode{
			Type: "expression",
			Name: token.String(),
		}
	case INDENTIFER:
		return &SyntaxNode{
			Type:  "variable",
			Name:  token.String(),
			Value: token.Raw,
		}
	case NUM:
		return &SyntaxNode{
			Type:  "literial",
			Name:  token.String(),
			Value: token.Raw,
		}
	default:
		return nil
	}
}
