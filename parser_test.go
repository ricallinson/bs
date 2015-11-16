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
