// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"github.com/basecomplextech/baseproto/compiler/internal/golang"
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type messageWriterWriter struct {
	writer.Writer
}

func newMessageWriterWriter(w writer.Writer) *messageWriterWriter {
	return &messageWriterWriter{w}
}

func (w *messageWriterWriter) write(def *model.Definition) error {
	msg, err := golang.NewMessage(def)
	if err != nil {
		return err
	}

	if err := w.def(msg); err != nil {
		return err
	}
	if err := w.new(msg); err != nil {
		return err
	}
	if err := w.fields(msg); err != nil {
		return err
	}
	if err := w.copy(msg); err != nil {
		return err
	}
	if err := w.end(msg); err != nil {
		return err
	}
	return nil
}

func (w *messageWriterWriter) def(msg *golang.Message) error {
	w.Linef(`// %v`, msg.Writer)
	w.Line()
	w.Writef(`type %v struct {`, msg.Writer)
	w.Write(`w baseproto.MessageWriter`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *messageWriterWriter) new(msg *golang.Message) error {
	w.Linef(`func New%v() %v {`, msg.Writer, msg.Writer)
	w.Linef(`return %v{baseproto.NewMessageWriter()}`, msg.Writer)
	w.Line(`}`)

	w.Linef(`func New%vBuffer(b buffer.Buffer) %v {`, msg.Name, msg.Writer)
	w.Linef(`return %v{baseproto.NewMessageWriterBuffer(b)}`, msg.Writer)
	w.Line(`}`)

	w.Linef(`func New%vWriterTo(w baseproto.MessageWriter) %v {`, msg.Name, msg.Writer)
	w.Linef(`return %v{w}`, msg.Writer)
	w.Line(`}`)
	w.Line()
	return nil
}

// fields

func (w *messageWriterWriter) fields(msg *golang.Message) error {
	for _, field := range msg.Fields {
		if err := w.field(msg, field); err != nil {
			return err
		}
	}

	w.Line()
	return nil
}

func (w *messageWriterWriter) field(msg *golang.Message, field *golang.MessageField) error {
	tag := field.Tag
	typ := field.Type

	switch t := typ.(type) {
	case golang.AnyType:
		w.Writef(`func (w %v) %v() baseproto.FieldWriter {`, msg.Writer, field.Name)
		w.Writef(`return w.w.Field(%d)`, tag)
		w.Line(`}`)

	case golang.AnyMessageType:
		w.Writef(`func (w %v) %v() baseproto.MessageWriter {`, msg.Writer, field.Name)
		w.Writef(`return w.w.Field(%d).Message()`, tag)
		w.Line(`}`)

	case golang.EnumType:
		encode := t.EncodeFunc()

		w.Writef(`func (w %v) %v(v %v) {`, msg.Writer, field.Name, t.InputName())
		w.Writef(`baseproto.WriteField(w.w.Field(%d), v, %v)`, tag, encode)
		w.Line(`}`)

	case golang.ListType:
		writer := t.Writer()
		newWriter := t.NewWriter()
		addElem := t.Elem().AddListElem()

		w.Linef(`func (w %v) %v() %v {`, msg.Writer, field.Name, writer)
		w.Linef(`w1 := w.w.Field(%d).List()`, tag)
		w.Linef(`return %v(w1, %v)`, newWriter, addElem)
		w.Line(`}`)

	case golang.MessageType:
		writer := t.Writer()
		newWriter := t.NewWriter()

		w.Linef(`func (w %v) %v() %v {`, msg.Writer, field.Name, writer)
		w.Linef(`w1 := w.w.Field(%d).Message()`, tag)
		w.Linef(`return %v(w1)`, newWriter)
		w.Line(`}`)

	case golang.StructType:
		encode := t.EncodeFunc()

		w.Writef(`func (w %v) %v(v %v) {`, msg.Writer, field.Name, t.InputName())
		w.Writef(`baseproto.WriteField(w.w.Field(%d), v, %v)`, tag, encode)
		w.Line(`}`)

	case golang.ValueType:
		w.Writef(`func (w %v) %v(v %v) {`, msg.Writer, field.Name, typ.InputName())
		t.WriteField(w, tag)
		w.Line(`}`)

	default:
		panic("unsupported field type")
	}
	return nil
}

// copy

func (w *messageWriterWriter) copy(msg *golang.Message) error {
	var n int

	for _, field := range msg.Fields {
		ok, err := w.copyField(msg, field)
		if err != nil {
			return err
		}
		if ok {
			n++
		}
	}

	if n > 0 {
		w.Line()
	}
	return nil
}

func (w *messageWriterWriter) copyField(msg *golang.Message, field *golang.MessageField) (
	bool, error) {

	tag := field.Tag
	typ := field.Type

	switch t := typ.(type) {
	case golang.AnyType:
		w.Writef(`func (w %v) Copy%v(v baseproto.Value) error {`, msg.Writer, field.Name)
		w.Writef(`return w.w.Field(%d).Any(v)`, tag)
		w.Line(`}`)

	case golang.AnyMessageType:
		w.Writef(`func (w %v) Copy%v(v baseproto.Message) error {`, msg.Writer, field.Name)
		w.Writef(`return w.w.Field(%d).Copy(v)`, tag)
		w.Line(`}`)

	case golang.MessageType:
		w.Linef(`func (w %v) Copy%v(v %v) error {`, msg.Writer, field.Name, t.InputName())
		w.Linef(`return w.w.Field(%d).Copy(v.Unwrap())`, tag)
		w.Line(`}`)

	default:
		return false, nil
	}

	return true, nil
}

// end

func (w *messageWriterWriter) end(msg *golang.Message) error {
	w.Linef(`func (w %v) Merge(msg %v) error {`, msg.Writer, msg.Name)
	w.Linef(`return w.w.Merge(msg.Unwrap())`)
	w.Line(`}`)
	w.Linef(`func (w %v) Unwrap() baseproto.MessageWriter {`, msg.Writer)
	w.Linef(`return w.w`)
	w.Line(`}`)
	w.Line()

	w.Linef(`func (w %v) End() error {`, msg.Writer)
	w.Linef(`return w.w.End()`)
	w.Line(`}`)

	w.Linef(`func (w %v) Build() (_ %v, err error) {`, msg.Writer, msg.Name)
	w.Linef(`bytes, err := w.w.Build()`)
	w.Linef(`if err != nil {
		return
	}`)
	w.Linef(`return Open%vErr(bytes)`, msg.Name)
	w.Line(`}`)
	w.Line()
	return nil
}
