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

package spew

import (
	"fmt"
	"io"
)

// Errorf is a wrapper for fmt.Errorf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the formatted string as a value that satisfies error.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Errorf(format, spew.NewFormatter(a), spew.NewFormatter(b))
func Errorf(format string, a ...any) (err error) {
	return fmt.Errorf(format, Default.convertArgs(a)...)
}

// Fprint is a wrapper for fmt.Fprint that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprint(w, spew.NewFormatter(a), spew.NewFormatter(b))
func Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, Default.convertArgs(a)...)
}

// Fprintf is a wrapper for fmt.Fprintf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprintf(w, format, spew.NewFormatter(a), spew.NewFormatter(b))
func Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return fmt.Fprintf(w, format, Default.convertArgs(a)...)
}

// Fprintln is a wrapper for fmt.Fprintln that treats each argument as if it
// passed with a default Formatter interface returned by NewFormatter.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprintln(w, spew.NewFormatter(a), spew.NewFormatter(b))
func Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprintln(w, Default.convertArgs(a)...)
}

// Print is a wrapper for fmt.Print that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Print(spew.NewFormatter(a), spew.NewFormatter(b))
func Print(a ...any) (n int, err error) {
	return fmt.Print(Default.convertArgs(a)...)
}

// Printf is a wrapper for fmt.Printf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Printf(format, spew.NewFormatter(a), spew.NewFormatter(b))
func Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, Default.convertArgs(a)...)
}

// Println is a wrapper for fmt.Println that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Println(spew.NewFormatter(a), spew.NewFormatter(b))
func Println(a ...any) (n int, err error) {
	return fmt.Println(Default.convertArgs(a)...)
}

// Sprint is a wrapper for fmt.Sprint that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the resulting string.  See NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Sprint(spew.NewFormatter(a), spew.NewFormatter(b))
func Sprint(a ...any) string {
	return fmt.Sprint(Default.convertArgs(a)...)
}

// Sprintf is a wrapper for fmt.Sprintf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the resulting string.  See NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Sprintf(format, spew.NewFormatter(a), spew.NewFormatter(b))
func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, Default.convertArgs(a)...)
}

// Sprintln is a wrapper for fmt.Sprintln that treats each argument as if it
// were passed with a default Formatter interface returned by NewFormatter.  It
// returns the resulting string.  See NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Sprintln(spew.NewFormatter(a), spew.NewFormatter(b))
func Sprintln(a ...any) string {
	return fmt.Sprintln(Default.convertArgs(a)...)
}
