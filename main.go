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
