package output

import (
	"bytes"
	"fmt"
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

func red(m string) string {
	return fmt.Sprintf("%v%v%v", redColor, m, resetColor)
}

func green(m string) string {
	return fmt.Sprintf("%v%v%v", greenColor, m, resetColor)
}

func yellow(m string) string {
	return fmt.Sprintf("%v%v%v", yellowColor, m, resetColor)
}
func blue(m string) string {
	return fmt.Sprintf("%v%v%v", blueColor, m, resetColor)
}

// func magenta(m string) string {
//     return fmt.Sprintf("%v%v%v", magentaColor, m, resetColor)
// }
func cyan(m string) string {
	return fmt.Sprintf("%v%v%v", cyanColor, m, resetColor)
}

func renderInterface(buf *bytes.Buffer, i interface{}) {
	if i == nil {
		return
	}
	v := reflect.Indirect(reflect.ValueOf(i))
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		buf.Write([]byte(fmt.Sprintf(" * %s: %s\n", fieldName, v.Field(i))))
	}
}

func renderResult(buf *bytes.Buffer, result *gornir.JobResult, renderHost bool) {
	if renderHost {
		var colorFunc func(string) string
		switch {
		case result.AnyErr() != nil:
			colorFunc = red
		case !result.AnyChanged():
			colorFunc = green
		default:
			colorFunc = yellow
		}
		buf.Write([]byte(colorFunc(fmt.Sprintf("@ %s\n", result.Context().Host.Hostname))))
	}
	renderInterface(buf, result.Data())
	buf.Write([]byte(fmt.Sprintf("  - err: %v\n\n", result.Err())))

	for i, r := range result.SubResults() {
		buf.Write([]byte(cyan(fmt.Sprintf("**** subtask %d\n", i))))
		renderResult(buf, r, false)
	}
}

func RenderResults(results chan *gornir.JobResult) string {
	var buf bytes.Buffer
	r := <-results

	title := blue(fmt.Sprintf("# %s\n", r.Context().Title()))
	buf.Write([]byte(title))
	renderResult(&buf, r, true)
	for r := range results {
		renderResult(&buf, r, true)
	}
	return buf.String()
}
