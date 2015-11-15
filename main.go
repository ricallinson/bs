package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	p := NewParser(r)
	fmt.Print(p.ParseToString())
}
