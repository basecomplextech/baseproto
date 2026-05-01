// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import (
	"fmt"

	"github.com/basecomplextech/baseproto/internal/format"
	"github.com/basecomplextech/baseproto/internal/values"
)

// MessageType defines a message type.
type MessageType[T values.MessageType] interface {
	Type[T]
	MessageTypeDyn
}

// MessageTypeDyn defines a message type.
type MessageTypeDyn interface {
	TypeDyn

	// Fields returns the number of fields in the message.
	Fields() uint16

	// FieldAt returns the field at the given index.
	FieldAt(i uint16) MessageFieldDyn

	// FieldByTag returns the field with the given tag.
	FieldByTag(tag uint16) (MessageFieldDyn, bool)

	// FieldByName returns the field with the given name.
	FieldByName(name string) (MessageFieldDyn, bool)

	// Methods

	// VerifyMessage verifies a message against the message type, skips unknown fields.
	VerifyMessage(msg values.Message) error
}

// NewMessageFunc is a function that returns a new message.
type NewMessageFunc[T values.MessageType] func(values.Message) T

// New

// NewMessageType returns a new message type with the given fields.
// The methods panics on duplicate field tags or names.
func NewMessageType[T values.MessageType](new NewMessageFunc[T],
	fields ...MessageFieldDyn) MessageType[T] {

	msg, err := newMessageTypeErr(new, fields...)
	if err != nil {
		panic(err)
	}
	return msg
}

// NewMessageTypeErr returns a new message type with the given fields.
// The methods returns an error on duplicate field tags or names.
func NewMessageTypeErr[T values.MessageType](new NewMessageFunc[T], fields ...MessageFieldDyn) (
	MessageType[T], error) {

	return newMessageTypeErr(new, fields...)
}

// internal

var _ MessageType[values.MessageType] = (*messageType[values.MessageType])(nil)

type messageType[T values.MessageType] struct {
	new NewMessageFunc[T]

	fields []MessageFieldDyn
	tags   map[uint16]MessageFieldDyn
	names  map[string]MessageFieldDyn
}

func newMessageTypeErr[T values.MessageType](new NewMessageFunc[T], fields ...MessageFieldDyn) (
	*messageType[T], error) {

	// Make message
	m := &messageType[T]{
		new: new,

		fields: make([]MessageFieldDyn, 0, len(fields)),
		tags:   make(map[uint16]MessageFieldDyn, len(fields)),
		names:  make(map[string]MessageFieldDyn, len(fields)),
	}

	// Add fields
	for _, field := range fields {
		tag := field.Tag()
		name := field.Name()

		if _, ok := m.tags[tag]; ok {
			return nil, fmt.Errorf("duplicate field tag: %d", tag)
		}
		if _, ok := m.names[name]; ok {
			return nil, fmt.Errorf("duplicate field name: %s", name)
		}

		m.fields = append(m.fields, field)
		m.tags[tag] = field
		m.names[name] = field
	}

	return m, nil
}

// Kind returns the type kind.
func (t *messageType[T]) Kind() format.Kind {
	return format.KindMessage
}

// String returns the string representation of the type.
func (t *messageType[T]) String() string {
	return format.KindMessage.String()
}

// Fields

// Fields returns the number of fields in the message.
func (t *messageType[T]) Fields() uint16 {
	return uint16(len(t.fields))
}

// FieldAt returns the field at the given index.
func (t *messageType[T]) FieldAt(i uint16) MessageFieldDyn {
	return t.fields[i]
}

// FieldByTag returns the field with the given tag.
func (t *messageType[T]) FieldByTag(tag uint16) (MessageFieldDyn, bool) {
	field, ok := t.tags[tag]
	return field, ok
}

// FieldByName returns the field with the given name.
func (t *messageType[T]) FieldByName(name string) (MessageFieldDyn, bool) {
	field, ok := t.names[name]
	return field, ok
}

// Methods

// Open opens a value.
func (t *messageType[T]) Open(b []byte) (v T, n int, err error) {
	// TODO: Fix size
	msg, err := values.OpenMessageErr(b)
	if err != nil {
		return v, 0, err
	}
	return t.new(msg), len(b), nil
}

// Parse parses and verifies a value.
func (t *messageType[T]) Parse(b []byte) (v T, n int, err error) {
	msg, err := values.OpenMessageErr(b)
	if err != nil {
		return v, 0, err
	}
	if err := t.verify(msg); err != nil {
		return v, 0, err
	}
	return t.new(msg), n, nil
}

// Verify

// Verify verifies a value against the type.
func (t *messageType[T]) Verify(b []byte) error {
	msg, err := values.OpenMessageErr(b)
	if err != nil {
		return err
	}
	return t.verify(msg)
}

// VerifyRaw verifies a raw, possibly untruncated, value against the type.
func (t *messageType[T]) VerifyRaw(b []byte) error {
	msg, err := values.OpenMessageErr(b)
	if err != nil {
		return err
	}
	return t.verify(msg)
}

// VerifyMessage verifies a message against the message type, skips unknown fields.
func (t *messageType[T]) VerifyMessage(msg values.Message) error {
	return t.verify(msg)
}

// private

func (t *messageType[T]) verify(msg values.Message) error {
	for _, field := range t.fields {
		tag := field.Tag()

		b, err := msg.FieldRawErr(tag)
		if err != nil {
			return err
		}

		if err := field.VerifyRaw(b); err != nil {
			return err
		}
	}
	return nil
}
