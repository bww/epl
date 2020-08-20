package epl

import (
	"strings"
)

func formatArgs(args []executable) string {
	b := &strings.Builder{}
	for i, e := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		e.print(b, 0, printState{})
	}
	return b.String()
}
