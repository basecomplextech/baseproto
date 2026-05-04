// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

// Type references a value type or a user-defined message/enum/struct/service.
type Type interface {
	// Kind returns a type kind.
	Kind() model.Kind

	// Name returns a type name.
	Name() string

	// Funcs

	// DecodeListElem returns a decode func for a list element.
	DecodeListElem() string

	// Fields

	// FieldInput returns an input field type, i.e. string (not baseproto.String).
	FieldInput() string

	// FieldOutput returns an output field type, i.e. baseproto.String (not string).
	FieldOutput() string

	// Write fields

	// ReturnField writes a field get.
	ReturnField(w writer.Writer, tag int) error

	// WriteField writes a field write.
	WriteField(w writer.Writer, tag int) error
}

// NewType returns a new type.
func NewType(typ *model.Type) (Type, error) {
	return newType(typ)
}

// internal

func newType(typ *model.Type) (Type, error) {
	kind := typ.Kind
	switch kind {
	case model.KindAny:
		return newAnyType(typ)
	case model.KindAnyMessage:
		return newAnyMessageType(typ)
	case model.KindEnum:
		return newEnumType(typ)
	case model.KindList:
		return newListType(typ)
	case model.KindMessage:
		return newMessageType(typ)
	case model.KindStruct:
		return newStructType(typ)
	default:
		return newValueType(typ)
	}
}
