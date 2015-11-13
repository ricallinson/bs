package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	p := NewParser(r)
	s, e := p.Parse(&Node{})
	if e != nil {
		panic(e)
	}
	for _, l := range s {
		fmt.Print(l.String())
	}
}
