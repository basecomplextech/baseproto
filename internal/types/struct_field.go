// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

// StructField defines a struct field.
type StructField[T any] interface {
	StructFieldDyn

	// Type returns the field type.
	Type() Type[T]
}

// StructFieldDyn defines a struct field.
type StructFieldDyn interface {
	// Index returns the field index in the struct.
	Index() uint16

	// Name returns the field name.
	Name() string

	// TypeDyn returns the field type.
	TypeDyn() TypeDyn
}

// NewStructField returns a new struct field.
func NewStructField[T any](index uint16, name string, typ Type[T]) StructField[T] {
	return newStructField(index, name, typ)
}

// internal

var _ StructField[any] = (*structField[any])(nil)

type structField[T any] struct {
	index uint16
	name  string
	typ   Type[T]
}

func newStructField[T any](index uint16, name string, typ Type[T]) *structField[T] {
	return &structField[T]{
		index: index,
		name:  name,
		typ:   typ,
	}
}

// Index returns the field index in the struct.
func (f *structField[T]) Index() uint16 {
	return f.index
}

// Name returns the field name.
func (f *structField[T]) Name() string {
	return f.name
}

// Type returns the field type.
func (f *structField[T]) Type() Type[T] {
	return f.typ
}

// TypeDyn returns the field type.
func (f *structField[T]) TypeDyn() TypeDyn {
	return f.typ
}
