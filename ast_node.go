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

type NodeI interface {
	Ident(...string) string
	Token(...Token) Token
	Prev(...NodeI) NodeI
	Next(...NodeI) NodeI
	String() (string, NodeI)
}

type Node struct {
	parser *Parser
	ident  string
	token  Token
	prev   NodeI
	next   NodeI
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
		str = "+"
	case SUB:
		str = "-"
	case MUL:
		str = "*"
	case DIV:
		str = "/"
	case LT:
		str = "<"
	case GT:
		str = ">"
	case LTE:
		str = "<="
	case GTE:
		str = ">="
	case EEQ:
		str = "=="
	case TRUE:
		str = "$((1))"
	case FALSE:
		str = "$((0))"
	case COMMA:
		// Do nothing.
	case LSQUARE:
		str = "("
	case RSQUARE:
		str = ")"
	case LS:
		str = "$(\"ls\")"
	case PRINT:
		str = "\"echo\" \"-ne\" "
	case PRINTLN:
		str = "\"echo\" \"-e\" "
	case RETURN:
		str = "\"echo\" \"-ne\" "
	}
	return str, this.Next()
}
