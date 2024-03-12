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
// Config a test applies to.
type spewFunc int

const (
	fnConfigFdump spewFunc = iota
	fnConfigFprint
	fnConfigFprintf
	fnConfigFprintln
	fnConfigPrint
	fnConfigPrintln
	fnConfigSdump
	fnConfigSprint
	fnConfigSprintf
	fnConfigSprintln
	fnConfigErrorf
	fnConfigNewFormatter
	fnErrorf
	fnFprint
	fnFprintln
	fnPrint
	fnPrintln
	fnSdump
	fnSprint
	fnSprintf
	fnSprintln
)

// Map of spewFunc values to names for pretty printing.
var spewFuncStrings = map[spewFunc]string{
	fnConfigFdump:        "Config.Fdump",
	fnConfigFprint:       "Config.Fprint",
	fnConfigFprintf:      "Config.Fprintf",
	fnConfigFprintln:     "Config.Fprintln",
	fnConfigSdump:        "Config.Sdump",
	fnConfigPrint:        "Config.Print",
	fnConfigPrintln:      "Config.Println",
	fnConfigSprint:       "Config.Sprint",
	fnConfigSprintf:      "Config.Sprintf",
	fnConfigSprintln:     "Config.Sprintln",
	fnConfigErrorf:       "Config.Errorf",
	fnConfigNewFormatter: "Config.NewFormatter",
	fnErrorf:             "spew.Errorf",
	fnFprint:             "spew.Fprint",
	fnFprintln:           "spew.Fprintln",
	fnPrint:              "spew.Print",
	fnPrintln:            "spew.Println",
	fnSdump:              "spew.Sdump",
	fnSprint:             "spew.Sprint",
	fnSprintf:            "spew.Sprintf",
	fnSprintln:           "spew.Sprintln",
}

func (f spewFunc) String() string {
	if s, ok := spewFuncStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown spewFunc (%d)", int(f))
}

// spewTest is used to describe a test to be performed against the public
// functions of the spew package or Config.
type spewTest struct {
	line   string // use line() to fill this
	cfg    *spew.Config
	f      spewFunc
	format string
	in     any
	want   string
}

// spewTests houses the tests to be performed against the public functions of
// the spew package and Config.
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
	cfgDefault := spew.NewDefaultConfig()
	cfgNoMethods := &spew.Config{Indent: " ", DisableMethods: true}
	cfgNoPmethods := &spew.Config{Indent: " ", DisablePointerMethods: true}
	cfgMaxDepth := &spew.Config{Indent: " ", MaxDepth: 1}
	cfgContinue := &spew.Config{Indent: " ", ContinueOnMethod: true}
	cfgNoPtrAddr := &spew.Config{DisablePointerAddresses: true}
	cfgNoCap := &spew.Config{DisableCapacities: true}
	cfgTrailingComma := &spew.Config{Indent: " ", TrailingCommas: true}
	cfgNoUnexported := &spew.Config{Indent: " ", DisableUnexported: true}
	cfgQuotes := &spew.Config{QuoteStrings: true}
	cfgClean := &spew.CleanConfig

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
		slice []any
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
		{line(), cfgDefault, fnConfigFdump, "", int8(127), "(int8) 127\n"},
		{line(), cfgDefault, fnConfigFprint, "", int16(32767), "32767"},
		{line(), cfgDefault, fnConfigFprintf, "%v", int32(2147483647), "2147483647"},
		{line(), cfgDefault, fnConfigFprintln, "", int(2147483647), "2147483647\n"},
		{line(), cfgDefault, fnConfigPrint, "", int64(9223372036854775807), "9223372036854775807"},
		{line(), cfgDefault, fnConfigPrintln, "", uint8(255), "255\n"},
		{line(), cfgDefault, fnConfigSdump, "", uint8(64), "(uint8) 64\n"},
		{line(), cfgDefault, fnConfigSprint, "", complex(1, 2), "(1+2i)"},
		{line(), cfgDefault, fnConfigSprintf, "%v", complex(float32(3), 4), "(3+4i)"},
		{line(), cfgDefault, fnConfigSprintln, "", complex(float64(5), 6), "(5+6i)\n"},
		{line(), cfgDefault, fnConfigErrorf, "%#v", uint16(65535), "(uint16)65535"},
		{line(), cfgDefault, fnConfigNewFormatter, "%v", uint32(4294967295), "4294967295"},
		{line(), cfgDefault, fnErrorf, "%v", uint64(18446744073709551615), "18446744073709551615"},
		{line(), cfgDefault, fnFprint, "", float32(3.14), "3.14"},
		{line(), cfgDefault, fnFprintln, "", float64(6.28), "6.28\n"},
		{line(), cfgDefault, fnPrint, "", true, "true"},
		{line(), cfgDefault, fnPrintln, "", false, "false\n"},
		{line(), cfgDefault, fnSdump, "", complex(-10, -20), "(complex128) (-10-20i)\n"},
		{line(), cfgDefault, fnSprint, "", complex(-1, -2), "(-1-2i)"},
		{line(), cfgDefault, fnSprintf, "%v", complex(float32(-3), -4), "(-3-4i)"},
		{line(), cfgDefault, fnSprintln, "", complex(float64(-5), -6), "(-5-6i)\n"},
		{line(), cfgNoMethods, fnConfigFprint, "", ts, "test"},
		{line(), cfgNoMethods, fnConfigFprint, "", &ts, "<*>test"},
		{line(), cfgNoMethods, fnConfigFprint, "", tps, "test"},
		{line(), cfgNoMethods, fnConfigFprint, "", &tps, "<*>test"},
		{line(), cfgNoPmethods, fnConfigFprint, "", ts, "stringer test"},
		{line(), cfgNoPmethods, fnConfigFprint, "", &ts, "<*>stringer test"},
		{line(), cfgNoPmethods, fnConfigFprint, "", tps, "test"},
		{line(), cfgNoPmethods, fnConfigFprint, "", &tps, "<*>stringer test"},
		{line(), cfgMaxDepth, fnConfigFprint, "", dt, "{{<max>} [<max>] [<max>] map[<max>]}"},
		{line(), cfgMaxDepth, fnConfigFdump, "", dt, "(spew_test.depthTester) {\n" +
			" ic: (spew_test.indirCir1) {\n  <max depth reached>\n },\n" +
			" arr: ([1]string) (len=1 cap=1) {\n  <max depth reached>\n },\n" +
			" slice: ([]string) (len=1 cap=1) {\n  <max depth reached>\n },\n" +
			" m: (map[string]int) (len=1) {\n  <max depth reached>\n }\n}\n"},
		{line(), cfgContinue, fnConfigFprint, "", ts, "(stringer test) test"},
		{line(), cfgContinue, fnConfigFdump, "", ts, "(spew_test.stringer) " +
			"(len=4) (stringer test) \"test\"\n"},
		{line(), cfgContinue, fnConfigFprint, "", te, "(error: 10) 10"},
		{line(), cfgContinue, fnConfigFdump, "", te, "(spew_test.customError) " +
			"(error: 10) 10\n"},
		{line(), cfgNoPtrAddr, fnConfigFprint, "", tptr, "<*>{<*>{}}"},
		{line(), cfgNoPtrAddr, fnConfigSdump, "", tptr, "(*spew_test.ptrTester)({\ns: (*struct {})({\n})\n})\n"},
		{line(), cfgNoCap, fnConfigSdump, "", make([]string, 0, 10), "([]string) {\n}\n"},
		{line(), cfgNoCap, fnConfigSdump, "", make([]string, 1, 10), "([]string) (len=1) {\n(string) \"\"\n}\n"},
		{line(), cfgTrailingComma, fnConfigFdump, "", commaTester{
			slice: []any{
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
		{line(), cfgNoUnexported, fnConfigSdump, "", tunexp, "(struct { X int; y int }) {\n X: (int) 123,\n}\n"},
		{line(), cfgNoUnexported, fnConfigSprintln, "", tunexp, "{123}\n"},
		{line(), cfgNoUnexported, fnConfigSprintf, "%v", tunexp, "{123}"},
		{line(), cfgNoUnexported, fnConfigSprintf, "%#v", tunexp, "(struct { X int; y int }){X:(int)123}"},
		{line(), cfgQuotes, fnConfigSdump, "", ts, "(spew_test.stringer) (len=4) \"stringer test\"\n"},
		{line(), cfgQuotes, fnConfigSprintln, "", ts, "\"stringer test\"\n"},
		{line(), cfgQuotes, fnConfigSprintf, "%v", ts, `"stringer test"`},
		{line(), cfgQuotes, fnConfigSprintf, "%#v", ts, `(spew_test.stringer)"stringer test"`},
		{line(), cfgClean, fnConfigSdump, "", make([]string, 0, 10), "[]\n"},
		{line(), cfgClean, fnConfigSdump, "", make([]string, 2, 10), "[\n  \"\",\n  \"\"\n]\n"},
		{line(), cfgClean, fnConfigSprintln, "", make([]int, 2, 10), "[0,0]\n"},
		{line(), cfgClean, fnConfigSprintf, "%v", make([]int, 2, 10), "[0,0]"},
		{line(), cfgClean, fnConfigSprintf, "%#v", make([]int, 2, 10), "([]int)[0,0]"},
		{line(), cfgClean, fnConfigSprintln, "", make([]string, 1, 10), "[\"\"]\n"},
		{line(), cfgClean, fnConfigSprintf, "%v", make([]string, 1, 10), `[""]`},
		{line(), cfgClean, fnConfigSprintf, "%#v", make([]string, 1, 10), `([]string)[""]`},
		{line(), cfgClean, fnConfigSdump, "", TestSpew,
			fmt.Sprintf("spew_test.TestSpew[spew_test.go:%d]\n", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), cfgClean, fnConfigSprintln, "", TestSpew,
			fmt.Sprintf("spew_test.TestSpew[spew_test.go:%d]\n", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), cfgClean, fnConfigSprintf, "%v", TestSpew,
			fmt.Sprintf("spew_test.TestSpew[spew_test.go:%d]", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), cfgClean, fnConfigSprintf, "%#v", TestSpew,
			fmt.Sprintf("(func(*testing.T))spew_test.TestSpew[spew_test.go:%d]", funcLine(reflect.ValueOf(TestSpew).Pointer()))},
		{line(), cfgClean, fnConfigSprintln, "", tfn,
			fmt.Sprintf("spew_test.initSpewTests.func1[spew_test.go:%d]\n", funcLine(reflect.ValueOf(tfn).Pointer()))},
		{line(), cfgClean, fnConfigSprintf, "%v", tfn,
			fmt.Sprintf("spew_test.initSpewTests.func1[spew_test.go:%d]", funcLine(reflect.ValueOf(tfn).Pointer()))},
		{line(), cfgClean, fnConfigSprintf, "%#v", tfn,
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
		case fnConfigFdump:
			test.cfg.Fdump(buf, test.in)

		case fnConfigFprint:
			test.cfg.Fprint(buf, test.in)

		case fnConfigFprintf:
			test.cfg.Fprintf(buf, test.format, test.in)

		case fnConfigFprintln:
			test.cfg.Fprintln(buf, test.in)

		case fnConfigPrint:
			b, err := redirStdout(func() { test.cfg.Print(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fnConfigPrintln:
			b, err := redirStdout(func() { test.cfg.Println(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fnConfigSdump:
			str := test.cfg.Sdump(test.in)
			buf.WriteString(str)

		case fnConfigSprint:
			str := test.cfg.Sprint(test.in)
			buf.WriteString(str)

		case fnConfigSprintf:
			str := test.cfg.Sprintf(test.format, test.in)
			buf.WriteString(str)

		case fnConfigSprintln:
			str := test.cfg.Sprintln(test.in)
			buf.WriteString(str)

		case fnConfigErrorf:
			err := test.cfg.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fnConfigNewFormatter:
			fmt.Fprintf(buf, test.format, test.cfg.NewFormatter(test.in))

		case fnErrorf:
			err := spew.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fnFprint:
			spew.Fprint(buf, test.in)

		case fnFprintln:
			spew.Fprintln(buf, test.in)

		case fnPrint:
			b, err := redirStdout(func() { spew.Print(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fnPrintln:
			b, err := redirStdout(func() { spew.Println(test.in) })
			if err != nil {
				t.Errorf("line %s: %v %v", test.line, test.f, err)
				continue
			}
			buf.Write(b)

		case fnSdump:
			str := spew.Sdump(test.in)
			buf.WriteString(str)

		case fnSprint:
			str := spew.Sprint(test.in)
			buf.WriteString(str)

		case fnSprintf:
			str := spew.Sprintf(test.format, test.in)
			buf.WriteString(str)

		case fnSprintln:
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
