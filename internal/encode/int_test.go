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

// Int16

func TestInt16__should_encode_decode_int16(t *testing.T) {
	b := buffer.New()
	EncodeInt16(b, math.MaxInt16)
	p := b.Bytes()

	v, n, err := DecodeInt16(p)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, n, b.Len())
	assert.Equal(t, int16(math.MaxInt16), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindInt16, kind)
	assert.Equal(t, size, len(p))
}

// Int32

func TestInt32__should_encode_decode_int32(t *testing.T) {
	b := buffer.New()
	EncodeInt32(b, math.MaxInt32)
	p := b.Bytes()

	v, n, err := DecodeInt32(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, int32(math.MaxInt32), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindInt32, kind)
	assert.Equal(t, size, len(p))
}

// Int64

func TestInt64__should_encode_decode_int64(t *testing.T) {
	b := buffer.New()
	EncodeInt64(b, math.MaxInt64)
	p := b.Bytes()

	v, n, err := DecodeInt64(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, int64(math.MaxInt64), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindInt64, kind)
	assert.Equal(t, size, len(p))
}

func TestInt64__should_encode_decode_int64_from_int32(t *testing.T) {
	b := buffer.New()
	EncodeInt32(b, math.MaxInt32)
	p := b.Bytes()

	v, n, err := DecodeInt64(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, int64(math.MaxInt32), v)
}
