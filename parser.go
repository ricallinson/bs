package main

import (
	// "fmt"
	"io"
	"os"
	"strings"
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
func NewParserFromString(in string) *Parser {
	p := NewParser(strings.NewReader(in))
	p.funcs = map[string]bool{}
	return p
}

func ParseString(in string) string {
	p := NewParserFromString(in)
	return p.ParseToString()
}

// NewStringParser returns a new instance of Parser.
func NewParserFromFile(path string) *Parser {
	r, _ := os.Open(path)
	p := NewParser(r)
	p.funcs = map[string]bool{}
	return p
}

func ParseFile(path string) string {
	p := NewParserFromFile(path)
	return p.ParseToString()
}

func (this *Parser) ParseToString() string {
	s, e := this.Parse()
	if e != nil {
		panic(e)
	}
	var str string
	for _, l := range s {
		str += l.String()
	}
	return str
}

// Parse parses a scripter string and returns a slice of statment AST objects.
func (this *Parser) Parse(n ...NodeI) ([]NodeI, error) {
	var prev NodeI
	if len(n) == 1 {
		prev = n[0]
	} else {
		prev = &Node{}
	}
	var nodes = []NodeI{}
	for {
		if tok, _ := this.scanIgnoreWhitespace(); tok == EOF || tok == RBRACE {
			return nodes, nil
		} else {
			this.unscan()
			node, err := this.ParseStatement(prev)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
}

// Parse parses a statement.
func (this *Parser) ParseStatement(prev NodeI) (NodeI, error) {
	var curr NodeI
	tok, lit := this.scanIgnoreWhitespace()
	if tok == EOL || tok == EOF {
		curr = &EndOf{&Node{}}
		curr.Token(tok)
		return curr, nil
	}

	switch tok {
	case IDENT:
		curr = &Variable{&Node{}}
		// Check if the ident is a function, if so return a function name.
		if ok, _ := this.funcs[lit]; ok {
			curr = &FunctionName{&Node{}}
		}
	case NUMBER:
		curr = &Integer{&Node{}}
	case STRING:
		curr = &String{&Node{}}
	case FUNCTION:
		curr = &Function{&Node{}}
	case IF:
		curr = &If{&Node{}}
	case ELSE:
		curr = &Else{&Node{}}
	case LBRACE:
		// We create a node here to represent the block.
		node := &Node{"{", LBRACE, prev, nil}
		block, _ := this.Parse(node)
		b := &Block{&Node{}, block}
		curr = b
	case WHILE:
		curr = &While{&Node{}}
	default:
		curr = &Node{}
	}

	next, e := this.ParseStatement(curr)
	if e != nil {
		return nil, e
	}

	curr.Ident(lit)
	curr.Token(tok)
	curr.Prev(prev)
	curr.Next(next)

	if curr.Token() == FUNCTION {
		if ok, _ := this.funcs[next.Ident()]; ok {
			panic("function " + next.Ident() + " already defined")
		}
		this.funcs[next.Ident()] = true
	}

	return curr, nil
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
