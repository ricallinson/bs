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
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	var script = flag.Bool("script", false, "print the parsed script")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("You must provide a source file")
		return
	}
	source := flag.Arg(0)
	s := ParseFile(source)
	if *script {
		fmt.Print(s)
		return
	}
	out, err := exec.Command("bash", "-c", s).Output()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Print(string(out))
}
