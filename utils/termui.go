package utils

import (
	"fmt"
	"io"
	"regexp"
)

// VTClean returns a string stripped from all its ANSI control codes
func VTClean(s string) string {
	c := regexp.MustCompile("\033\\[(\\d+(;\\d+)?)?\\w")
	return c.ReplaceAllString(s, "")
}

func RenderProgres(w io.Writer, state string, progress int) {
	switch progress {
	case 100:
		fmt.Fprintf(w, "\u2713 %s\n", state)
	case 1:
		fmt.Fprintf(w, "\u2026 %s\n", state)

	default:
		fmt.Fprintf(w, "%d%% %s\n", progress, state)
	}

}
