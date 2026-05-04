// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type AnyType interface {
	Type

	any()
}

// internal

var _ AnyType = (*anyType)(nil)

type anyType struct{}

func newAnyType(typ *model.Type) (*anyType, error) {
	kind := typ.Kind
	if kind != model.KindAny {
		panic("not any kind")
	}

	t := &anyType{}
	return t, nil
}

func (t *anyType) any() {}

// Kind returns a type kind.
func (t *anyType) Kind() model.Kind {
	return model.KindAny
}

// Name returns a type name.
func (t *anyType) Name() string {
	return "baseproto.Value"
}

// Funcs

// DecodeListElem returns a decode func for a list element.
func (t *anyType) DecodeListElem() string {
	return "baseproto.OpenValue"
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *anyType) FieldInput() string {
	return "baseproto.Value"
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *anyType) FieldOutput() string {
	return "baseproto.Value"
}

// Write fields

// ReturnField writes a field get.
func (t *anyType) ReturnField(w writer.Writer, tag int) error {
	w.Writef(`return m.msg.Field(%d)`, tag)
	return nil
}

// WriteField writes a field write.
func (t *anyType) WriteField(w writer.Writer, tag int) error {
	return nil
}
