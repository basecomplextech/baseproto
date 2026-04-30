// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import (
	"github.com/basecomplextech/baseproto/internal/values"
)

// Value is a raw value.
type Value = values.Value

// OpenValue opens and returns a value from bytes, or nil on error.
//
// The method only checks the type, but does not recursively parse the value.
// See [ParseValue] for recursive parsing.
func OpenValue(b []byte) Value {
	return values.OpenValue(b)
}

// OpenValueErr opens and returns a value from bytes, or an error.
//
// The method only checks the type, but does not recursively parse the value.
// See [ParseValue] for recursive parsing.
func OpenValueErr(b []byte) (Value, error) {
	return values.OpenValueErr(b)
}

// ParseValue recursively parses and returns a value.
func ParseValue(b []byte) (_ Value, n int, err error) {
	return values.ParseValue(b)
}
