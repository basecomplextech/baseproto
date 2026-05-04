// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package writer

import (
	"bytes"
	"fmt"
)

type Writer interface {
	// Bytes returns the written bytes.
	Bytes() []byte

	// Line writes args and inserts a newline.
	Line(args ...string)

	// Linef writes a formatted string and inserts a newline.
	Linef(format string, args ...any)

	// Write writes a string.
	Write(args ...string)

	// Writef writes a formatted string.
	Writef(format string, args ...any)
}

// New returns a new writer.
func New() Writer {
	return newWriter()
}

// internal

var _ Writer = (*writer)(nil)

type writer struct {
	b bytes.Buffer
}

func newWriter() *writer {
	return &writer{
		b: bytes.Buffer{},
	}
}

func (w *writer) Bytes() []byte {
	return w.b.Bytes()
}

// Line writes args and inserts a newline.
func (w *writer) Line(args ...string) {
	w.Write(args...)
	w.b.WriteString("\n")
}

// Linef writes a formatted string and inserts a newline.
func (w *writer) Linef(format string, args ...any) {
	w.Writef(format, args...)
	w.b.WriteString("\n")
}

// Write writes a string.
func (w *writer) Write(args ...string) {
	for _, s := range args {
		w.b.WriteString(s)
	}
}

// Writef writes a formatted string.
func (w *writer) Writef(format string, args ...any) {
	if len(args) == 0 {
		w.Write(format)
		return
	}

	s := fmt.Sprintf(format, args...)
	w.b.WriteString(s)
}
