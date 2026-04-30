// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import (
	"github.com/basecomplextech/baseproto/internal/values"
)

// List is a raw list of elements.
type List = values.List

// OpenList opens and returns a list from bytes, or an empty list on error.
// The method decodes the list table, but not the elements, see [ParseList].
func OpenList(b []byte) List {
	return values.OpenList(b)
}

// OpenListErr opens and returns a list from bytes, or an error.
// The method decodes the list table, but not the elements, see [ParseList].
func OpenListErr(b []byte) (List, error) {
	return values.OpenListErr(b)
}

// ParseList recursively parses and returns a list.
func ParseList(b []byte) (l List, size int, err error) {
	return values.ParseList(b)
}
