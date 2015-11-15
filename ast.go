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
	PRINTLN            // println
	LS                 // ls
	EXISTS             // exists
	CONCAT             // concat
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
	keywords["return"] = RETURN
	// Built-ins
	keywords["println"] = PRINTLN
	keywords["ls"] = LS
	keywords["exists"] = EXISTS
	keywords["concat"] = CONCAT
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

func IsConstruct(node NodeI) bool {
	// If there is a left brace then we are in the construct.
	if Find(LBRACE, node, PREV) != nil {
		return false
	}
	if (Find(FUNCTION, node, PREV) != nil || Find(WHILE, node, PREV) != nil) && Find(LPAREN, node, PREV) != nil && Find(RPAREN, node, NEXT) != nil {
		return true
	}
	return false
}

func IsArithmetic(node NodeI) bool {
	switch node.Token() {
	case ADD, SUB, DIV, MUL:
		return true
	}
	return false
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
	// Check to see if we are in a construct, if we are do not print any vars.
	if IsConstruct(this) {
		return this.Next().String()
	}
	str := ""
	switch this.Token() {
	case EQ:
		str = "="
	case ADD:
		str = " + "
	case SUB:
		str = " - "
	case MUL:
		str = " * "
	case DIV:
		str = " / "
	case LT:
		str = " < "
	case GT:
		str = " > "
	case LTE:
		str = " <= "
	case GTE:
		str = " >= "
	case RETURN:
		str = "\"echo\" \"-ne\" "
	case TRUE:
		str = "$((1))"
	case FALSE:
		str = "$((0))"
	case COMMA:
		if Find(LSQUARE, this, PREV) != nil {
			// If we are in a list then add a space.
			str = " "
		}
	case LSQUARE:
		str = "("
	case RSQUARE:
		str = ")"
	case LS:
		str = "$(\"ls\")"
	case PRINTLN:
		str = "\"echo\" \"-e\" "
	}
	return str + this.Next().String()
}

type Block struct {
	*Node
	next []NodeI
}

func (this Block) String() string {
	str := ""
	if Find(WHILE, this, PREV) == nil {
		str += " {"
	}
	// Check back to see if we are in a function, if we are then print out the function arguments here.
	if n := Find(FUNCTION, this, PREV); n != nil {
		// Search the block for variables and then check if they should be local.
		for _, b := range this.next {
			if b.Token() == IDENT {
				str += "\nlocal " + b.Ident()
			}
		}
		n = n.Next().Next().Next() // function->foo->(
		i := 1
		for n.Token() != RPAREN {
			if n.Token() == IDENT {
				str += fmt.Sprintf("\nlocal %s\n%s=\"$%d\"", n.Ident(), n.Ident(), i)
				i++
			}
			n = n.Next()
		}
	}
	// Print each statement in the block.
	for _, b := range this.next {
		str += "" + b.String()
	}
	// Check to see what type of block we are in.
	if Find(FUNCTION, this, PREV) != nil {
		// If we are at the end of a function then print return.
		str += "return\n}\n"
	} else if Find(WHILE, this, PREV) != nil {
		// If we are at the end of a while block then print done.
		str += "done\n"
	}
	return str
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

type While struct {
	*Node
}

func (this While) String() string {
	str := ""
	n := this.Next().Next() // while->(
	i := 1
	for n.Token() != RPAREN {
		switch n.Token() {
		case IDENT:
			str += fmt.Sprintf("$%s", n.Ident())
		case NUMBER:
			str += fmt.Sprintf("%s", n.Ident())
		case LT:
			str += " < "
		case GT:
			str += " > "
		case LTE:
			str += " <= "
		case GTE:
			str += " >= "
		}
		n = n.Next()
		i++
	}
	return "while [ $((" + str + ")) == 1 ]; do" + this.Next().String()
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
	// Check to see if we are in a construct, if we are do not print any vars.
	if IsConstruct(this) {
		return this.Next().String()
	}
	var str string
	// Are we in an exists function call, if so do something special.
	if Find(EXISTS, this, NEXT) != nil { // a=exists("str")
		str = "[ -e \"file.txt\" ]\n" + this.Ident() + "=$((!$?))"
		return str + this.Next().Next().Next().Next().Next().String()
	}
	// If we got here then then it is a normal variable so check the previous token to see whats going on.
	switch this.Prev().Token() {
	case 0, LBRACE, EOL:
		// If there is no parent, a { or an EOL then this must be the first var to be assigned.
		str = this.Ident()
	case FUNCTION:
		// If the parent is a function then this is a function name not a variable.
		str = this.Ident()
	default:
		// Are there any operators that makes this variable a number?
		if IsArithmetic(this.Prev()) == false && IsArithmetic(this.Next()) {
			str = fmt.Sprintf("$(($%s", this.Ident())
		} else if IsArithmetic(this.Prev()) && IsArithmetic(this.Next()) == false {
			str = fmt.Sprintf("$%s))", this.Ident())
		} else if IsArithmetic(this.Prev()) && IsArithmetic(this.Next()) {
			str = fmt.Sprintf("$%s", this.Ident())
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
	// Check to see if we are in a construct, if we are do not print any vars.
	if IsConstruct(this) {
		return this.Next().String()
	}
	str := ""
	if IsArithmetic(this.Prev()) == false {
		str = "$(("
	}
	i, _ := strconv.Atoi(this.Ident())
	str += fmt.Sprintf("%d", i)
	if IsArithmetic(this.Next()) == false {
		str += "))"
	}
	return str + this.Next().String()
}
