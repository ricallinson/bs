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
	"strings"
)

type Token int // Token represents a lexical token.

const (
	PREV               = false
	NEXT               = true
	ILLEGAL      Token = iota
	EOF                // end of file
	EOL                // end of line
	WS                 // whitespace
	COMMENT            // #
	IDENT              //
	NUMBER             // 12345.67
	FUNCTION           // func
	RETURN             // return
	WHILE              // while
	FORIN              // for in
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
	PRINT              // print
	PRINTLN            // println
	LS                 // ls
	EXISTS             // exists
	CONCAT             // concat
	CALL               // call
	BASH               // bash
	ARRAY              // array
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
	keywords["print"] = PRINT
	keywords["println"] = PRINTLN
	keywords["ls"] = LS
	keywords["exists"] = EXISTS
	keywords["concat"] = CONCAT
	keywords["call"] = CALL
	keywords["bash"] = BASH
	keywords["array"] = ARRAY
	keywords["for"] = FORIN
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}
