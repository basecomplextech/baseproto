// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type MessageType interface {
	Type

	// Funcs

	// NewFunc returns a new message function.
	NewFunc() string
}

// internal

var _ MessageType = (*messageType)(nil)

type messageType struct {
	name string
	imp  string
}

func newMessageType(typ *model.Type) (*messageType, error) {
	kind := typ.Kind
	if kind != model.KindMessage {
		panic("not message kind")
	}

	t := &messageType{
		name: typ.Name,
		imp:  typ.ImportName,
	}
	return t, nil
}

// Kind returns a type kind.
func (t *messageType) Kind() model.Kind {
	return model.KindMessage
}

// Name returns a type name.
func (t *messageType) Name() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.%v", t.imp, t.name)
	}
	return t.name
}

// Funcs

// NewFunc returns a new message function.
func (t *messageType) NewFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.New%v", t.imp, t.name)
	}
	return fmt.Sprintf("New%v", t.name)
}

// DecodeListElem returns a decode func for a list element.
func (t *messageType) DecodeListElem() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Open%vErr", t.imp, t.name)
	}
	return fmt.Sprintf("Open%vErr", t.name)
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *messageType) FieldInput() string {
	return t.Name()
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *messageType) FieldOutput() string {
	return t.Name()
}

// Write fields

// ReturnField writes a field get.
func (t *messageType) ReturnField(w writer.Writer, tag int) error {
	new := t.NewFunc()
	w.Writef(`return %v(m.msg.Message(%d))`, new, tag)
	return nil
}

// WriteField writes a field write.
func (t *messageType) WriteField(w writer.Writer, tag int) error {
	return nil
}
