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

	// Write

	// Write returns a list writer type.
	Writer() string

	// NewWriter returns a new list writer func.
	NewWriter() string
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

// InputName returns an input type, i.e. string (not baseproto.String).
func (t *listType) InputName() string {
	panic("unsupported list field input")
}

// OutputName returns an output type, i.e. baseproto.String (not string).
func (t *listType) OutputName() string {
	elem := t.elem.OutputName()

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
	elem := t.elem.OutputName()
	return fmt.Sprintf("baseproto.ParseTypedList[%v]", elem)
}

// ParseFunc returns a parse func.
func (t *listType) ParseFunc() string {
	return t.DecodeFunc()
}

// List

// AddListElem returns an encode func for a list element.
func (t *listType) AddListElem() string {
	elem := t.elem
	if elem.Kind() == model.KindMessage {
		return "baseproto.NewMessageListWriter"
	}
	return "baseproto.NewValueListWriter"
}

// GetListElem returns a decode func for a list element.
func (t *listType) GetListElem() string {
	elem := t.elem.GetListElem()

	if t.elem.Kind() == model.KindMessage {
		return fmt.Sprintf("baseproto.OpenMessageListErr[%v]", elem)
	}
	return fmt.Sprintf("baseproto.OpenValueListErr[%v]", elem)
}

// Message

// GetField writes a field get.
func (t *listType) GetField(w writer.Writer, tag int) error {
	kind := t.elem.Kind()
	decode := t.elem.GetListElem()

	switch kind {
	case model.KindMessage:
		w.Writef(`return baseproto.NewMessageList(m.msg.List(%d), %v)`, tag, decode)
	default:
		w.Writef(`return baseproto.NewValueList(m.msg.List(%d), %v)`, tag, decode)
	}

	return nil
}

// Write

// Write returns a list writer type.
func (t *listType) Writer() string {
	elem := t.elem

	if m, ok := elem.(MessageType); ok {
		writer := m.Writer()
		return fmt.Sprintf("baseproto.MessageListWriter[%v]", writer)
	}

	elemName := elem.InputName()
	return fmt.Sprintf("baseproto.ValueListWriter[%v]", elemName)
}

// NewWriter returns a new list writer func.
func (t *listType) NewWriter() string {
	elem := t.elem
	if elem.Kind() == model.KindMessage {
		return "baseproto.NewMessageListWriter"
	}
	return "baseproto.NewValueListWriter"
}
