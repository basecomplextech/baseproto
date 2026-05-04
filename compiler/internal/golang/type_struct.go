// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type StructType interface {
	Type

	// Funcs

	// DecodeFunc returns a decode struct func.
	DecodeFunc() string

	// OpenFunc returns an open struct func.
	OpenFunc() string
}

// internal

var _ StructType = (*structType)(nil)

type structType struct {
	name string
	imp  string
}

func newStructType(typ *model.Type) (*structType, error) {
	kind := typ.Kind
	if kind != model.KindStruct {
		panic("not struct kind")
	}

	t := &structType{
		name: typ.Name,
		imp:  typ.ImportName,
	}
	return t, nil
}

// Kind returns a type kind.
func (t *structType) Kind() model.Kind {
	return model.KindStruct
}

// Name returns a type name.
func (t *structType) Name() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.%v", t.imp, t.name)
	}
	return t.name
}

// Funcs

// DecodeFunc returns a decode struct func.
func (t *structType) DecodeFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Decode%v", t.imp, t.name)
	}
	return "Decode" + t.name
}

// DecodeListElem returns a decode func for a list element.
func (t *structType) DecodeListElem() string {
	return t.DecodeFunc()
}

// OpenFunc returns an open struct func.
func (t *structType) OpenFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Open%v", t.imp, t.name)
	}
	return "Open" + t.name
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *structType) FieldInput() string {
	return t.Name()
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *structType) FieldOutput() string {
	return t.Name()
}

// Write fields

// ReturnField writes a field get.
func (t *structType) ReturnField(w writer.Writer, tag int) error {
	open := t.OpenFunc()
	w.Writef(`return %v(m.msg.FieldRaw(%d))`, open, tag)
	return nil
}

// WriteField writes a field write.
func (t *structType) WriteField(w writer.Writer, tag int) error {
	return nil
}
