package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
	"io/ioutil"
	"strings"
	"path"
	"fmt"
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
					It(fmt.Sprintf("%s - %d", filepath, i + 1), func() {
						p := strings.Split(string(test), "---")
						s := ParseString(p[0])
						AssertEqual(strings.TrimSpace(s), strings.TrimSpace(p[1]))
					})
				}
			}
		}
	})

	Describe("String", func() {

		It("should return assignment of string", func() {
			s := ParseString("a = \"foo\"")
			AssertEqual(s, "a=\"foo\"\n")
		})

		It("should return assignment of strings", func() {
			s := ParseString("a = \"foo\" + \"bar\"")
			AssertEqual(s, "a=\"foo\"\"bar\"\n")
		})
	})

	Report(t)
}
