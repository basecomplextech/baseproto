// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import (
	"fmt"

	"github.com/basecomplextech/baseproto/internal/format"
)

// StructType defines a struct type.
type StructType[T any] interface {
	Type[T]
	StructTypeDyn
}

// StructTypeDyn defines a struct type.
type StructTypeDyn interface {
	TypeDyn

	// Fields returns the number of fields in the struct.
	Fields() uint16

	// FieldAt returns the field at the given index.
	FieldAt(i uint16) StructFieldDyn
}

// DecodeStructFunc is a function that decodes a struct.
type DecodeStructFunc[T any] func([]byte) (T, int, error)

// New

// NewStructType returns a new struct type with the given fields.
// The methods panics on duplicate field names or invalid field indexes.
func NewStructType[T any](decode DecodeStructFunc[T], fields ...StructFieldDyn) StructType[T] {
	s, err := newStructType(decode, fields...)
	if err != nil {
		panic(err)
	}
	return s
}

// NewStructTypeErr returns a new struct type with the given fields.
// The methods returns an error on duplicate field names or invalid field indexes.
func NewStructTypeErr[T any](decode DecodeStructFunc[T], fields ...StructFieldDyn) (
	StructType[T], error) {

	return newStructType(decode, fields...)
}

// internal

type structType[T any] struct {
	decode DecodeStructFunc[T]
	fields []StructFieldDyn
	names  map[string]StructFieldDyn
}

func newStructType[T any](decode DecodeStructFunc[T], fields ...StructFieldDyn) (
	*structType[T], error) {

	// Make struct
	s := &structType[T]{
		decode: decode,
		fields: make([]StructFieldDyn, 0, len(fields)),
		names:  make(map[string]StructFieldDyn, len(fields)),
	}

	// Add fields
	for i, f := range fields {
		// Check index
		index := f.Index()
		if index != uint16(i) {
			return nil, fmt.Errorf("invalid field index %d, expected %d", index, i)
		}

		// Check name
		name := f.Name()
		if _, ok := s.names[name]; ok {
			return nil, fmt.Errorf("duplicate field name %q", name)
		}

		// Add field
		s.fields = append(s.fields, f)
		s.names[name] = f
	}

	return s, nil
}

// Kind returns the type kind.
func (s *structType[T]) Kind() format.Kind {
	return format.KindStruct
}

// String returns the string representation of the type.
func (s *structType[T]) String() string {
	return format.KindStruct.String()
}

// Fields

// Fields returns the number of fields in the struct.
func (s *structType[T]) Fields() uint16 {
	return uint16(len(s.fields))
}

// FieldAt returns the field at the given index.
func (s *structType[T]) FieldAt(i uint16) StructFieldDyn {
	return s.fields[i]
}

// Methods

// Open opens a value.
func (s *structType[T]) Open(b []byte) (v T, n int, err error) {
	return s.decode(b)
}

// Parse parses and verifies a value.
func (s *structType[T]) Parse(b []byte) (v T, n int, err error) {
	return s.decode(b)
}

// Verify

// Verify verifies a value against the type.
func (s *structType[T]) Verify(b []byte) error {
	_, _, err := s.decode(b)
	return err
}

// VerifyRaw verifies a raw, possibly untruncated, value against the type.
func (s *structType[T]) VerifyRaw(b []byte) error {
	_, _, err := s.decode(b)
	return err
}

// Internal

// Resolve resolves internal field type references.
func (s *structType[T]) Resolve() {
	for _, field := range s.fields {
		field.Resolve()
	}
}
