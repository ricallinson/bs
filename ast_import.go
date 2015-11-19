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

type Import struct {
	*Node
}

func (this *Import) String() (string, NodeI) {
	var str string
	path, node := GetArgs(this.Next()) // len->(
	if len(path) == 1 {
		str = ParseDir(strings.Trim(path[0], "\""), this.Node.parser.module, this.Node.parser.funcs)
	}
	return str, node
}
