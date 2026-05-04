// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import "github.com/basecomplextech/baseproto/compiler/internal/model"

type MessageField struct {
	Tag  int
	Name string
	Type Type
}

func newMessageField(field *model.Field) (*MessageField, error) {
	name := toUpperCamelCase(field.Name)

	// Make type
	typ, err := newType(field.Type)
	if err != nil {
		return nil, err
	}

	// Make field
	f := &MessageField{
		Tag:  field.Tag,
		Name: name,
		Type: typ,
	}
	return f, nil
}
