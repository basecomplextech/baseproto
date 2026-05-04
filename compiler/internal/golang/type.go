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

	// InputName returns an input type, i.e. string (not baseproto.String).
	InputName() string

	// OutputName returns an output type, i.e. baseproto.String (not string).
	OutputName() string

	// Funcs

	// ParseFunc returns a parse func.
	ParseFunc() string

	// List

	// AddListElem returns an encode func for a list element.
	AddListElem() string

	// GetListElem returns a decode func for a list element.
	GetListElem() string

	// Message

	// GetField writes a field get.
	GetField(w writer.Writer, tag int) error
}

// EncodableType is a common interface for values/enums/structs.
type EncodableType interface {
	Type

	// EncodeFunc returns an encode func.
	EncodeFunc() string

	// DecodeFunc returns a decode func.
	DecodeFunc() string

	// DecodeCloneFunc returns a decode func, which returns string clones.
	DecodeCloneFunc() string
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
