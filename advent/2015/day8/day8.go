package day8

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range strings.SplitSeq(string(bytes.TrimSpace(input)), "\n") {
		unquoted, quoted, err := unquotedQuoted(line)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to unquote string %s %v", line, err)
		}
		// count size of line - size in characters
		s1 += len(line) - unquoted
		// count the quoted characters - the size of the line
		s2 += quoted - len(line)
	}
	return s1, s2, nil
}

// unquotedQuoted counts unquoted and quoted lengths. thanks golang stdlib
func unquotedQuoted(line string) (int, int, error) {
	unquoted, err := strconv.Unquote(line)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to unquote string %s %v", line, err)
	}
	quoted := strconv.Quote(line)
	slog.Debug(line, "unquoted", unquoted, "quoted", quoted)

	return len(unquoted), len(quoted), nil

}
