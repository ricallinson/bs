/*
Copyright 2015 Richard Allinson

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	// "fmt"
	"io"
	"os"
	"strings"
)

const (
	INDENT = "    "
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
	funcs map[string]bool
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// NewStringParser returns a new instance of Parser.
func NewParserFromFile(path string) *Parser {
	r, err := os.Open(path)
	if err != nil {
		panic("File not found")
	}
	p := NewParser(r)
	return p
}

// NewParserFromString returns a new instance of Parser.
func NewParserFromString(in string) *Parser {
	p := NewParser(strings.NewReader(in))
	return p
}

func ParseString(in string) string {
	p := NewParserFromString(in)
	return p.ParseToString()
}

func ParseFile(path string) string {
	p := NewParserFromFile(path)
	return p.ParseToString()
}

func (this *Parser) ParseToString() string {
	root, err := this.Parse()
	if err != nil {
		panic(err)
	}
	// Print each statement in the block.
	var str string
	s, n := root.String()
	str += s
	for n != nil {
		s, n = n.String()
		str += s
	}
	return str
}

// Parse parses a scripter string and returns a slice of statment AST objects.
func (this *Parser) Parse(n ...NodeI) (NodeI, error) {
	if this.funcs == nil {
		this.funcs = map[string]bool{}
	}
	var root NodeI
	if len(n) == 1 {
		root = n[0]
	} else {
		root = &Node{parser: this}
	}
	if err := this.parseStatement(root); err != nil {
		return nil, err
	}
	return root, nil
}

// Parse parses a statement.
func (this *Parser) parseStatement(prev NodeI) error {
	var node *Node = &Node{parser: this}
	tok, lit := this.scanIgnoreWhitespace()
	if tok == EOF || tok == RBRACE {
		node.Token(tok)
		node.Prev(prev)
		prev.Next(node)
		return nil
	}
	var curr NodeI
	switch tok {
	case IDENT:
		curr = &Variable{node}
		// Check if the ident is a function, if so return a function name.
		if ok, _ := this.funcs[lit]; ok {
			curr = &FunctionName{node}
		}
	case NUMBER:
		curr = &Integer{node}
	case STRING:
		curr = &String{node}
	case FUNCTION:
		curr = &Function{node}
	case IF:
		curr = &If{node}
	case ELSE:
		curr = &Else{node}
	case WHILE:
		curr = &While{node}
	case FORIN:
		curr = &ForIn{node}
	case CALL:
		curr = &Call{node}
	case BASH:
		curr = &Bash{node}
	case ARRAY:
		curr = &Array{node}
	case COMMENT:
		curr = &Comment{node}
	case APPEND:
		curr = &Append{node}
	case LEN:
		curr = &Len{node}
	case LBRACE:
		// We create a node here to represent the block.
		n := &Node{this, "", LBRACE, prev, nil}
		block, _ := this.Parse(n)
		curr = &Block{node, block}
	default:
		curr = node
	}

	curr.Ident(lit)
	curr.Token(tok)
	curr.Prev(prev)
	prev.Next(curr)

	if prev.Token() == FUNCTION {
		if ok, _ := this.funcs[curr.Ident()]; ok {
			panic("function " + curr.Ident() + " already defined")
		}
		this.funcs[curr.Ident()] = true
	}

	return this.parseStatement(curr)
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (this *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if this.buf.n != 0 {
		this.buf.n = 0
		return this.buf.tok, this.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, _, lit = this.s.Scan()

	// Save it to the buffer in case we unscan later.
	this.buf.tok, this.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (this *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = this.scan()
	if tok == WS {
		tok, lit = this.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (this *Parser) unscan() { this.buf.n = 1 }
