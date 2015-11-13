package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
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

		It("should return simple function", func() {
			s := ParseString("func foo () {return 1}")
			AssertEqual(s, "function foo {\n\"echo\" \"-ne\" $((1))\nreturn\n}\n")
		})

		It("should return simple function with parameters", func() {
			s := ParseString("func foo (a b) {return 1}")
			AssertEqual(s, "function foo {\n\"echo\" \"-ne\" $((1))\nreturn\n}\n")
		})
	})

	Report(t)
}
