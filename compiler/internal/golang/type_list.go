// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type ListType interface {
	Type

	// Elem returns an element type.
	Elem() Type

	// Funcs

	// DecodeFunc returns a decode list func.
	DecodeFunc() string
}

// internal

var _ ListType = (*listType)(nil)

type listType struct {
	elem Type
}

func newListType(typ *model.Type) (*listType, error) {
	kind := typ.Kind
	if kind != model.KindList {
		panic("not list kind")
	}

	elem, err := newType(typ.Element)
	if err != nil {
		return nil, err
	}

	t := &listType{elem: elem}
	return t, nil
}

// Kind returns a type kind.
func (t *listType) Kind() model.Kind {
	return model.KindList
}

// Name returns a type name.
func (t *listType) Name() string {
	elem := t.elem.Name()

	if t.elem.Kind() == model.KindMessage {
		return fmt.Sprintf("baseproto.MessageList[%v]", elem)
	}
	return fmt.Sprintf("baseproto.ValueList[%v]", elem)
}

// Elem returns an element type.
func (t *listType) Elem() Type {
	return t.elem
}

// Funcs

// DecodeFunc returns a decode list func.
func (t *listType) DecodeFunc() string {
	elem := t.elem.FieldOutput()
	return fmt.Sprintf("baseproto.ParseTypedList[%v]", elem)
}

// DecodeListElem returns a decode func for a list element.
func (t *listType) DecodeListElem() string {
	elem := t.elem.DecodeListElem()

	if t.elem.Kind() == model.KindMessage {
		return fmt.Sprintf("baseproto.OpenMessageListErr[%v]", elem)
	}
	return fmt.Sprintf("baseproto.OpenValueListErr[%v]", elem)
}

// Fields

// FieldInput returns an input field type, i.e. string (not baseproto.String).
func (t *listType) FieldInput() string {
	panic("unsupported list field input")
}

// FieldOutput returns an output field type, i.e. baseproto.String (not string).
func (t *listType) FieldOutput() string {
	elem := t.elem.FieldOutput()

	if t.elem.Kind() == model.KindMessage {
		return fmt.Sprintf("baseproto.MessageList[%v]", elem)
	}
	return fmt.Sprintf("baseproto.ValueList[%v]", elem)
}

// Write fields

// ReturnField writes a field get.
func (t *listType) ReturnField(w writer.Writer, tag int) error {
	kind := t.elem.Kind()
	decode := t.elem.DecodeListElem()

	switch kind {
	case model.KindMessage:
		w.Writef(`return baseproto.NewMessageList(m.msg.List(%d), %v)`, tag, decode)
	default:
		w.Writef(`return baseproto.NewValueList(m.msg.List(%d), %v)`, tag, decode)
	}

	return nil
}

// WriteField writes a field write.
func (t *listType) WriteField(w writer.Writer, tag int) error {
	return nil
}
