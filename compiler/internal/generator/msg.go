// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"github.com/basecomplextech/baseproto/compiler/internal/golang"
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type messageWriter struct {
	writer.Writer
}

func newMessageWriter(w writer.Writer) *messageWriter {
	return &messageWriter{w}
}

func (w *messageWriter) write(def *model.Definition) error {
	msg, err := golang.NewMessage(def)
	if err != nil {
		return err
	}

	if err := w.def(msg); err != nil {
		return err
	}
	if err := w.new_methods(msg); err != nil {
		return err
	}
	if err := w.parse_method(msg); err != nil {
		return err
	}
	if err := w.fields(msg); err != nil {
		return err
	}
	if err := w.has_fields(msg); err != nil {
		return err
	}
	if err := w.methods(msg); err != nil {
		return err
	}
	return nil
}

func (w *messageWriter) def(msg *golang.Message) error {
	w.Linef(`// %v`, msg.Name)
	w.Line()
	w.Linef(`type %v struct {`, msg.Name)
	w.Line(`msg baseproto.Message`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) new_methods(msg *golang.Message) error {
	w.Linef(`func New%v(msg baseproto.Message) %v {`, msg.Name, msg.Name)
	w.Linef(`return %v{msg}`, msg.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func Open%v(b []byte) %v {`, msg.Name, msg.Name)
	w.Linef(`msg := baseproto.OpenMessage(b)`)
	w.Linef(`return %v{msg}`, msg.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func Open%vErr(b []byte) (_ %v, err error) {`, msg.Name, msg.Name)
	w.Linef(`msg, err := baseproto.OpenMessageErr(b)`)
	w.Linef(`return %v{msg}, err`, msg.Name)
	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) parse_method(msg *golang.Message) error {
	w.Linef(`func Parse%v(b []byte) (_ %v, size int, err error) {`, msg.Name, msg.Name)
	w.Linef(`msg, size, err := baseproto.ParseMessage(b)`)
	w.Linef(`return %v{msg}, size, err`, msg.Name)
	w.Linef(`}`)
	w.Line()
	return nil
}

// fields

func (w *messageWriter) fields(msg *golang.Message) error {
	for _, field := range msg.Fields {
		if err := w.field(msg, field); err != nil {
			return err
		}
	}

	if len(msg.Fields) > 1 {
		w.Line()
	}
	return nil
}

func (w *messageWriter) field(msg *golang.Message, field *golang.MessageField) error {
	typ := field.Type

	w.Writef(`func (m %v) %v() %v {`, msg.Name, field.Name, typ.FieldOutput())
	typ.ReturnField(w, field.Tag)
	w.Writef(`}`)
	w.Line()
	return nil
}

// has fields

func (w *messageWriter) has_fields(msg *golang.Message) error {
	for _, field := range msg.Fields {
		if err := w.has_field(msg, field); err != nil {
			return err
		}
	}

	if len(msg.Fields) > 1 {
		w.Line()
	}
	return nil
}

func (w *messageWriter) has_field(msg *golang.Message, field *golang.MessageField) error {
	w.Writef(`func (m %v) Has%v() bool {`, msg.Name, field.Name)
	w.Writef(`return m.msg.HasField(%d)`, field.Tag)
	w.Writef(`}`)
	w.Line()
	return nil
}

// methods

func (w *messageWriter) methods(msg *golang.Message) error {
	w.Writef(`func (m %v) Clone() %v {`, msg.Name, msg.Name)
	w.Writef(`return %v{m.msg.Clone()}`, msg.Name)
	w.Writef(`}`)
	w.Line()

	w.Writef(`func (m %v) CloneToArena(a alloc.Arena) %v {`, msg.Name, msg.Name)
	w.Writef(`return %v{m.msg.CloneToArena(a)}`, msg.Name)
	w.Writef(`}`)
	w.Line()

	w.Writef(`func (m %v) CloneToBuffer(b buffer.Buffer) %v {`, msg.Name, msg.Name)
	w.Writef(`return %v{m.msg.CloneToBuffer(b)}`, msg.Name)
	w.Writef(`}`)
	w.Line()

	w.Line()

	w.Writef(`func (m %v) IsEmpty() bool {`, msg.Name)
	w.Writef(`return m.msg.Empty()`)
	w.Writef(`}`)
	w.Line()

	w.Writef(`func (m %v) Unwrap() baseproto.Message {`, msg.Name)
	w.Writef(`return m.msg`)
	w.Writef(`}`)
	w.Line()
	return nil
}

// writer

// util

func messageFieldName(field *model.Field) string {
	return toUpperCamelCase(field.Name)
}
