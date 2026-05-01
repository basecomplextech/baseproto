// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import (
	"fmt"

	"github.com/basecomplextech/baseproto/internal/format"
	"github.com/basecomplextech/baseproto/internal/values"
)

// ListType defines a list type.
type ListType[T any] interface {
	Type[T]
	ListTypeDyn

	// Elem returns the element type of the list.
	Elem() Type[T]
}

// ListTypeDyn defines a list type.
type ListTypeDyn interface {
	TypeDyn

	// ElemDyn returns the element type of the list.
	ElemDyn() TypeDyn
}

// NewListFunc is a function that returns a new list.
type NewListFunc[T any] func(list values.List) T

// New

// NewListType returns a new list type.
func NewListType[T any](new NewListFunc[T], elem Type[T]) ListType[T] {
	return newListType(new, elem, nil)
}

// NewListTypePtr returns a new list type with an element type pointer.
func NewListTypePtr[T any](new NewListFunc[T], elemPtr *Type[T]) ListType[T] {
	return newListType(new, nil, elemPtr)
}

// internal

var _ ListType[any] = (*listType[any])(nil)

type listType[T any] struct {
	new     NewListFunc[T]
	elem    Type[T]
	elemPtr *Type[T]
}

func newListType[T any](new NewListFunc[T], elem Type[T], elemPtr *Type[T]) *listType[T] {
	return &listType[T]{
		new:     new,
		elem:    elem,
		elemPtr: elemPtr,
	}
}

// Kind returns the type kind.
func (t *listType[T]) Kind() format.Kind {
	return format.KindList
}

// String returns the string representation of the type.
func (t *listType[T]) String() string {
	return fmt.Sprintf("list<%s>", t.elem.String())
}

// Elements

// Elem returns the element type of the list.
func (t *listType[T]) Elem() Type[T] {
	return t.elem
}

// ElemDyn returns the element type of the list.
func (t *listType[T]) ElemDyn() TypeDyn {
	return t.elem
}

// Methods

// Open opens a value.
func (t *listType[T]) Open(b []byte) (v T, n int, err error) {
	// TODO: Fix size
	list, err := values.OpenListErr(b)
	if err != nil {
		return v, n, err
	}
	return t.new(list), n, nil
}

// Parse parses and verifies a value.
func (t *listType[T]) Parse(b []byte) (v T, n int, err error) {
	list, err := values.OpenListErr(b)
	if err != nil {
		return v, n, err
	}
	return t.new(list), n, nil
}

// Verify

// Verify verifies a value against the type.
func (t *listType[T]) Verify(b []byte) error {
	return t.VerifyRaw(b)
}

// VerifyRaw verifies a raw, possibly untruncated, value against the type.
func (t *listType[T]) VerifyRaw(b []byte) error {
	list, err := values.OpenListErr(b)
	if err != nil {
		return err
	}
	return t.verify(list)
}

// Internal

// Resolve resolves internal type references.
func (t *listType[T]) Resolve() {
	if t.elem != nil {
		return
	}

	t.elem = *t.elemPtr
}

// private

func (t *listType[T]) verify(list values.List) error {
	num := list.Len()
	for i := range num {
		b := list.GetBytes(i)

		if err := t.elem.VerifyRaw(b); err != nil {
			return fmt.Errorf("invalid list element %d: %w", i, err)
		}
	}
	return nil
}
