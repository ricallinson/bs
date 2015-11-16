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
		root = &Node{}
	}
	if err := this.ParseStatement(root); err != nil {
		return nil, err
	}
	return root, nil
}

// Parse parses a statement.
func (this *Parser) ParseStatement(prev NodeI) error {
	var curr NodeI
	tok, lit := this.scanIgnoreWhitespace()
	if tok == EOF || tok == RBRACE {
		curr = &Node{}
		curr.Token(tok)
		curr.Prev(prev)
		prev.Next(curr)
		return nil
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
	case WHILE:
		curr = &While{&Node{}}
	case CALL:
		curr = &Call{&Node{}}
	case BASH:
		curr = &Bash{&Node{}}
	case LBRACE:
		// We create a node here to represent the block.
		node := &Node{"", LBRACE, prev, nil}
		block, _ := this.Parse(node)
		curr = &Block{&Node{}, block}
	default:
		curr = &Node{}
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

	return this.ParseStatement(curr)
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
