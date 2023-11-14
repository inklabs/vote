package vote_test

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// TrimmingWriter trims trailing whitespace for each line before writing to the underlying writer.
// It also ignores multiple newlines due to this [issue](https://github.com/golang/go/issues/59191)
type TrimmingWriter struct {
	w               io.Writer
	isLastLineEmpty bool
}

// NewTrimmingWriter creates a new TrimmingWriter.
func NewTrimmingWriter(w io.Writer) *TrimmingWriter {
	return &TrimmingWriter{w: w}
}

// Write trims trailing whitespace for each line before writing to the underlying writer.
func (tw *TrimmingWriter) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimRight(line, " \t\r\n")
		if trimmedLine == "" {
			if tw.isLastLineEmpty {
				continue
			}
			tw.isLastLineEmpty = true
		} else {
			tw.isLastLineEmpty = false
		}
		_, err := tw.w.Write([]byte(trimmedLine + "\n"))
		if err != nil {
			return n, err
		}
		n += len(trimmedLine) + 1 // Include the newline character
	}
	return n, nil
}