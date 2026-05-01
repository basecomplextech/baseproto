// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"testing"

	"github.com/basecomplextech/baselibrary/bin"
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
	"github.com/stretchr/testify/assert"
)

// Bin64/128/256

func TestBin64__should_encode_decode_bin64(t *testing.T) {
	b := buffer.New()
	v := bin.Random64()
	EncodeBin64(b, v)
	p := b.Bytes()

	v1, n, err := DecodeBin64(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, v, v1)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindBin64, kind)
	assert.Equal(t, size, len(p))
}

func TestBin128__should_encode_decode_bin128(t *testing.T) {
	b := buffer.New()
	v := bin.Random128()
	EncodeBin128(b, v)
	p := b.Bytes()

	v1, n, err := DecodeBin128(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, v, v1)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindBin128, kind)
	assert.Equal(t, size, len(p))
}

func TestBin192__should_encode_decode_bin192(t *testing.T) {
	b := buffer.New()
	v := bin.Random192()
	EncodeBin192(b, v)
	p := b.Bytes()

	v1, n, err := DecodeBin192(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, v, v1)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindBin192, kind)
	assert.Equal(t, size, len(p))
}

func TestBin256__should_encode_decode_bin256(t *testing.T) {
	b := buffer.New()
	v := bin.Random256()
	EncodeBin256(b, v)
	p := b.Bytes()

	v1, n, err := DecodeBin256(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, v, v1)

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindBin256, kind)
	assert.Equal(t, size, len(p))
}
