/*
Copyright 2022 The logr Authors.

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

package spew_test

import (
	"testing"

	"github.com/thockin/go-spew/spew"
)

//go:noinline
func doSprintOneArg(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		spew.Sprint("one arg")
	}
}

//go:noinline
func doSprintManyArgs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		spew.Sprint(true, "str", 42, 3.14,
			struct{ X, Y int }{93, 76},
			&struct{ X, Y int }{118, 78},
			[]int{8, 6, 7, 5, 3, 0, 9},
			map[string]int{"one": 1, "two": 2})
	}
}

//go:noinline
func doSprintfOneArgV(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		spew.Sprint("%v", "one arg")
	}
}

//go:noinline
func doSprintfManyArgsV(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		spew.Sprintf("%v %v %v %v %v %v %v %v",
			true, "str", 42, 3.14,
			struct{ X, Y int }{93, 76},
			&struct{ X, Y int }{118, 78},
			[]int{8, 6, 7, 5, 3, 0, 9},
			map[string]int{"one": 1, "two": 2})
	}
}

func BenchmarkSprintOneArg(b *testing.B) {
	doSprintOneArg(b)
}

func BenchmarkSprintManyArgs(b *testing.B) {
	doSprintManyArgs(b)
}

func BenchmarkSprintfOneArgV(b *testing.B) {
	doSprintOneArg(b)
}

func BenchmarkSprintfManyArgsV(b *testing.B) {
	doSprintManyArgs(b)
}
