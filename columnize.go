package columnize

import (
	"fmt"
	"strings"
)

type Config struct {
	// The string by which the lines of input will be split.
	Delim string

	// The string by which columns of output will be separated.
	Glue string
}

// Returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Delim: "|",
		Glue:  "  ",
	}
}

// Returns a list of elements, each representing a single item which will
// belong to a column of output.
func getElementsFromLine(line string, delim string) []interface{} {
	elements := make([]interface{}, 0)
	for _, field := range strings.Split(line, delim) {
		elements = append(elements, strings.TrimSpace(field))
	}
	return elements
}

// Examines a list of strings and determines how wide each column should be
// considering all of the elements that need to be printed within it.
func getWidthsFromLines(lines []string, delim string) []int {
	var widths []int

	for _, line := range lines {
		elems := getElementsFromLine(line, delim)
		for i := 0; i < len(elems); i++ {
			if len(widths) <= i {
				widths = append(widths, len(elems[i].(string)))
			} else if widths[i] < len(elems[i].(string)) {
				widths[i] = len(elems[i].(string))
			}
		}
	}
	return widths
}

// Given a set of column widths and the number of columns in the current line,
// returns a sprintf-style format string which can be used to print output
// aligned properly with other lines using the same widths set.
func getStringFormat(widths []int, columns int, space string) string {
	var stringfmt string

	// Create the format string from the discovered widths
	for i := 0; i < columns && i < len(widths); i++ {
		if i == columns-1 {
			stringfmt += "%s\n"
		} else {
			stringfmt += fmt.Sprintf("%%-%ds%s", widths[i], space)
		}
	}
	return stringfmt
}

// Format is the public-facing interface that takes either a plain string
// or a list of strings, plus a delimiter, and returns nicely aligned output.
func Format(input interface{}, config *Config) (string, error) {
	var result string
	var lines []string

	switch in := input.(type) {
	case string:
		for _, line := range strings.Split(in, "\n") {
			lines = append(lines, line)
		}

	case []string:
		lines = in

	default:
		return "", fmt.Errorf("columnize: Expected string or []string")
	}

	widths := getWidthsFromLines(lines, config.Delim)

	// Create the formatted output using the format string
	for _, line := range lines {
		elems := getElementsFromLine(line, config.Delim)
		stringfmt := getStringFormat(widths, len(elems), config.Glue)
		result += fmt.Sprintf(stringfmt, elems...)
	}
	return strings.TrimSpace(result), nil
}

// Convenience function for using Columnize as easy as possible.
func SimpleFormat(input interface{}) (string, error) {
	config := DefaultConfig()
	return Format(input, config)
}
