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
	EncodableType
	enum()

	// Funcs

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

// InputName returns an input type, i.e. string (not baseproto.String).
func (t *enumType) InputName() string {
	return t.Name()
}

// OutputName returns an output type, i.e. baseproto.String (not string).
func (t *enumType) OutputName() string {
	return t.Name()
}

// Funcs

// OpenFunc returns an open func for an enum.
func (t *enumType) OpenFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Open%v", t.imp, t.name)
	}
	return fmt.Sprintf("Open%v", t.name)
}

// EncodeFunc returns an encode func for an enum.
func (t *enumType) EncodeFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Encode%v", t.imp, t.name)
	}
	return fmt.Sprintf("Encode%v", t.name)
}

// DecodeFunc returns a decode func for an enum.
func (t *enumType) DecodeFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Decode%v", t.imp, t.name)
	}
	return fmt.Sprintf("Decode%v", t.name)
}

// DecodeCloneFunc returns a decode func, which returns string clones.
func (t *enumType) DecodeCloneFunc() string {
	return t.DecodeFunc()
}

// ParseFunc returns a parse func.
func (t *enumType) ParseFunc() string {
	return t.DecodeFunc()
}

// List

// AddListElem returns an encode func for a list element.
func (t *enumType) AddListElem() string {
	return t.EncodeFunc()
}

// GetListElem returns a decode func for a list element.
func (t *enumType) GetListElem() string {
	return t.DecodeFunc()
}

// Message

// GetField writes a field get.
func (t *enumType) GetField(w writer.Writer, tag int) error {
	open := t.OpenFunc()

	w.Writef(`return %v(m.msg.FieldRaw(%d))`, open, tag)
	return nil
}
