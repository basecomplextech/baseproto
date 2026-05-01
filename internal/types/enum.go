// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import "github.com/basecomplextech/baseproto/internal/format"

// EnumType defines an enum type.
type EnumType[T int32] interface {
	Type[T]
	EnumTypeDyn
}

// EnumTypeDyn defines a non-generic enum type.
type EnumTypeDyn interface {
	TypeDyn
}

// DecodeEnumFunc is a function that decodes an enum value from bytes.
type DecodeEnumFunc[T int32] func(b []byte) (T, int, error)

// New

// NewEnumType returns a new enum type.
func NewEnumType[T int32](decode DecodeEnumFunc[T]) EnumType[T] {
	return newEnumType(decode)
}

// internal

var _ EnumType[int32] = (*enumType[int32])(nil)

type enumType[T int32] struct {
	decode DecodeEnumFunc[T]
}

func newEnumType[T int32](decode DecodeEnumFunc[T]) *enumType[T] {
	return &enumType[T]{decode: decode}
}

// Kind returns the type kind.
func (t *enumType[T]) Kind() format.Kind {
	return format.KindInt32
}

// String returns the string representation of the type.
func (t *enumType[T]) String() string {
	return "enum"
}

// Methods

// Open opens a value.
func (t *enumType[T]) Open(b []byte) (v T, n int, err error) {
	return t.decode(b)
}

// Parse parses and verifies a value.
func (t *enumType[T]) Parse(b []byte) (v T, n int, err error) {
	return t.decode(b)
}

// Verify

// Verify verifies a value against the type.
func (t *enumType[T]) Verify(b []byte) error {
	_, _, err := t.decode(b)
	return err
}

// VerifyRaw verifies a raw, possibly untruncated, value against the type.
func (t *enumType[T]) VerifyRaw(b []byte) error {
	_, _, err := t.decode(b)
	return err
}

// Internal

// Resolve resolves internal type references.
func (t *enumType[T]) Resolve() {}
