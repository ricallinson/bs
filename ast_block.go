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

func (this Block) String() (string, NodeI) {
    var str string
    // Check back to see if we are in a function, if we are then print out the function arguments here.
    if n := GetMyFunction(this); n != nil {
        str += " {"
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
    } else if Find(WHILE, this, PREV) != nil || Find(FORIN, this, PREV) != nil {
        // If we are at the end of a while block then print done.
        str += "done\n"
    }
    return str, this.Next()
}
