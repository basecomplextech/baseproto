// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
)

type Service struct {
	Name string
	Sub  bool // subservice

	Handler     string
	HandlerImpl string

	Methods []*Method
}

func NewService(def *model.Definition) (*Service, error) {
	srv := def.Service
	handler := fmt.Sprintf("%vHandler", toUpperCamelCase(def.Name))
	handlerImpl := fmt.Sprintf("%vHandler", toLowerCamelCase(def.Name))

	// Make service
	s := &Service{
		Name: def.Name,
		Sub:  srv.Sub,

		Handler:     handler,
		HandlerImpl: handlerImpl,
	}

	// Make methods
	for _, m := range srv.Methods {
		method, err := newMethod(srv, m)
		if err != nil {
			return nil, err
		}
		s.Methods = append(s.Methods, method)
	}
	return s, nil
}
