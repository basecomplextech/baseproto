// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package types

import (
	"github.com/basecomplextech/baselibrary/bin"
	"github.com/basecomplextech/baseproto/internal/encode"
	"github.com/basecomplextech/baseproto/internal/format"
)

var (
	Bool = &valueType[bool]{
		kind:   format.KindTrue,
		decode: encode.DecodeBool,
	}
	Byte = &valueType[byte]{
		kind:   format.KindByte,
		decode: encode.DecodeByte,
	}

	Int16 = &valueType[int16]{
		kind:   format.KindInt16,
		decode: encode.DecodeInt16,
	}
	Int32 = &valueType[int32]{
		kind:   format.KindInt32,
		decode: encode.DecodeInt32,
	}
	Int64 = &valueType[int64]{
		kind:   format.KindInt64,
		decode: encode.DecodeInt64,
	}

	Uint16 = &valueType[uint16]{
		kind:   format.KindUint16,
		decode: encode.DecodeUint16,
	}
	Uint32 = &valueType[uint32]{
		kind:   format.KindUint32,
		decode: encode.DecodeUint32,
	}
	Uint64 = &valueType[uint64]{
		kind:   format.KindUint64,
		decode: encode.DecodeUint64,
	}

	Float32 = &valueType[float32]{
		kind:   format.KindFloat32,
		decode: encode.DecodeFloat32,
	}
	Float64 = &valueType[float64]{
		kind:   format.KindFloat64,
		decode: encode.DecodeFloat64,
	}

	Bin64 = &valueType[bin.Bin64]{
		kind:   format.KindBin64,
		decode: encode.DecodeBin64,
	}
	Bin128 = &valueType[bin.Bin128]{
		kind:   format.KindBin128,
		decode: encode.DecodeBin128,
	}
	Bin192 = &valueType[bin.Bin192]{
		kind:   format.KindBin192,
		decode: encode.DecodeBin192,
	}
	Bin256 = &valueType[bin.Bin256]{
		kind:   format.KindBin256,
		decode: encode.DecodeBin256,
	}

	Bytes = &valueType[format.Bytes]{
		kind:   format.KindBytes,
		decode: encode.DecodeBytes,
	}
	String = &valueType[format.String]{
		kind:   format.KindString,
		decode: encode.DecodeString,
	}
)

// private

var _ Type[any] = (*valueType[any])(nil)

type valueType[T any] struct {
	kind   format.Kind
	decode func(b []byte) (_ T, n int, err error)
}

// Kind returns the type kind.
func (t *valueType[T]) Kind() format.Kind {
	return t.kind
}

// String returns the string representation of the type.
func (t *valueType[T]) String() string {
	return t.kind.String()
}

// Methods

// Open opens a value.
func (t *valueType[T]) Open(b []byte) (v T, n int, err error) {
	return t.decode(b)
}

// Parse parses and verifies a value.
func (t *valueType[T]) Parse(b []byte) (v T, n int, err error) {
	return t.decode(b)
}

// Verify

// Verify verifies a value against the type.
func (t *valueType[T]) Verify(b []byte) error {
	_, _, err := t.decode(b)
	return err
}

// VerifyRaw verifies a raw, possibly untruncated, value against the type.
func (t *valueType[T]) VerifyRaw(b []byte) error {
	_, _, err := t.decode(b)
	return err
}

// Internal

// Resolve resolves internal field type references.
func (t *valueType[T]) Resolve() {}
