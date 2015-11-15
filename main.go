package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	p := NewParser(r)
	root, err := p.Parse(&Node{})
	if err != nil {
		panic(err)
	}
	fmt.Print(root.String())
}
