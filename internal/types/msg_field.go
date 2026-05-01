// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import "fmt"

// MessageField defines a message field.
type MessageField[T any] interface {
	MessageFieldDyn

	// Type returns the field type.
	Type() Type[T]
}

// MessageFieldDyn defines a message field.
type MessageFieldDyn interface {
	// Tag returns the field tag.
	Tag() uint16

	// Name returns the field name.
	Name() string

	// TypeDyn returns the field type.
	TypeDyn() TypeDyn

	// Methods

	// Verify verifies a value against the field.
	Verify(b []byte) error

	// VerifyRaw verifies a raw, possibly untruncated, value against the field.
	VerifyRaw(b []byte) error
}

// NewMessageField returns a new message field with the given tag, name and type.
func NewMessageField[T any](tag uint16, name string, typ Type[T]) MessageField[T] {
	return newMessageField(tag, name, typ)
}

// internal

var _ MessageField[any] = (*messageField[any])(nil)

type messageField[T any] struct {
	tag  uint16
	name string
	typ  Type[T]
}

func newMessageField[T any](tag uint16, name string, typ Type[T]) *messageField[T] {
	return &messageField[T]{
		tag:  tag,
		name: name,
		typ:  typ,
	}
}

// Tag returns the field tag.
func (f *messageField[T]) Tag() uint16 {
	return f.tag
}

// Name returns the field name.
func (f *messageField[T]) Name() string {
	return f.name
}

// Type returns the field type.
func (f *messageField[T]) Type() Type[T] {
	return f.typ
}

// TypeDyn returns the field type.
func (f *messageField[T]) TypeDyn() TypeDyn {
	return f.typ
}

// Methods

// Verify verifies a value against the field.
func (f *messageField[T]) Verify(value []byte) error {
	if err := f.typ.Verify(value); err != nil {
		return fmt.Errorf("invalid field %q:%d: %w", f.name, f.tag, err)
	}
	return nil
}

// VerifyRaw verifies a raw, possibly untruncated, value against the field.
func (f *messageField[T]) VerifyRaw(b []byte) error {
	if err := f.typ.VerifyRaw(b); err != nil {
		return fmt.Errorf("invalid field %q:%d: %w", f.name, f.tag, err)
	}
	return nil
}
