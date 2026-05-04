// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type EnumType interface {
	Type
	enum()

	// Funcs

	// DecodeFunc returns a decode func for an enum.
	DecodeFunc() string

	// OpenFunc returns an open func for an enum.
	OpenFunc() string
}

// internal

var _ EnumType = (*enumType)(nil)

type enumType struct {
	name string
	imp  string
}

func newEnumType(typ *model.Type) (*enumType, error) {
	kind := typ.Kind
	if kind != model.KindEnum {
		panic("not enum kind")
	}

	t := &enumType{
		name: typ.Name,
		imp:  typ.ImportName,
	}
	return t, nil
}

func (t *enumType) enum() {}

// Kind returns a type kind.
func (t *enumType) Kind() model.Kind {
	return model.KindEnum
}

// Name returns a type name.
func (t *enumType) Name() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.%v", t.imp, t.name)
	}
	return t.name
}

// Funcs

// DecodeFunc returns a decode func for an enum.
func (t *enumType) DecodeFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Decode%v", t.imp, t.name)
	}
	return fmt.Sprintf("Decode%v", t.name)
}

// OpenFunc returns an open func for an enum.
func (t *enumType) OpenFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Open%v", t.imp, t.name)
	}
	return fmt.Sprintf("Open%v", t.name)
}

// DecodeListElem returns a decode func for a list element.
func (t *enumType) DecodeListElem() string {
	return t.DecodeFunc()
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *enumType) FieldInput() string {
	return t.Name()
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *enumType) FieldOutput() string {
	return t.Name()
}

// Write fields

// ReturnField writes a field get.
func (t *enumType) ReturnField(w writer.Writer, tag int) error {
	open := t.OpenFunc()

	w.Writef(`return %v(m.msg.FieldRaw(%d))`, open, tag)
	return nil
}

// WriteField writes a field write.
func (t *enumType) WriteField(w writer.Writer, tag int) error {
	return nil
}
