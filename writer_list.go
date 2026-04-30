// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import (
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/writer"
)

type (
	// ListWriter is a raw list writer.
	ListWriter = writer.ListWriter

	// MessageListWriter writes a list of messages.
	MessageListWriter[T any] = writer.MessageListWriter[T]

	// ValueListWriter writes a list of primitive values.
	ValueListWriter[T any] = writer.ValueListWriter[T]
)

// NewListWriter returns a new list writer with a new empty buffer.
//
// The writer is released on end.
func NewListWriter() ListWriter {
	w := writer.New(true /* release */)
	return w.List()
}

// NewListWriterBuffer returns a new list writer with the given buffer.
//
// The writer is freed on end.
func NewListWriterBuffer(buf buffer.Buffer) ListWriter {
	w := writer.Acquire(buf)
	return w.List()
}

// Message

// NewMessageListWriter returns a new message list writer.
func NewMessageListWriter[T any](w ListWriter, next func(w MessageWriter) T) MessageListWriter[T] {
	return writer.NewMessageListWriter(w, next)
}

// Value

// NewValueListWriter returns a new value list writer.
func NewValueListWriter[T any](w ListWriter, write writer.WriteFunc[T]) (_ ValueListWriter[T]) {
	return writer.NewValueListWriter(w, write)
}
