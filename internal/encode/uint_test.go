// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"math"
	"testing"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
	"github.com/stretchr/testify/assert"
)

// Uint16

func TestUint16__should_encode_decode_int16(t *testing.T) {
	b := buffer.New()
	EncodeUint16(b, math.MaxUint16)
	p := b.Bytes()

	v, n, err := DecodeUint16(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, uint16(math.MaxUint16), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindUint16, kind)
	assert.Equal(t, size, len(p))
}

// Uint32

func TestUint32__should_encode_decode_int32(t *testing.T) {
	b := buffer.New()
	EncodeUint32(b, math.MaxUint32)
	p := b.Bytes()

	v, n, err := DecodeUint32(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, uint32(math.MaxUint32), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindUint32, kind)
	assert.Equal(t, size, len(p))
}

// Uint64

func TestUint64__should_encode_decode_int64(t *testing.T) {
	b := buffer.New()
	EncodeUint64(b, math.MaxUint64)
	p := b.Bytes()

	v, n, err := DecodeUint64(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, uint64(math.MaxUint64), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindUint64, kind)
	assert.Equal(t, size, len(p))
}

func TestUint64__should_encode_decode_uint64_from_uint32(t *testing.T) {
	b := buffer.New()
	EncodeUint32(b, math.MaxUint32)
	p := b.Bytes()

	v, n, err := DecodeUint64(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, uint64(math.MaxUint32), v)
}
