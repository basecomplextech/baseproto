// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import (
	"github.com/basecomplextech/baseproto/internal/values"
)

type (
	// List is a raw list of elements.
	List = values.List

	// MessageList is a list of messages.
	MessageList[T any] = values.MessageList[T]

	// ValueList is a list of primitive values.
	ValueList[T any] = values.ValueList[T]
)

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

// MessageList

// NewMessageList returns a new message list.
func NewMessageList[T any](list List, open func([]byte) (T, error)) MessageList[T] {
	return values.NewMessageList(list, open)
}

// OpenMessageList opens and returns a message list, or an empty list on error.
func OpenMessageList[T any](b []byte, open func([]byte) (T, error)) MessageList[T] {
	return values.OpenMessageList(b, open)
}

// OpenMessageListErr opens and returns a message list, or an error.
func OpenMessageListErr[T any](b []byte, open func([]byte) (T, error)) (
	_ MessageList[T], err error) {

	return values.OpenMessageListErr(b, open)
}

// ParseMessageList decodes, recursively validates and returns a list.
func ParseMessageList[T any](b []byte, open func([]byte) (T, error)) (
	_ MessageList[T], size int, err error) {

	return values.ParseMessageList(b, open)
}

// ValueList

// NewValueList returns a new value list.
func NewValueList[T any](list List, decode func([]byte) (T, int, error)) ValueList[T] {
	return values.NewValueList(list, decode)
}

// OpenValueList opens and returns a value list, or an empty list on error.
func OpenValueList[T any](b []byte, decode func([]byte) (T, int, error)) ValueList[T] {
	return values.OpenValueList(b, decode)
}

// OpenValueListErr opens and returns a value list, or an error.
func OpenValueListErr[T any](b []byte, decode func([]byte) (T, int, error)) (
	_ ValueList[T], err error) {

	return values.OpenValueListErr(b, decode)
}

// ParseValueList decodes, recursively validates and returns a list.
func ParseValueList[T any](b []byte, decode func([]byte) (T, int, error)) (
	_ ValueList[T], size int, err error) {

	return values.ParseValueList(b, decode)
}
