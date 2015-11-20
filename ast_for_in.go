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

type ForIn struct {
	*Node
}

func (this ForIn) String() (string, NodeI) {
	item := this.ModulePrefix() + this.Next().Ident()
	list := this.ModulePrefix() + this.Next().Next().Next().Ident() // for->var->in->
	next := this.Next().Next().Next().Next()                        // for->item->in->list->
	return "for " + item + " in \"${" + list + "[@]}\"; do", next
}
