// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

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

func (w *messageWriter) message(def *model.Definition) error {
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

func (w *messageWriter) messageWriter(def *model.Definition) error {
	if err := w.writer_def(def); err != nil {
		return err
	}
	if err := w.writer_new_method(def); err != nil {
		return err
	}
	if err := w.writer_fields(def); err != nil {
		return err
	}
	if err := w.writer_end(def); err != nil {
		return err
	}
	return nil
}

func (w *messageWriter) writer_def(def *model.Definition) error {
	w.Linef(`// %vWriter`, def.Name)
	w.Line()
	w.Linef(`type %vWriter struct {`, def.Name)
	w.Line(`w baseproto.MessageWriter`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) writer_new_method(def *model.Definition) error {
	w.Linef(`func New%vWriter() %vWriter {`, def.Name, def.Name)
	w.Linef(`w := baseproto.NewMessageWriter()`)
	w.Linef(`return %vWriter{w}`, def.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func New%vWriterBuffer(b buffer.Buffer) %vWriter {`, def.Name, def.Name)
	w.Linef(`w := baseproto.NewMessageWriterBuffer(b)`)
	w.Linef(`return %vWriter{w}`, def.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func New%vWriterTo(w baseproto.MessageWriter) %vWriter {`, def.Name, def.Name)
	w.Linef(`return %vWriter{w}`, def.Name)
	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) writer_end(def *model.Definition) error {
	w.Linef(`func (w %vWriter) Merge(msg %v) error {`, def.Name, def.Name)
	w.Linef(`return w.w.Merge(msg.Unwrap())`)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func (w %vWriter) End() error {`, def.Name)
	w.Linef(`return w.w.End()`)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func (w %vWriter) Build() (_ %v, err error) {`, def.Name, def.Name)
	w.Linef(`bytes, err := w.w.Build()`)
	w.Linef(`if err != nil {
		return
	}`)
	w.Linef(`return Open%vErr(bytes)`, def.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func (w %vWriter) Unwrap() baseproto.MessageWriter {`, def.Name)
	w.Linef(`return w.w`)
	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) writer_fields(msg *model.Definition) error {
	fields := msg.Message.Fields.List

	for _, field := range fields {
		if err := w.writer_field(msg, field); err != nil {
			return err
		}
	}

	w.Line()
	return nil
}

func (w *messageWriter) writer_field(msg *model.Definition, field *model.Field) error {
	fname := messageFieldName(field)
	tname := inTypeName(field.Type)
	wname := fmt.Sprintf("%vWriter", msg.Name)

	tag := field.Tag
	kind := field.Type.Kind

	switch kind {
	default:
		w.Writef(`func (w %vWriter) %v(v %v) {`, msg.Name, fname, tname)

		switch kind {
		case model.KindBool:
			w.Writef(`w.w.Field(%d).Bool(v)`, tag)
		case model.KindByte:
			w.Writef(`w.w.Field(%d).Byte(v)`, tag)

		case model.KindInt16:
			w.Writef(`w.w.Field(%d).Int16(v)`, tag)
		case model.KindInt32:
			w.Writef(`w.w.Field(%d).Int32(v)`, tag)
		case model.KindInt64:
			w.Writef(`w.w.Field(%d).Int64(v)`, tag)

		case model.KindUint16:
			w.Writef(`w.w.Field(%d).Uint16(v)`, tag)
		case model.KindUint32:
			w.Writef(`w.w.Field(%d).Uint32(v)`, tag)
		case model.KindUint64:
			w.Writef(`w.w.Field(%d).Uint64(v)`, tag)

		case model.KindBin64:
			w.Writef(`w.w.Field(%d).Bin64(v)`, tag)
		case model.KindBin128:
			w.Writef(`w.w.Field(%d).Bin128(v)`, tag)
		case model.KindBin192:
			w.Writef(`w.w.Field(%d).Bin192(v)`, tag)
		case model.KindBin256:
			w.Writef(`w.w.Field(%d).Bin256(v)`, tag)

		case model.KindFloat32:
			w.Writef(`w.w.Field(%d).Float32(v)`, tag)
		case model.KindFloat64:
			w.Writef(`w.w.Field(%d).Float64(v)`, tag)

		case model.KindBytes:
			w.Writef(`w.w.Field(%d).Bytes(v)`, tag)
		case model.KindString:
			w.Writef(`w.w.Field(%d).String(v)`, tag)
		}
		w.Linef(`}`)

	case model.KindAny:
		w.Writef(`func (w %v) %v() baseproto.FieldWriter {`, wname, fname)
		w.Writef(`return w.w.Field(%d)`, tag)
		w.Linef(`}`)

		w.Writef(`func (w %v) Copy%v(v baseproto.Value) error {`, wname, fname)
		w.Writef(`return w.w.Field(%d).Any(v)`, tag)
		w.Linef(`}`)

	case model.KindAnyMessage:
		w.Writef(`func (w %v) %v() baseproto.MessageWriter {`, wname, fname)
		w.Writef(`return w.w.Field(%d).Message()`, tag)
		w.Linef(`}`)

		w.Writef(`func (w %v) Copy%v(v baseproto.Message) error {`, wname, fname)
		w.Writef(`return w.w.Field(%d).Any(v.Raw())`, tag)
		w.Linef(`}`)

	case model.KindEnum:
		writeFunc := typeWriteFunc(field.Type)

		w.Writef(`func (w %v) %v(v %v) {`, wname, fname, tname)
		w.Writef(`baseproto.WriteField(w.w.Field(%d), v, %v)`, tag, writeFunc)
		w.Linef(`}`)

	case model.KindStruct:
		writeFunc := typeWriteFunc(field.Type)

		w.Writef(`func (w %v) %v(v %v) {`, wname, fname, tname)
		w.Writef(`baseproto.WriteField(w.w.Field(%d), v, %v)`, tag, writeFunc)
		w.Linef(`}`)

	case model.KindList:
		writer := typeWriter(field.Type)
		buildList := typeWriteFunc(field.Type)
		encodeElement := typeWriteFunc(field.Type.Element)

		w.Linef(`func (w %v) %v() %v {`, wname, fname, writer)
		w.Linef(`w1 := w.w.Field(%d).List()`, tag)
		w.Linef(`return %v(w1, %v)`, buildList, encodeElement)
		w.Linef(`}`)

	case model.KindMessage:
		writer := typeWriter(field.Type)
		writer_new_method := typeWriteFunc(field.Type)
		w.Linef(`func (w %v) %v() %v {`, wname, fname, writer)
		w.Linef(`w1 := w.w.Field(%d).Message()`, tag)
		w.Linef(`return %v(w1)`, writer_new_method)
		w.Linef(`}`)

		tname := typeName(field.Type)
		w.Linef(`func (w %v) Copy%v(v %v) error {`, wname, fname, tname)
		w.Linef(`return w.w.Field(%d).Any(v.Unwrap().Raw())`, tag)
		w.Linef(`}`)
	}
	return nil
}

// util

func messageFieldName(field *model.Field) string {
	return toUpperCamelCase(field.Name)
}
