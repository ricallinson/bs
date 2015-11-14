package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
	"io/ioutil"
)

func TestApplication(t *testing.T) {

	Describe("Variable", func() {

		It("should return assignment of var", func() {
			s := ParseString("a = b")
			AssertEqual(s, "a=\"$b\"\n")
		})

		It("should return assignment of vars", func() {
			s := ParseString("a = b + c")
			AssertEqual(s, "a=\"$b\"\"$c\"\n")
		})
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

	Describe("Integer", func() {

		It("should return assignment of integer", func() {
			s := ParseString("a = 1")
			AssertEqual(s, "a=$((1))\n")
		})

		It("should return assignment of integer add", func() {
			s := ParseString("a = 1 + 1")
			AssertEqual(s, "a=$((1 + 1))\n")
		})

		It("should return assignment of integer multiply", func() {
			s := ParseString("a = 1 * 1")
			AssertEqual(s, "a=$((1 * 1))\n")
		})

		It("should return assignment of integer subtract", func() {
			s := ParseString("a = 1 - 1")
			AssertEqual(s, "a=$((1 - 1))\n")
		})

		It("should return assignment of integer divide", func() {
			s := ParseString("a = 1 / 1")
			AssertEqual(s, "a=$((1 / 1))\n")
		})
	})

	Describe("Function", func() {

		It("should test func_001", func() {
			s := ParseFile("./fixtures/func_001.go")
			t, _ := ioutil.ReadFile("./fixtures/func_001.bash")
			AssertEqual(s, string(t))
		})

		It("should test func_002", func() {
			s := ParseFile("./fixtures/func_002.go")
			t, _ := ioutil.ReadFile("./fixtures/func_002.bash")
			AssertEqual(s, string(t))
		})
	})

	Report(t)
}
