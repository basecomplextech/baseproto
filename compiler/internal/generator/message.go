// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

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
	if err := w.def(def); err != nil {
		return err
	}
	if err := w.new_methods(def); err != nil {
		return err
	}
	if err := w.parse_method(def); err != nil {
		return err
	}
	if err := w.fields(def); err != nil {
		return err
	}
	if err := w.has_fields(def); err != nil {
		return err
	}
	if err := w.methods(def); err != nil {
		return err
	}
	return nil
}

func (w *messageWriter) def(def *model.Definition) error {
	w.Linef(`// %v`, def.Name)
	w.Line()
	w.Linef(`type %v struct {`, def.Name)
	w.Line(`msg baseproto.Message`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) new_methods(def *model.Definition) error {
	w.Linef(`func New%v(msg baseproto.Message) %v {`, def.Name, def.Name)
	w.Linef(`return %v{msg}`, def.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func Open%v(b []byte) %v {`, def.Name, def.Name)
	w.Linef(`msg := baseproto.OpenMessage(b)`)
	w.Linef(`return %v{msg}`, def.Name)
	w.Linef(`}`)
	w.Line()

	w.Linef(`func Open%vErr(b []byte) (_ %v, err error) {`, def.Name, def.Name)
	w.Linef(`msg, err := baseproto.OpenMessageErr(b)`)
	w.Linef(`return %v{msg}, err`, def.Name)
	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) parse_method(def *model.Definition) error {
	w.Linef(`func Parse%v(b []byte) (_ %v, size int, err error) {`, def.Name, def.Name)
	w.Linef(`msg, size, err := baseproto.ParseMessage(b)`)
	w.Linef(`return %v{msg}, size, err`, def.Name)
	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) fields(def *model.Definition) error {
	fields := def.Message.Fields.List

	for _, field := range fields {
		if err := w.field(def, field); err != nil {
			return err
		}
	}

	if len(fields) > 1 {
		w.Line()
	}
	return nil
}

func (w *messageWriter) field(def *model.Definition, field *model.Field) error {
	fieldName := messageFieldName(field)
	typeName := typeRefName(field.Type)

	tag := field.Tag
	kind := field.Type.Kind

	switch kind {
	default:
		w.Writef(`func (m %v) %v() %v {`, def.Name, fieldName, typeName)

		switch kind {
		case model.KindBool:
			w.Writef(`return m.msg.Bool(%d)`, tag)
		case model.KindByte:
			w.Writef(`return m.msg.Byte(%d)`, tag)

		case model.KindInt16:
			w.Writef(`return m.msg.Int16(%d)`, tag)
		case model.KindInt32:
			w.Writef(`return m.msg.Int32(%d)`, tag)
		case model.KindInt64:
			w.Writef(`return m.msg.Int64(%d)`, tag)

		case model.KindUint16:
			w.Writef(`return m.msg.Uint16(%d)`, tag)
		case model.KindUint32:
			w.Writef(`return m.msg.Uint32(%d)`, tag)
		case model.KindUint64:
			w.Writef(`return m.msg.Uint64(%d)`, tag)

		case model.KindBin64:
			w.Writef(`return m.msg.Bin64(%d)`, tag)
		case model.KindBin128:
			w.Writef(`return m.msg.Bin128(%d)`, tag)
		case model.KindBin192:
			w.Writef(`return m.msg.Bin192(%d)`, tag)
		case model.KindBin256:
			w.Writef(`return m.msg.Bin256(%d)`, tag)

		case model.KindFloat32:
			w.Writef(`return m.msg.Float32(%d)`, tag)
		case model.KindFloat64:
			w.Writef(`return m.msg.Float64(%d)`, tag)

		case model.KindBytes:
			w.Writef(`return m.msg.Bytes(%d)`, tag)
		case model.KindString:
			w.Writef(`return m.msg.String(%d)`, tag)

		case model.KindAny:
			w.Writef(`return m.msg.Field(%d)`, tag)
		case model.KindAnyMessage:
			w.Writef(`return m.msg.Field(%d).Message()`, tag)
		}

		w.Writef(`}`)
		w.Line()

	case model.KindList:
		elem := field.Type.Element
		decodeFunc := typeDecodeRefFunc(field.Type.Element)

		w.Writef(`func (m %v) %v() %v {`, def.Name, fieldName, typeName)
		if elem.Kind == model.KindMessage {
			w.Writef(`return baseproto.NewMessageList(m.msg.List(%d), %v)`, tag, decodeFunc)
		} else {
			w.Writef(`return baseproto.NewValueList(m.msg.List(%d), %v)`, tag, decodeFunc)
		}

		w.Writef(`}`)
		w.Line()

	case model.KindMessage:
		makeFunc := typeMakeMessageFunc(field.Type)

		w.Writef(`func (m %v) %v() %v {`, def.Name, fieldName, typeName)
		w.Writef(`return %v(m.msg.Message(%d))`, makeFunc, tag)
		w.Writef(`}`)
		w.Line()

	case model.KindEnum,
		model.KindStruct:
		newFunc := typeNewFunc(field.Type)

		w.Writef(`func (m %v) %v() %v {`, def.Name, fieldName, typeName)
		w.Writef(`return %v(m.msg.FieldRaw(%d))`, newFunc, tag)
		w.Writef(`}`)
		w.Line()
	}
	return nil
}

func (w *messageWriter) has_fields(def *model.Definition) error {
	fields := def.Message.Fields.List

	for _, field := range fields {
		if err := w.has_field(def, field); err != nil {
			return err
		}
	}

	if len(fields) > 1 {
		w.Line()
	}
	return nil
}

func (w *messageWriter) has_field(def *model.Definition, field *model.Field) error {
	fieldName := messageFieldName(field)
	tag := field.Tag

	w.Writef(`func (m %v) Has%v() bool {`, def.Name, fieldName)
	w.Writef(`return m.msg.HasField(%d)`, tag)
	w.Writef(`}`)
	w.Line()
	return nil
}

func (w *messageWriter) methods(def *model.Definition) error {
	w.Writef(`func (m %v) Clone() %v {`, def.Name, def.Name)
	w.Writef(`return %v{m.msg.Clone()}`, def.Name)
	w.Writef(`}`)
	w.Line()

	w.Writef(`func (m %v) CloneToArena(a alloc.Arena) %v {`, def.Name, def.Name)
	w.Writef(`return %v{m.msg.CloneToArena(a)}`, def.Name)
	w.Writef(`}`)
	w.Line()

	w.Writef(`func (m %v) CloneToBuffer(b buffer.Buffer) %v {`, def.Name, def.Name)
	w.Writef(`return %v{m.msg.CloneToBuffer(b)}`, def.Name)
	w.Writef(`}`)
	w.Line()

	w.Line()

	w.Writef(`func (m %v) IsEmpty() bool {`, def.Name)
	w.Writef(`return m.msg.Empty()`)
	w.Writef(`}`)
	w.Line()

	w.Writef(`func (m %v) Unwrap() baseproto.Message {`, def.Name)
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

func (w *messageWriter) writer_fields(def *model.Definition) error {
	fields := def.Message.Fields.List

	for _, field := range fields {
		if err := w.writer_field(def, field); err != nil {
			return err
		}
	}

	w.Line()
	return nil
}

func (w *messageWriter) writer_field(def *model.Definition, field *model.Field) error {
	fname := messageFieldName(field)
	tname := inTypeName(field.Type)
	wname := fmt.Sprintf("%vWriter", def.Name)

	tag := field.Tag
	kind := field.Type.Kind

	switch kind {
	default:
		w.Writef(`func (w %vWriter) %v(v %v) {`, def.Name, fname, tname)

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
