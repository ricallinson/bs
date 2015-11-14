package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Token represents a lexical token.
type Token int

const (
	PREV               = false
	NEXT               = true
	ILLEGAL      Token = iota
	EOF                // end of file
	EOL                // end of line
	WS                 // whitespace
	IDENT              //
	NUMBER             // 12345.67
	FUNCTION           // func
	RETURN             // return
	WHILE              // while
	IF                 // if
	ELSE               // else
	TRUE               // true
	FALSE              // false
	DURATION_VAL       // 13h
	STRING             // "abc"
	BADSTRING          // "abc
	BADESCAPE          // \q
	REGEX              // Regular expressions
	BADREGEX           // `.*
	PRINTLN            // println
	ADD                // +
	SUB                // -
	MUL                // *
	DIV                // /
	AND                // AND
	OR                 // OR
	EQ                 // =
	NEQ                // !=
	EQREGEX            // =~
	NEQREGEX           // !~
	LT                 // <
	LTE                // <=
	GT                 // >
	GTE                // >=
	LSQUARE            // [
	RSQUARE            // ]
	LPAREN             // (
	RPAREN             // )
	LBRACE             // {
	RBRACE             // }
	COMMA              // ,
	COLON              // :
	SEMICOLON          // ;
	DOT                // .
)

var keywords map[string]Token

func init() {
	keywords = map[string]Token{}
	keywords["func"] = FUNCTION
	keywords["while"] = WHILE
	keywords["if"] = IF
	keywords["else"] = ELSE
	keywords["true"] = TRUE
	keywords["false"] = FALSE
	keywords["println"] = PRINTLN
	keywords["return"] = RETURN
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

// true = Next, false = Prev
func Find(t Token, n NodeI, dir bool) NodeI {
	for n.Token() != 0 {
		if n.Token() == t {
			return n
		}
		if dir {
			n = n.Next()
		} else {
			n = n.Prev()
		}
		if n == nil {
			return nil
		}
	}
	return nil
}

type NodeI interface {
	Ident(...string) string
	Token(...Token) Token
	Prev(...NodeI) NodeI
	Next(...NodeI) NodeI
	String() string
}

type Node struct {
	ident string
	token Token
	prev  NodeI
	next  NodeI
}

func (this *Node) Ident(v ...string) string {
	if len(v) == 1 {
		this.ident = v[0]
	}
	return this.ident
}

func (this *Node) Token(v ...Token) Token {
	if len(v) == 1 {
		this.token = v[0]
	}
	return this.token
}

func (this *Node) Prev(v ...NodeI) NodeI {
	if len(v) == 1 {
		this.prev = v[0]
	}
	return this.prev
}

func (this *Node) Next(v ...NodeI) NodeI {
	if len(v) == 1 {
		this.next = v[0]
	}
	return this.next
}

func (this *Node) String() string {
	str := ""
	switch this.Token() {
	case EQ:
		str = "="
	case ADD:
		if this.Prev().Token() != STRING && this.Prev().Token() != IDENT {
			str = " + "
		}
	case SUB:
		str = " - "
	case MUL:
		str = " * "
	case DIV:
		str = " / "
	case RETURN:
		str = "\"echo\" \"-ne\" "
	case TRUE:
		str = "$((1))"
	case FALSE:
		str = "$((0))"
	case COMMA:
		str = " "
	case LSQUARE:
		str = "("
	case RSQUARE:
		str = ")"
	}
	return str + this.Next().String()
}

type Block struct {
	*Node
	next []NodeI
}

func (this Block) String() string {
	str := " {"
	// Search the block for varibales and then check if they should be local.
	for _, b := range this.next {
		if b.Token() == IDENT {
			str += "\nlocal " + b.Ident()
		}
	}
	// Check back to see if we are in a function, if we are then print out the function arguments here.
	if n := Find(FUNCTION, this, PREV); n != nil {
		n = n.Next().Next().Next() // function->foo->(
		i := 1
		for n.Token() != RPAREN {
			str += fmt.Sprintf("\nlocal %s\n%s=\"$%d\"", n.Ident(), n.Ident(), i)
			n = n.Next()
			i++
		}
	}
	// Print the statments in the block.
	for _, b := range this.next {
		str += "" + b.String()
	}
	// If we are at the end of a function then print return.
	if Find(FUNCTION, this, PREV) != nil {
		str += "return\n"
	}
	return str + "}\n"
}

type EndOf struct {
	*Node
}

func (this EndOf) String() string {
	return "\n"
}

type Function struct {
	*Node
}

func (this *Function) String() string {
	return "function " + this.Next().String()
}

type FunctionName struct {
	*Node
}

func (this *FunctionName) String() string {
	return "\"" + this.Ident() + "\" " + this.Next().String()
}

type If struct {
	*Node
}

func (this If) String() string {
	return this.Next().String()
}

type Else struct {
	*Node
}

func (this Else) String() string {
	return this.Next().String()
}

type Variable struct {
	*Node
}

func (this *Variable) String() string {
	str := ""
	// Check to see if we are in a function, if we are do not print any vars.
	if Find(FUNCTION, this, PREV) != nil && Find(LPAREN, this, PREV) != nil && Find(RPAREN, this, NEXT) != nil {
		return this.Next().String()
	}
	switch this.Prev().Token() {
	case 0, LBRACE, EOL:
		// If there is no parent or EOL then this must be the first var to be assigned.
		// Check if this variable is a function.
		if false {
			str = fmt.Sprintf("\"%s\" ", this.Ident())
		} else {
			str = this.Ident()
		}
	case FUNCTION:
		// If the parent is a function then this is a function name not a var.
		str = this.Ident()
	default:
		// Are there any operators that makes this var a number?
		if Find(SUB, this, NEXT) != nil || Find(MUL, this, NEXT) != nil || Find(DIV, this, NEXT) != nil {
			str = fmt.Sprintf("$(($%s", this.Ident())
		} else {
			// Otherwise assume the var needs to be wrapped.
			str = fmt.Sprintf("\"$%s\"", this.Ident())
		}
	}
	return str + this.Next().String()
}

type String struct {
	*Node
}

func (this String) String() string {
	return fmt.Sprintf("\"%s\"", this.Ident()) + this.Next().String()
}

type Integer struct {
	*Node
}

func (this Integer) String() string {
	str := ""
	switch this.Prev().Token() {
	case ADD, SUB, DIV, MUL:
	default:
		str = "$(("
	}
	i, _ := strconv.Atoi(this.Ident())
	str += fmt.Sprintf("%d", i)
	switch this.Next().Token() {
	case ADD, SUB, DIV, MUL:
	default:
		str += "))"
	}
	return str + this.Next().String()
}
