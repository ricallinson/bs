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

type FunctionName struct {
	*Node
}

func (this *FunctionName) String() (string, NodeI) {
	// Get the full variable name.
	ident, next := GetModuleIdent(this)
	// If the line starts with the function call then print it out.
	if this.Prev().Token() == EOL || this.Prev().Token() == ILLEGAL {
		return "\"" + ident + "\" ", next
	}
	// Otherwise the function call is a value and it must be wrapped in '$("")'.
	str := "$(\"" + ident + "\""
	var args []string
	args, next = GetArgs(next) // foo->(->
	if len(args) > 0 {
		str += " "
	}
	str += strings.Join(args, " ") + ")"
	// If the function call is part of some arithmetic then it must be wrapped again in '$(())'.
	if IsArithmetic(this.Prev()) == false && IsArithmetic(next.Next()) {
		str = "$((" + str
	} else if IsArithmetic(this.Prev()) && IsArithmetic(next.Next()) == false {
		str += "))"
	}
	return str, next
}
