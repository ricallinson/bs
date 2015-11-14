package main

import (
	"fmt"
	. "github.com/ricallinson/simplebdd"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestApplication(t *testing.T) {

	Describe("File tests", func() {
		dir := "./fixtures"
		files, _ := ioutil.ReadDir(dir)
		for _, file := range files {
			if name := file.Name(); strings.HasSuffix(name, ".bs") {
				filepath := path.Join(dir, name)
				file, _ := ioutil.ReadFile(filepath)
				tests := strings.Split(string(file), "===")
				for i, test := range tests {
					It(fmt.Sprintf("%s - %d", filepath, i+1), func() {
						p := strings.Split(string(test), "---")
						s := ParseString(p[0])
						AssertEqual(strings.TrimSpace(s), strings.TrimSpace(p[1]))
					})
				}
			}
		}
	})

	Report(t)
}
