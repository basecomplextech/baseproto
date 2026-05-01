// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import "github.com/basecomplextech/baseproto/internal/types"

type (
	// Type defines a type.
	Type[T any] = types.Type[T]

	// TypeDyn defines a non-generic type.
	TypeDyn = types.TypeDyn

	// Enum

	// EnumType defines an enum type.
	EnumType[T int32] = types.EnumType[T]

	// EnumTypeDyn defines a non-generic enum type.
	EnumTypeDyn = types.EnumTypeDyn

	// List

	// StructType defines a struct type.
	ListType[T any] = types.ListType[T]

	// StructTypeDyn defines a non-generic struct type.
	ListTypeDyn = types.ListTypeDyn

	// Message

	// MessageField defines a message field.
	MessageField[T any] = types.MessageField[T]

	// MessageFieldDyn defines a non-generic message field.
	MessageFieldDyn = types.MessageFieldDyn

	// MessageType defines a message type.
	MessageType2[T MessageType] = types.MessageType[T]

	// MessageTypeDyn defines a non-generic message type.
	MessageTypeDyn = types.MessageTypeDyn

	// Struct

	// StructField defines a struct field.
	StructField[T any] = types.StructField[T]

	// StructFieldDyn defines a non-generic struct field.
	StructFieldDyn = types.StructFieldDyn

	// StructType defines a struct type.
	StructType[T any] = types.StructType[T]

	// StructTypeDyn defines a non-generic struct type.
	StructTypeDyn = types.StructTypeDyn
)

type (
	// NewListFunc is a function that returns a new list.
	NewListFunc[T any] = types.NewListFunc[T]

	// NewMessageFunc is a function that returns a new message.
	NewMessageFunc[T MessageType] = types.NewMessageFunc[T]

	// DecodeEnumFunc is a function that decodes an enum value.
	DecodeEnumFunc[T int32] = types.DecodeEnumFunc[T]

	// DecodeStructFunc is a function that decodes a struct.
	DecodeStructFunc[T any] = types.DecodeStructFunc[T]
)

var (
	TypeBool = types.Bool
	TypeByte = types.Byte

	TypeInt16 = types.Int16
	TypeInt32 = types.Int32
	TypeInt64 = types.Int64

	TypeUint16 = types.Uint16
	TypeUint32 = types.Uint32
	TypeUint64 = types.Uint64

	TypeFloat32 = types.Float32
	TypeFloat64 = types.Float64

	TypeBin64  = types.Bin64
	TypeBin128 = types.Bin128
	TypeBin256 = types.Bin256

	TypeBytes  = types.Bytes
	TypeString = types.String
)

// Enum

// NewEnumType returns a new enum type.
func NewEnumType[T int32](decode DecodeEnumFunc[T]) EnumType[T] {
	return types.NewEnumType(decode)
}

// List

// NewListType returns a new list type.
func NewListType[T any](new NewListFunc[T], elem Type[T]) ListType[T] {
	return types.NewListType(new, elem)
}

// Message

// NewMessageField returns a new message field.
func NewMessageField[T any](tag uint16, name string, typ Type[T]) MessageField[T] {
	return types.NewMessageField(tag, name, typ)
}

// NewMessageType returns a new message type with the given fields.
func NewMessageType[T MessageType](new NewMessageFunc[T],
	fields ...MessageFieldDyn) MessageType2[T] {

	return types.NewMessageType(new, fields...)
}

// Struct

// NewStructField returns a new struct field.
func NewStructField[T any](index uint16, name string, typ Type[T]) StructField[T] {
	return types.NewStructField(index, name, typ)
}

// NewStructType returns a new struct type with the given fields.
func NewStructType[T any](decode DecodeStructFunc[T], fields ...StructFieldDyn) StructType[T] {
	return types.NewStructType(decode, fields...)
}
