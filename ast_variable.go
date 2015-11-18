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
	"fmt"
	"strings"
)

type Variable struct {
	*Node
}

func (this *Variable) String() (string, NodeI) {
	var str string
	// Are we in an exists function call, if so do something special.
	// TODO: Move this block out to an object.
	if this.Next().Next().Token() == EXISTS { // a->=->exists->(->"str"->)
		param, _ := this.Next().Next().Next().Next().String()
		str = "[ -e " + param + " ]\n" + this.Ident() + "=$((!$?))"
		return str, this.Next().Next().Next().Next().Next()
	}
	// Get the full variable name.
	ident, next := GetModuleIdent(this)
	// Is this variable a list with an index?
	if this.Next().Token() == LSQUARE {
		var args []string
		args, next = GetArgs(next.Next())
		index := strings.Join(args, " ")
		return "\"${" + ident + "[" + index + "]}\"", next
	}
	// If we got here then then it is a normal variable so check the previous token to see whats going on.
	switch this.Prev().Token() {
	case 0, EOL:
		// If there is no parent, a '{' or an 'EOL' then this must be the first var to be assigned.
		str = ident
	default:
		// Are there any operators that makes this variable a number?
		if IsArithmetic(this.Prev()) == false && IsArithmetic(next) && this.Prev().Token() != LSQUARE {
			str = fmt.Sprintf("$(($%s", ident)
		} else if IsArithmetic(this.Prev()) && IsArithmetic(next) == false && next.Token() != RSQUARE {
			str = fmt.Sprintf("$%s))", ident)
		} else if IsArithmetic(this.Prev()) && IsArithmetic(next) {
			str = fmt.Sprintf("$%s", ident)
		} else {
			// Otherwise assume the var needs to be wrapped as a string in '""'.
			str = fmt.Sprintf("\"$%s\"", ident)
		}
	}
	return str, next
}
