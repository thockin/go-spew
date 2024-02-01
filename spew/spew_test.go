/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package spew_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/thockin/go-spew/spew"
)

// spewFunc is used to identify which public function of the spew package or
// ConfigState a test applies to.
type spewFunc int

const (
	fCSFdump spewFunc = iota
	fCSFprint
	fCSFprintf
	fCSFprintln
	fCSPrint
	fCSPrintln
	fCSSdump
	fCSSprint
	fCSSprintf
	fCSSprintln
	fCSErrorf
	fCSNewFormatter
	fErrorf
	fFprint
	fFprintln
	fPrint
	fPrintln
	fSdump
	fSprint
	fSprintf
	fSprintln
)

// Map of spewFunc values to names for pretty printing.
var spewFuncStrings = map[spewFunc]string{
	fCSFdump:        "ConfigState.Fdump",
	fCSFprint:       "ConfigState.Fprint",
	fCSFprintf:      "ConfigState.Fprintf",
	fCSFprintln:     "ConfigState.Fprintln",
	fCSSdump:        "ConfigState.Sdump",
	fCSPrint:        "ConfigState.Print",
	fCSPrintln:      "ConfigState.Println",
	fCSSprint:       "ConfigState.Sprint",
	fCSSprintf:      "ConfigState.Sprintf",
	fCSSprintln:     "ConfigState.Sprintln",
	fCSErrorf:       "ConfigState.Errorf",
	fCSNewFormatter: "ConfigState.NewFormatter",
	fErrorf:         "spew.Errorf",
	fFprint:         "spew.Fprint",
	fFprintln:       "spew.Fprintln",
	fPrint:          "spew.Print",
	fPrintln:        "spew.Println",
	fSdump:          "spew.Sdump",
	fSprint:         "spew.Sprint",
	fSprintf:        "spew.Sprintf",
	fSprintln:       "spew.Sprintln",
}

func (f spewFunc) String() string {
	if s, ok := spewFuncStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown spewFunc (%d)", int(f))
}

// spewTest is used to describe a test to be performed against the public
// functions of the spew package or ConfigState.
type spewTest struct {
	line   string // use line() to fill this
	cs     *spew.ConfigState
	f      spewFunc
	format string
	in     interface{}
	want   string
}

// spewTests houses the tests to be performed against the public functions of
// the spew package and ConfigState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var spewTests []spewTest

// redirStdout is a helper function to return the standard output from f as a
// byte slice.
func redirStdout(f func()) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "ss-test")
	if err != nil {
		return nil, err
	}
	fileName := tempFile.Name()
	defer os.Remove(fileName) // Ignore error

	origStdout := os.Stdout
	os.Stdout = tempFile
	f()
	os.Stdout = origStdout
	tempFile.Close()

	return os.ReadFile(fileName)
}

func initSpewTests() {
	// Config states with various settings.
	scsDefault := spew.NewDefaultConfig()
	scsNoMethods := &spew.ConfigState{Indent: " ", DisableMethods: true}
	scsNoPmethods := &spew.ConfigState{Indent: " ", DisablePointerMethods: true}
	scsMaxDepth := &spew.ConfigState{Indent: " ", MaxDepth: 1}
	scsContinue := &spew.ConfigState{Indent: " ", ContinueOnMethod: true}
	scsNoPtrAddr := &spew.ConfigState{DisablePointerAddresses: true}
	scsNoCap := &spew.ConfigState{DisableCapacities: true}
	scsTrailingComma := &spew.ConfigState{Indent: " ", TrailingCommas: true}
	scsNoUnexported := &spew.ConfigState{Indent: " ", DisableUnexported: true}
	scsQuotes := &spew.ConfigState{QuoteStrings: true}
	scsClean := &spew.CleanConfig

	// Variables for tests on types which implement Stringer interface with and
	// without a pointer receiver.
	ts := stringer("test")
	tps := pstringer("test")

	type ptrTester struct {
		s *struct{}
	}
	tptr := &ptrTester{s: &struct{}{}}

	// depthTester is used to test max depth handling for structs, array, slices
	// and maps.
	type depthTester struct {
		ic    indirCir1
		arr   [1]string
		slice []string
		m     map[string]int
	}
	dt := depthTester{indirCir1{nil}, [1]string{"arr"}, []string{"slice"},
		map[string]int{"one": 1}}

	// commatester is to test trailing commas
	type commaTester struct {
		slice []interface{}
		m     map[string]int
	}

	// Variable for tests on types which implement error interface.
	te := customError(10)

	// unexported fields.
	tunexp := struct {
		X int
		y int
	}{123, 456}

	// Variable for tests on anonymous functions.
	tfn := func() {}

	spewTests = []spewTest{
		{line(), scsDefault, fCSFdump, "", int8(127), "(int8) 127\n"},
		{line(), scsDefault, fCSFprint, "", int16(32767), "32767"},
		{line(), scsDefault, fCSFprintf, "%v", int32(2147483647), "2147483647"},
		{line(), scsDefault, fCSFprintln, "", int(2147483647), "2147483647\n"},
		{line(), scsDefault, fCSPrint, "", int64(9223372036854775807), "9223372036854775807"},
		{line(), scsDefault, fCSPrintln, "", uint8(255), "255\n"},
		{line(), scsDefault, fCSSdump, "", uint8(64), "(uint8) 64\n"},
		{line(), scsDefault, fCSSprint, "", complex(1, 2), "(1+2i)"},
		{line(), scsDefault, fCSSprintf, "%v", complex(float32(3), 4), "(3+4i)"},
		{line(), scsDefault, fCSSprintln, "", complex(float64(5), 6), "(5+6i)\n"},
		{line(), scsDefault, fCSErrorf, "%#v", uint16(65535), "(uint16)65535"},
		{line(), scsDefault, fCSNewFormatter, "%v", uint32(4294967295), "4294967295"},
		{line(), scsDefault, fErrorf, "%v", uint64(18446744073709551615), "18446744073709551615"},
		{line(), scsDefault, fFprint, "", float32(3.14), "3.14"},
		{line(), scsDefault, fFprintln, "", float64(6.28), "6.28\n"},
		{line(), scsDefault, fPrint, "", true, "true"},
		{line(), scsDefault, fPrintln, "", false, "false\n"},
		{line(), scsDefault, fSdump, "", complex(-10, -20), "(complex128) (-10-20i)\n"},
		{line(), scsDefault, fSprint, "", complex(-1, -2), "(-1-2i)"},
		{line(), scsDefault, fSprintf, "%v", complex(float32(-3), -4), "(-3-4i)"},
		{line(), scsDefault, fSprintln, "", complex(float64(-5), -6), "(-5-6i)\n"},
		{line(), scsNoMethods, fCSFprint, "", ts, "test"},
		{line(), scsNoMethods, fCSFprint, "", &ts, "<*>test"},
		{line(), scsNoMethods, fCSFprint, "", tps, "test"},
		{line(), scsNoMethods, fCSFprint, "", &tps, "<*>test"},
		{line(), scsNoPmethods, fCSFprint, "", ts, "stringer test"},
		{line(), scsNoPmethods, fCSFprint, "", &ts, "<*>stringer test"},
		{line(), scsNoPmethods, fCSFprint, "", tps, "test"},
		{line(), scsNoPmethods, fCSFprint, "", &tps, "<*>stringer test"},
		{line(), scsMaxDepth, fCSFprint, "", dt, "{{<max>} [<max>] [<max>] map[<max>]}"},
		{line(), scsMaxDepth, fCSFdump, "", dt, "(spew_test.depthTester) {\n" +
			" ic: (spew_test.indirCir1) {\n  <max depth reached>\n },\n" +
			" arr: ([1]string) (len=1 cap=1) {\n  <max depth reached>\n },\n" +
			" slice: ([]string) (len=1 cap=1) {\n  <max depth reached>\n },\n" +
			" m: (map[string]int) (len=1) {\n  <max depth reached>\n }\n}\n"},
		{line(), scsContinue, fCSFprint, "", ts, "(stringer test) test"},
		{line(), scsContinue, fCSFdump, "", ts, "(spew_test.stringer) " +
			"(len=4) (stringer test) \"test\"\n"},
		{line(), scsContinue, fCSFprint, "", te, "(error: 10) 10"},
		{line(), scsContinue, fCSFdump, "", te, "(spew_test.customError) " +
			"(error: 10) 10\n"},
		{line(), scsNoPtrAddr, fCSFprint, "", tptr, "<*>{<*>{}}"},
		{line(), scsNoPtrAddr, fCSSdump, "", tptr, "(*spew_test.ptrTester)({\ns: (*struct {})({\n})\n})\n"},
		{line(), scsNoCap, fCSSdump, "", make([]string, 0, 10), "([]string) {\n}\n"},
		{line(), scsNoCap, fCSSdump, "", make([]string, 1, 10), "([]string) (len=1) {\n(string) \"\"\n}\n"},
		{line(), scsTrailingComma, fCSFdump, "", commaTester{
			slice: []interface{}{
				map[string]int{"one": 1},
			},
			m: map[string]int{"one": 1},
		},
			"(spew_test.commaTester) {\n" +
				" slice: ([]interface {}) (len=1 cap=1) {\n" +
				"  (map[string]int) (len=1) {\n" +
				"   (string) (len=3) \"one\": (int) 1,\n" +
				"  },\n" +
				" },\n" +
				" m: (map[string]int) (len=1) {\n" +
				"  (string) (len=3) \"one\": (int) 1,\n" +
				" },\n" +
				"}\n"},
		{line(), scsNoUnexported, fCSSdump, "", tunexp, "(struct { X int; y int }) {\n X: (int) 123,\n}\n"},
		{line(), scsNoUnexported, fCSSprintln, "", tunexp, "{123}\n"},
		{line(), scsNoUnexported, fCSSprintf, "%v", tunexp, "{123}"},
		{line(), scsNoUnexported, fCSSprintf, "%#v", tunexp, "(struct { X int; y int }){X:(int)123}"},
		{line(), scsQuotes, fCSSdump, "", ts, "(spew_test.stringer) (len=4) \"stringer test\"\n"},
		{line(), scsQuotes, fCSSprintln, "", ts, "\"stringer test\"\n"},
		{line(), scsQuotes, fCSSprintf, "%v", ts, `"stringer test"`},
		{line(), scsQuotes, fCSSprintf, "%#v", ts, `(spew_test.stringer)"stringer test"`},
		{line(), scsClean, fCSSdump, "", make([]string, 0, 10), "[]\n"},
		{line(), scsClean, fCSSdump, "", make([]string, 2, 10), "[\n  \"\",\n  \"\"\n]\n"},
		{line(), scsClean, fCSSprintln, "", make([]int, 2, 10), "[0,0]\n"},
		{line(), scsClean, fCSSprintf, "%v", make([]int, 2, 10), "[0,0]"},
		{line(), scsClean, fCSSprintf, "%#v", make([]int, 2, 10), "([]int)[0,0]"},
		{line(), scsClean, fCSSprintln, "", make([]string, 1, 10), "[\"\"]\n"},
		{line(), scsClean, fCSSprintf, "%v", make([]string, 1, 10), `[""]`},
		{line(), scsClean, fCSSprintf, "%#v", make([]string, 1, 10), `([]string)[""]`},
		{line(), scsClean, fCSSdump, "", TestSpew,
			fmt.Sprintf("spew_test.TestSpew[spew_test.go:%d]\n", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), scsClean, fCSSprintln, "", TestSpew,
			fmt.Sprintf("spew_test.TestSpew[spew_test.go:%d]\n", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), scsClean, fCSSprintf, "%v", TestSpew,
			fmt.Sprintf("spew_test.TestSpew[spew_test.go:%d]", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), scsClean, fCSSprintf, "%#v", TestSpew,
			fmt.Sprintf("(func(*testing.T))spew_test.TestSpew[spew_test.go:%d]", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), scsClean, fCSSprintln, "", tfn,
			fmt.Sprintf("spew_test.initSpewTests.func1[spew_test.go:%d]\n", funcLine(reflect.ValueOf(tfn).Pointer()))},
		{line(), scsClean, fCSSprintf, "%v", tfn,
			fmt.Sprintf("spew_test.initSpewTests.func1[spew_test.go:%d]", funcLine(reflect.ValueOf(tfn).Pointer()))},
		{line(), scsClean, fCSSprintf, "%#v", tfn,
			fmt.Sprintf("(func())spew_test.initSpewTests.func1[spew_test.go:%d]", funcLine(reflect.ValueOf(tfn).Pointer()))},
	}
}

func funcLine(p uintptr) int {
	fn := runtime.FuncForPC(p)
	if fn == nil {
		return -1
	}
	_, line := fn.FileLine(p)
	return line
}

// TestSpew executes all of the tests described by spewTests.
func TestSpew(t *testing.T) {
	initSpewTests()

	t.Logf("Running %d tests", len(spewTests))
	for _, test := range spewTests {
		buf := new(bytes.Buffer)
		switch test.f {
		case fCSFdump:
			test.cs.Fdump(buf, test.in)

		case fCSFprint:
			test.cs.Fprint(buf, test.in)

		case fCSFprintf:
			test.cs.Fprintf(buf, test.format, test.in)

		case fCSFprintln:
			test.cs.Fprintln(buf, test.in)

		case fCSPrint:
			b, err := redirStdout(func() { test.cs.Print(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fCSPrintln:
			b, err := redirStdout(func() { test.cs.Println(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fCSSdump:
			str := test.cs.Sdump(test.in)
			buf.WriteString(str)

		case fCSSprint:
			str := test.cs.Sprint(test.in)
			buf.WriteString(str)

		case fCSSprintf:
			str := test.cs.Sprintf(test.format, test.in)
			buf.WriteString(str)

		case fCSSprintln:
			str := test.cs.Sprintln(test.in)
			buf.WriteString(str)

		case fCSErrorf:
			err := test.cs.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fCSNewFormatter:
			fmt.Fprintf(buf, test.format, test.cs.NewFormatter(test.in))

		case fErrorf:
			err := spew.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fFprint:
			spew.Fprint(buf, test.in)

		case fFprintln:
			spew.Fprintln(buf, test.in)

		case fPrint:
			b, err := redirStdout(func() { spew.Print(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fPrintln:
			b, err := redirStdout(func() { spew.Println(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fSdump:
			str := spew.Sdump(test.in)
			buf.WriteString(str)

		case fSprint:
			str := spew.Sprint(test.in)
			buf.WriteString(str)

		case fSprintf:
			str := spew.Sprintf(test.format, test.in)
			buf.WriteString(str)

		case fSprintln:
			str := spew.Sprintln(test.in)
			buf.WriteString(str)

		default:
			t.Errorf("line %s: %v unrecognized function", test.line, test.f)
			continue
		}
		s := buf.String()
		if test.want != s {
			nl := ""
			if !strings.HasSuffix(s, "\n") {
				nl = "\n"
			}
			t.Errorf("testcase on line %s:\n got: %s%swant: %s", test.line, s, nl, test.want)
			continue
		}
	}
}
