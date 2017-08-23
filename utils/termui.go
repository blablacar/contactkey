package utils

import "regexp"

// VTClean returns a string stripped from all its ANSI control codes
func VTClean(s string) string {
	c := regexp.MustCompile("\033\\[(\\d+(;\\d+)?)?\\w")
	return c.ReplaceAllString(s, "")
}
