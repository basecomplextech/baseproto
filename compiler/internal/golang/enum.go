// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"
	"strings"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
)

type Enum struct {
	Name   string
	Values []*EnumValue
}

type EnumValue struct {
	Tag    int
	Name   string
	String string
}

func NewEnum(def *model.Definition) (*Enum, error) {
	enum := def.Enum
	name := toUpperCamelCase(def.Name)

	// Make enum
	e := &Enum{
		Name:   name,
		Values: make([]*EnumValue, 0, len(enum.Values)),
	}

	// Make values
	for _, val := range enum.Values {
		v, err := newEnumValue(def, val)
		if err != nil {
			return nil, err
		}
		e.Values = append(e.Values, v)
	}

	return e, nil
}

// private

func newEnumValue(def *model.Definition, val *model.EnumValue) (*EnumValue, error) {
	name := toUpperCamelCase(val.Name)
	name = fmt.Sprintf("%v_%v", def.Name, name)
	str := strings.ToLower(val.Name)

	v := &EnumValue{
		Tag:    val.Number,
		Name:   name,
		String: str,
	}
	return v, nil
}
