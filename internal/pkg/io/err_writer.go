package io

import (
	"fmt"
	"strings"
)

// ErrWriter wraps a strings.Builder and any errors that have occurred whilst
// writing multiple values; meaning errors only have to be handled once
// at the end of multiple writes.
type ErrWriter struct {
	sb  strings.Builder
	err error
}

// NewErrWriter returns a pointer to a new instance of ErrWriter.
func NewErrWriter(sb strings.Builder) *ErrWriter {
	return &ErrWriter{
		sb: sb,
	}
}

// WriteString writes a string to the underlying strings.Builder.
func (ew *ErrWriter) WriteString(format string, args ...any) {
	if ew.err != nil {
		return
	}

	_, ew.err = ew.sb.Write([]byte(fmt.Sprintf(format, args...)))
}

// String returns the underlying error.
func (ew *ErrWriter) String() string {
	return ew.sb.String()
}

// Err returns the underlying error.
func (ew *ErrWriter) Err() error {
	return ew.err
}
