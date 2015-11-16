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
)

type Block struct {
	*Node
	block NodeI
}

// Look back from the given node and return a 'function' node if the given node is with a 'function' block.
func (this Block) GetFunction(n NodeI) NodeI {
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
func (this Block) GetIf(n NodeI) NodeI {
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
func (this Block) GetWhile(n NodeI) NodeI {
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
func (this Block) GetForIn(n NodeI) NodeI {
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

func (this Block) GetArguments(n NodeI) string {
	// Get the arguments from the function.
	var strArgs string
	var strLocal string
	argsList := map[string]bool{}
	args := n.Next().Next().Next() // function->foo->(
	i := 1                         // Count the given arguments.
	for args.Token() != LBRACE {   // function->foo->(->...)->{<-
		if args.Token() == IDENT {
			argsList[args.Ident()] = true
			strArgs += fmt.Sprintf("\nlocal %s\n%s=\"$%d\"", args.Ident(), args.Ident(), i)
			i++
		}
		args = args.Next()
	}
	// Search the block for variables and then check if they should be local.
	_, local := this.block.String()
	for local != nil {
		if local.Token() == IDENT {
			if ok := argsList[local.Ident()]; ok == false {
				strLocal += "\nlocal " + local.Ident()
			}
		}
		_, local = local.String()
	}
	return strLocal + strArgs
}

// Look back to see if the given node was in a function with a return.
func (this Block) GetReturn(n NodeI) bool {
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

// Print each statement in the block as a string.
func (this Block) BlockString() string {
	var str string
	s, line := this.block.String()
	str += s
	for line != nil {
		s, line = line.String()
		str += s
	}
	return str
}

// Count the number of prior 'if' declarations in an 'if else' chain.
func (this Block) CountIfs(n NodeI) int {
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

func (this Block) IfTerminators() string {
	var str string
	if this.Next().Token() != ELSE {
		// We are at the end of an 'if' block so count the total 'if' statements and then print an 'fi' for each .
		for i := 0; i < this.CountIfs(this); i++ {
			str += "fi\n"
		}
		str = str[:len(str)-1]
	}
	return str
}

func (this Block) String() (string, NodeI) {
	var str string
	// Check back to see if we are in a function, if we are then print out the function arguments here.
	if n := this.GetFunction(this); n != nil {
		str += " {" + this.GetArguments(n)
	}
	// Print each statement in the block.
	str += this.BlockString()
	// Was there a return in this block?
	if this.GetReturn(this.block) {
		str += "return\n"
	}
	// Check to see what type of block we are in.
	if this.GetFunction(this) != nil {
		// If we are at the end of a function then print return.
		str += "}"
	} else if this.GetIf(this) != nil {
		str += this.IfTerminators()
	} else if this.GetWhile(this) != nil || this.GetForIn(this) != nil {
		// If we are at the end of a while block then print done.
		str += "done\n"
	}
	return str, this.Next()
}
