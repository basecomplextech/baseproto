// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package writer

// ValueListWriter writes a list of primitive values.
type ValueListWriter[T any] struct {
	w     ListWriter
	write WriteFunc[T]
}

// NewValueListWriter returns a new value list writer.
func NewValueListWriter[T any](w ListWriter, write WriteFunc[T]) ValueListWriter[T] {
	return ValueListWriter[T]{
		w:     w,
		write: write,
	}
}

// Add adds the next element.
func (b ValueListWriter[T]) Add(value T) error {
	return WriteElement(b.w, value, b.write)
}

// Len returns the number of written elements.
// The method is only valid when there is no pending element.
func (b ValueListWriter[T]) Len() int {
	return b.w.Len()
}

// End ends the list.
func (b ValueListWriter[T]) End() error {
	return b.w.End()
}
