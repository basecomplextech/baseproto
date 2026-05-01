// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import (
	"github.com/basecomplextech/baseproto/internal/format"
)

// Type defines a type.
type Type[T any] interface {
	TypeDyn

	// Open opens a value.
	Open(b []byte) (v T, n int, err error)

	// Parse parses and verifies a value.
	Parse(b []byte) (v T, n int, err error)
}

// TypeDyn defines a type.
type TypeDyn interface {
	// Kind returns the type kind.
	Kind() format.Kind

	// String returns the string representation of the type.
	String() string

	// Verify

	// Verify verifies a value against the type.
	Verify(b []byte) error

	// VerifyRaw verifies a raw, possibly untruncated, value against the type.
	VerifyRaw(b []byte) error
}
