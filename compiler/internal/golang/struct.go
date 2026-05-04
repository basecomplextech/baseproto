// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import "github.com/basecomplextech/baseproto/compiler/internal/model"

type Struct struct {
	Name   string
	Fields []*StructField
}

type StructField struct {
	Name string
	Type EncodableType
}

func NewStruct(def *model.Definition) (*Struct, error) {
	str := def.Struct
	name := def.Name

	// Make struct
	s := &Struct{
		Name:   name,
		Fields: make([]*StructField, 0, str.Fields.Len()),
	}

	// Make fields
	for _, field := range str.Fields.Values() {
		f, err := newStructField(field)
		if err != nil {
			return nil, err
		}
		s.Fields = append(s.Fields, f)
	}

	return s, nil
}

// private

func newStructField(field *model.StructField) (*StructField, error) {
	name := toUpperCamelCase(field.Name)

	typ, err := newType(field.Type)
	if err != nil {
		return nil, err
	}
	typ1, ok := typ.(EncodableType)
	if !ok {
		panic("not encodable type")
	}

	f := &StructField{
		Name: name,
		Type: typ1,
	}
	return f, nil
}
