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

// String

func TestString__should_encode_decode_string(t *testing.T) {
	v := "hello, world"

	b := buffer.New()
	_, err := EncodeString(b, v)
	if err != nil {
		t.Fatal(err)
	}
	p := b.Bytes()

	v1, n, err := DecodeString(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, b.Len())
	assert.Equal(t, v, v1.Unwrap())

	kind, size, err := DecodeKindSize(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindString, kind)
	assert.Equal(t, size, len(p))
}
