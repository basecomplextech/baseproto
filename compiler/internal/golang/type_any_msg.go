// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type AnyMessageType interface {
	Type

	anyMessage()
}

// internal

var _ AnyMessageType = (*anyMessageType)(nil)

type anyMessageType struct{}

func newAnyMessageType(typ *model.Type) (*anyMessageType, error) {
	kind := typ.Kind
	if kind != model.KindAnyMessage {
		panic("not any message kind")
	}

	t := &anyMessageType{}
	return t, nil
}

func (t *anyMessageType) anyMessage() {}

// Kind returns a type kind.
func (t *anyMessageType) Kind() model.Kind {
	return model.KindAnyMessage
}

// Name returns a type name.
func (t *anyMessageType) Name() string {
	return "baseproto.Message"
}

// Funcs

// DecodeListElem returns a decode func for a list element.
func (t *anyMessageType) DecodeListElem() string {
	return "baseproto.ParseMessage"
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *anyMessageType) FieldInput() string {
	return "baseproto.Message"
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *anyMessageType) FieldOutput() string {
	return "baseproto.Message"
}

// Write fields

// ReturnField writes a field get.
func (t *anyMessageType) ReturnField(w writer.Writer, tag int) error {
	w.Writef(`return m.msg.Field(%d).Message()`, tag)
	return nil
}

// WriteField writes a field write.
func (t *anyMessageType) WriteField(w writer.Writer, tag int) error {
	return nil
}
