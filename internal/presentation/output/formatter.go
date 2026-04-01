package output

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Password(s string) string {
	return fmt.Sprintf("%s", s)
}

func (f *Formatter) Success(msg string) {
	fmt.Printf("✓ %s\n", msg)
}

func (f *Formatter) Error(msg string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
}

func (f *Formatter) Config(cfg interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Current configuration:")

	fmt.Fprintln(w, "  (see ep config show commands)")
	w.Flush()
}
