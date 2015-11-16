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

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
    if tok, ok := keywords[strings.ToLower(ident)]; ok {
        return tok
    }
    return IDENT
}

// Look back from the given node and return a 'function' node if the given node is with a 'function' block.
func GetMyFunction(n NodeI) NodeI {
    for n.Token() != ILLEGAL {
        if n.Token() == IF || n.Token() == WHILE || n.Token() == FORIN {
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

// Look back from the given node and return an 'if' node if the given node is with a 'if' block.
func GetMyIf(n NodeI) NodeI {
    for n.Token() != ILLEGAL {
        if n.Token() == FUNCTION || n.Token() == WHILE || n.Token() == FORIN {
            return nil
        } else if n.Token() == IF {
            return n
        }
        n = n.Prev()
        if n == nil {
            return nil
        }
    }
    return nil
}

// Look back from the given node and return an 'while' node if the given node is with a 'while' block.
func GetMyWhile(n NodeI) NodeI {
    for n.Token() != ILLEGAL {
        if n.Token() == FUNCTION || n.Token() == IF || n.Token() == FORIN {
            return nil
        } else if n.Token() == WHILE {
            return n
        }
        n = n.Prev()
        if n == nil {
            return nil
        }
    }
    return nil
}

// Look back from the given node and return an 'for in' node if the given node is with a 'for in' block.
func GetMyForIn(n NodeI) NodeI {
    for n.Token() != ILLEGAL {
        if n.Token() == FUNCTION || n.Token() == IF || n.Token() == WHILE {
            return nil
        } else if n.Token() == FORIN {
            return n
        }
        n = n.Prev()
        if n == nil {
            return nil
        }
    }
    return nil
}

// Is the given node and arithmetic operator.
func IsArithmetic(n NodeI) bool {
    switch n.Token() {
    case ADD, SUB, DIV, MUL, LT, GT, LTE, GTE:
        return true
    }
    return false
}

// Look back to see if the given node is part of an arithmetic statement.
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

// Look back to see if the given node was in a function with a return.
func WasReturn(n NodeI) bool {
    for n.Token() != FUNCTION && n.Token() != ILLEGAL {
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

// Count the number of variables that have been used in the statement prior to the given node.
func CountVariables(n NodeI) int {
    var i int
    for n.Token() != EOL && n.Token() != ILLEGAL {
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

// Count the number of prior if declarations in an if else chain.
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

// Get the function arguments by consuming all following nodes 
// until the close of the function bracket or a { is found.
func GetArgs(n NodeI, sep string) (string, NodeI) {
    var str string
    var s string
    count := 1
    for n != nil && n.Token() != LBRACE && count > 0 {
        s, n = n.String()
        if n.Token() == LPAREN {
            count++
        } else if n.Token() == RPAREN {
            count--
        }
        if len(s) > 0 {
            str += sep + s
        }
    }
    return str, n
}
