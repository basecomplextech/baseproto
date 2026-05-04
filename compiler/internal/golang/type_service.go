// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
)

type ServiceType interface {
	// Kind returns a type kind.
	Kind() model.Kind

	// Name returns a type name.
	Name() string
}

// internal

type serviceType struct {
	name string
	imp  string
}

func newServiceType(typ *model.Type) (ServiceType, error) {
	kind := typ.Kind
	if kind != model.KindService {
		panic("not service")
	}

	t := &serviceType{
		name: typ.Name,
		imp:  typ.ImportName,
	}
	return t, nil
}

// Kind returns a type kind.
func (t *serviceType) Kind() model.Kind {
	return model.KindService
}

// Name returns a type name.
func (t *serviceType) Name() string {
	if t.imp != "" {
		return fmt.Sprintf("%s.%s", t.imp, t.name)
	}
	return t.name
}
