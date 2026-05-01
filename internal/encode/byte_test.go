// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"testing"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
	"github.com/stretchr/testify/assert"
)

// DecodeBool

func TestBool__should_encode_decode_bool_value(t *testing.T) {
	b := []byte{byte(format.KindTrue)}
	v, n, err := DecodeBool(b)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, n, len(b))
	assert.Equal(t, true, v)

	b = []byte{byte(format.KindFalse)}
	v, n, err = DecodeBool(b)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, n, len(b))
	assert.Equal(t, false, v)

	kind, size, err := DecodeKindSize(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindFalse, kind)
	assert.Equal(t, size, len(b))
}

// DecodeByte

func TestByte__should_encode_decode_byte(t *testing.T) {
	b := buffer.New()
	EncodeByte(b, 1)
	p := b.Bytes()

	v, n, err := DecodeByte(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, byte(1), v)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindByte, kind)
	assert.Equal(t, size, len(p))
}
