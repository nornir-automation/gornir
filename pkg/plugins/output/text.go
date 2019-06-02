package output

import (
	"fmt"
	"io"
	"reflect"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

const (
	redColor    = "\u001b[31m"
	greenColor  = "\u001b[32m"
	yellowColor = "\u001b[33m"
	blueColor   = "\u001b[34m"
	// magentaColor = "\u001b[35m"
	cyanColor  = "\u001b[36m"
	resetColor = "\u001b[0m"
)

func red(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", redColor, m, resetColor)
	}
	return m
}

func green(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", greenColor, m, resetColor)
	}
	return m
}

func yellow(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", yellowColor, m, resetColor)
	}
	return m
}
func blue(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", blueColor, m, resetColor)
	}
	return m
}

// func magenta(m string) string {
//     return fmt.Sprintf("%v%v%v", magentaColor, m, resetColor)
// }
func cyan(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", cyanColor, m, resetColor)
	}
	return m
}

func renderInterface(wr io.Writer, i interface{}) error {
	if i == nil {
		return nil
	}
	v := reflect.Indirect(reflect.ValueOf(i))
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		if _, err := wr.Write([]byte(fmt.Sprintf(" * %s: %s\n", fieldName, v.Field(i)))); err != nil {
			return err
		}
	}
	return nil
}

func renderResult(wr io.Writer, result *gornir.JobResult, renderHost bool, color bool) error {
	if renderHost {
		var colorFunc func(string, bool) string
		switch {
		case result.AnyErr() != nil:
			colorFunc = red
		case !result.AnyChanged():
			colorFunc = green
		default:
			colorFunc = yellow
		}
		if _, err := wr.Write([]byte(colorFunc(fmt.Sprintf("@ %s\n", result.Context().Host().Hostname), color))); err != nil {
			return err
		}
	}
	if err := renderInterface(wr, result.Data()); err != nil {
		return err
	}
	if _, err := wr.Write([]byte(fmt.Sprintf("  - err: %v\n\n", result.Err()))); err != nil {
		return err
	}

	for i, r := range result.SubResults() {
		if _, err := wr.Write([]byte(cyan(fmt.Sprintf("**** subtask %d\n", i), color))); err != nil {
			return err
		}
		if err := renderResult(wr, r, false, color); err != nil {
			return err
		}
	}
	return nil
}

// RenderResults writes the contents of the results to an io.Writer in either color or b/w. The
// output will be similar to:
//     # What's my ip?
//     @ dev5.no_group
//       - err: failed to dial: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain
//
//     @ dev1.group_1
//      * Stdout: 10.21.33.101/24
//
//      * Stderr:
//       - err: <nil>
//
//     @ dev2.group_1
//      * Stdout: 10.21.33.102/24
//
//      * Stderr:
//       - err: <nil>
//
//     @ dev3.group_2
//      * Stdout: 10.21.33.103/24
//
//      * Stderr:
//       - err: <nil>
//
//     @ dev4.group_2
//      * Stdout: 10.21.33.104/24
//
//      * Stderr:
//       - err: <nil>
func RenderResults(wr io.Writer, results chan *gornir.JobResult, color bool) error {
	r := <-results

	title := blue(fmt.Sprintf("# %s\n", r.Context().Title()), color)
	if _, err := wr.Write([]byte(title)); err != nil {
		return err
	}
	if err := renderResult(wr, r, true, color); err != nil {
		return err
	}
	for r := range results {
		if err := renderResult(wr, r, true, color); err != nil {
			return err
		}
	}
	return nil
}
