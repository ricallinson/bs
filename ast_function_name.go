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

type FunctionName struct {
	*Node
}

func (this *FunctionName) String() (string, NodeI) {
	// If the line stars with a function then print it out.
	if this.Prev().Token() == EOL || this.Prev().Token() == ILLEGAL {
		return "\"" + this.Ident() + "\" ", this.Next()
	}
	// If the function call is a value then it must be wrapped in $().
	str := "$(\"" + this.Ident() + "\""
	s, node := GetArgs(this.Next().Next(), " ") // foo->(->
	str += s
	if IsArithmetic(this.Prev()) == false && IsArithmetic(node.Next()) {
		str = "$((" + str
	} else if IsArithmetic(this.Prev()) && IsArithmetic(node.Next()) == false {
		str += "))"
	}
	return str + ")", node.Next()
}
