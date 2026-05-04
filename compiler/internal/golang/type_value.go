// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type ValueType interface {
	Type
	value()

	// Funcs

	// DecodeFunc returns a decode func.
	DecodeFunc() string
}

// internal

var _ ValueType = (*valueType)(nil)

type valueType struct {
	kind model.Kind
	name string
}

func newValueType(typ *model.Type) (*valueType, error) {
	kind := typ.Kind
	t := &valueType{
		kind: kind,
	}

	switch t.kind {
	case model.KindBool:
		t.name = "bool"

	case model.KindByte:
		t.name = "byte"

	case model.KindInt16:
		t.name = "int16"
	case model.KindInt32:
		t.name = "int32"
	case model.KindInt64:
		t.name = "int64"

	case model.KindUint16:
		t.name = "uint16"
	case model.KindUint32:
		t.name = "uint32"
	case model.KindUint64:
		t.name = "uint64"

	case model.KindBin64:
		t.name = "bin.Bin64"
	case model.KindBin128:
		t.name = "bin.Bin128"
	case model.KindBin192:
		t.name = "bin.Bin192"
	case model.KindBin256:
		t.name = "bin.Bin256"

	case model.KindFloat32:
		t.name = "float32"
	case model.KindFloat64:
		t.name = "float64"

	case model.KindBytes:
		t.name = "[]byte"
	case model.KindString:
		t.name = "string"

	default:
		panic("unsupported value type")
	}

	return t, nil
}

func (t *valueType) value() {}

// Kind returns a type kind.
func (t *valueType) Kind() model.Kind {
	return t.kind
}

// Name returns a type name.
func (t *valueType) Name() string {
	return t.name
}

// Funcs

// DecodeFunc returns a decode func.
func (t *valueType) DecodeFunc() string {
	switch t.kind {
	case model.KindAny:
		return "baseproto.ParseValue"

	case model.KindBool:
		return "baseproto.DecodeBool"
	case model.KindByte:
		return "baseproto.DecodeByte"

	case model.KindInt16:
		return "baseproto.DecodeInt16"
	case model.KindInt32:
		return "baseproto.DecodeInt32"
	case model.KindInt64:
		return "baseproto.DecodeInt64"

	case model.KindUint16:
		return "baseproto.DecodeUint16"
	case model.KindUint32:
		return "baseproto.DecodeUint32"
	case model.KindUint64:
		return "baseproto.DecodeUint64"

	case model.KindFloat32:
		return "baseproto.DecodeFloat32"
	case model.KindFloat64:
		return "baseproto.DecodeFloat64"

	case model.KindBin64:
		return "baseproto.DecodeBin64"
	case model.KindBin128:
		return "baseproto.DecodeBin128"
	case model.KindBin192:
		return "baseproto.DecodeBin192"
	case model.KindBin256:
		return "baseproto.DecodeBin256"

	case model.KindBytes:
		return "baseproto.DecodeBytes"
	case model.KindString:
		return "baseproto.DecodeString"
	}

	panic("unsupported value type")
}

// DecodeListElem returns a decode func for a list element.
func (t *valueType) DecodeListElem() string {
	return t.DecodeFunc()
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *valueType) FieldInput() string {
	return t.name
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *valueType) FieldOutput() string {
	switch t.kind {
	case model.KindBytes:
		return "baseproto.Bytes"
	case model.KindString:
		return "baseproto.String"
	}

	return t.name
}

// Write fields

// ReturnField writes a field get.
func (t *valueType) ReturnField(w writer.Writer, tag int) error {
	switch t.kind {
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

	default:
		panic("unsupported value type")
	}

	return nil
}

// WriteField writes a field write.
func (t *valueType) WriteField(w writer.Writer, tag int) error {
	switch t.kind {
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

	default:
		panic("unsupported value type")
	}

	return nil
}
