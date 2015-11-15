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
	EEQ                // ==
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
	for n.Token() != ILLEGAL {
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

func GetMyFunction(n NodeI) NodeI {
	for n.Token() != ILLEGAL {
		if n.Token() == IF || n.Token() == WHILE {
			return nil
		} else if n.Token() == FUNCTION {
			return n
		}
		n = n.Prev()
		if n == nil {
			return nil
		}
	}
	return nil
}

func IsArithmetic(n NodeI) bool {
	switch n.Token() {
	case ADD, SUB, DIV, MUL, LT, GT, LTE, GTE:
		return true
	}
	return false
}

func WasArithmetic(n NodeI) bool {
	for n.Token() != IF && n.Token() != WHILE && n.Token() != ILLEGAL {
		switch n.Token() {
		case ADD, SUB, DIV, MUL, LT, GT, LTE, GTE:
			return true
		}
		n = n.Prev()
		if n == nil {
			return false
		}
	}
	return false
}

func WasReturn(n NodeI) bool {
	for n.Token() != ILLEGAL {
		if n.Token() == RETURN {
			return true
		}
		n = n.Next()
		if n == nil {
			return false
		}
	}
	return false
}

func InList(n NodeI) bool {
	for n.Token() != EOL && n.Token() != ILLEGAL {
		if n.Token() == LSQUARE {
			return true
		}
		n = n.Prev()
		if n == nil {
			return false
		}
	}
	return false
}

func CountVariables(n NodeI) int {
	var i int
	for n.Token() != IF && n.Token() != WHILE && n.Token() != ILLEGAL {
		switch n.Token() {
		case IDENT, NUMBER:
			i++
		}
		n = n.Prev()
		if n == nil {
			return i
		}
	}
	return i
}

func CountIfs(n NodeI) int {
	var i int
	for n.Token() != ILLEGAL {
		switch n.Token() {
		case IF:
			i++
			if n.Prev().Token() != ELSE {
				return i
			}
		}
		n = n.Prev()
		if n == nil {
			return i
		}
	}
	return i
}

func Things(node NodeI) (string, NodeI) {
	var str string
	var s string
	for node.Token() != RPAREN {
		s, node = node.String()
		str += s
	}
	return str, node
}

type NodeI interface {
	Ident(...string) string
	Token(...Token) Token
	Prev(...NodeI) NodeI
	Next(...NodeI) NodeI
	String() (string, NodeI)
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

func (this *Node) String() (string, NodeI) {
	var str string
	switch this.Token() {
	case EOL:
		str = "\n"
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
	case EEQ:
		str = " == "
	case TRUE:
		str = "$((1))"
	case FALSE:
		str = "$((0))"
	case COMMA:
		if InList(this) {
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
	case RETURN:
		str = "\"echo\" \"-ne\" "
	}
	return str, this.Next()
}

type Block struct {
	*Node
	block NodeI
}

func (this Block) String() (string, NodeI) {
	var str string
	if Find(IF, this, PREV) == nil && Find(WHILE, this, PREV) == nil {
		str += " {"
	}
	// Check back to see if we are in a function, if we are then print out the function arguments here.
	if n := GetMyFunction(this); n != nil {
		// Get the arguments from the function.
		var argsStr string
		argsList := map[string]bool{}
		args := n.Next().Next().Next() // function->foo->(
		i := 1
		for args.Token() != RPAREN {
			if args.Token() == IDENT {
				argsList[args.Ident()] = true
				argsStr += fmt.Sprintf("\nlocal %s\n%s=\"$%d\"", args.Ident(), args.Ident(), i)
				i++
			}
			args = args.Next()
		}
		// Search the block for variables and then check if they should be local.
		_, local := this.block.String()
		for local != nil {
			if local.Token() == IDENT {
				if ok := argsList[local.Ident()]; ok == false {
					str += "\nlocal " + local.Ident()
				}
			}
			_, local = local.String()
		}
		str += argsStr
	}
	// Print each statement in the block.
	s, line := this.block.String()
	str += s
	for line != nil {
		s, line = line.String()
		str += s
	}
	// Was there a return in this block?
	if WasReturn(this.block) {
		str += "return\n"
	}
	// Check to see what type of block we are in.
	if GetMyFunction(this) != nil {
		// If we are at the end of a function then print return.
		str += "}"
	} else if Find(IF, this, PREV) != nil {
		if this.Next().Token() != ELSE {
			// If we are at the end of a if block so count the total if then print an for each fi.
			for i := 0; i < CountIfs(this); i++ {
				str += "fi\n"
			}
			str = str[:len(str)-1]
		}
	} else if Find(WHILE, this, PREV) != nil {
		// If we are at the end of a while block then print done.
		str += "done\n"
	}
	return str, this.Next()
}

type Function struct {
	*Node
}

func (this *Function) String() (string, NodeI) {
	_, node := Things(this.Next().Next()) // func->name->(
	return "function " + this.Next().Ident(), node.Next()
}

type FunctionName struct {
	*Node
}

func (this *FunctionName) String() (string, NodeI) {
	return "\"" + this.Ident() + "\" ", this.Next()
}

type While struct {
	*Node
}

func (this While) String() (string, NodeI) {
	str, node := Things(this.Next()) // while->(
	check := ""
	if WasArithmetic(node.Prev()) || CountVariables(node.Prev()) == 1 {
		check = " == 1"
	}
	return "while [ " + str + check + " ]; do", node.Next()
}

type If struct {
	*Node
}

func (this If) String() (string, NodeI) {
	str, node := Things(this.Next()) // if->(
	check := ""
	if WasArithmetic(node.Prev()) || CountVariables(node.Prev()) == 1 {
		check = " == 1"
	}
	return "if [ " + str + check + " ]; then", node.Next()
}

type Else struct {
	*Node
}

func (this Else) String() (string, NodeI) {
	if this.Next().Token() == IF {
		return "else\n", this.Next()
	}
	return "else", this.Next()
}

type Variable struct {
	*Node
}

func (this *Variable) String() (string, NodeI) {
	var str string
	// Are we in an exists function call, if so do something special.
	if this.Next().Next().Token() == EXISTS { // a->=->exists->(->"str"->)
		param, _ := this.Next().Next().Next().Next().String()
		str = "[ -e " + param + " ]\n" + this.Ident() + "=$((!$?))"
		return str, this.Next().Next().Next().Next().Next()
	}
	// If we got here then then it is a normal variable so check the previous token to see whats going on.
	switch this.Prev().Token() {
	case 0, LBRACE, EOL:
		// If there is no parent, a { or an EOL then this must be the first var to be assigned.
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
	return str, this.Next()
}

type String struct {
	*Node
}

func (this String) String() (string, NodeI) {
	return fmt.Sprintf("\"%s\"", this.Ident()), this.Next()
}

type Integer struct {
	*Node
}

func (this Integer) String() (string, NodeI) {
	var str string
	if IsArithmetic(this.Prev()) == false {
		str = "$(("
	}
	i, _ := strconv.Atoi(this.Ident())
	str += fmt.Sprintf("%d", i)
	if IsArithmetic(this.Next()) == false {
		str += "))"
	}
	return str, this.Next()
}
