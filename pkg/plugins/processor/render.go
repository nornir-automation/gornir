package processor

import (
	"context"
	"fmt"
	"io"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

const (
	redColor   = "\u001b[31m"
	greenColor = "\u001b[32m"
	// yellowColor = "\u001b[33m"
	blueColor = "\u001b[34m"
	// magentaColor = "\u001b[35m"
	// cyanColor  = "\u001b[36m"
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

// func yellow(m string, color bool) string {
//     if color {
//         return fmt.Sprintf("%v%v%v", yellowColor, m, resetColor)
//     }
//     return m
// }

func blue(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", blueColor, m, resetColor)
	}
	return m
}

// func magenta(m string) string {
//     return fmt.Sprintf("%v%v%v", magentaColor, m, resetColor)
// }
// func cyan(m string, color bool) string {
//     if color {
//         return fmt.Sprintf("%v%v%v", cyanColor, m, resetColor)
//     }
//     return m
// }

type render struct {
	wr    io.Writer
	color bool
}

// Render returns a processor that renders the output
func Render(wr io.Writer, color bool) *render {
	return &render{
		wr:    wr,
		color: color,
	}
}

func (r *render) TaskStart(ctx context.Context, logger gornir.Logger, task gornir.Task) error {
	fmt.Println(blue("Starting Task", r.color))
	return nil
}

func (r *render) TaskCompleted(ctx context.Context, logger gornir.Logger, task gornir.Task) error {
	return nil
}

func (r *render) HostStart(ctx context.Context, logger gornir.Logger, host *gornir.Host, task gornir.Task) error {
	return nil
}

func (r *render) HostCompleted(ctx context.Context, logger gornir.Logger, jobResult *gornir.JobResult, host *gornir.Host, task gornir.Task) error {
	switch jobResult.Err() {
	case nil:
		if _, err := r.wr.Write([]byte(green(fmt.Sprintf("@ %s\n", host.Hostname), r.color))); err != nil {
			return err
		}
		if _, err := r.wr.Write([]byte(fmt.Sprintf("%v\n\n", jobResult.Data()))); err != nil {
			return err
		}
	default:
		if _, err := r.wr.Write([]byte(red(fmt.Sprintf("@ %s\n", host.Hostname), r.color))); err != nil {
			return err
		}
		if _, err := r.wr.Write([]byte(fmt.Sprintf("  - err: %v\n\n", jobResult.Err()))); err != nil {
			return err
		}
	}
	return nil
}
