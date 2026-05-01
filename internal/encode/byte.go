// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"errors"
	"fmt"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeBool(b buffer.Buffer, v bool) (int, error) {
	p := b.Grow(1)
	if v {
		p[0] = byte(format.KindTrue)
	} else {
		p[0] = byte(format.KindFalse)
	}
	return 1, nil
}

func EncodeByte(b buffer.Buffer, v byte) (int, error) {
	p := b.Grow(2)
	p[0] = v
	p[1] = byte(format.KindByte)
	return 2, nil
}

// Decode

func DecodeByte(b []byte) (byte, int, error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		return 0, 0, errors.New("decode byte: invalid data")
	}
	if kind != format.KindByte {
		return 0, 0, fmt.Errorf("decode byte: invalid kind, kind=%v", kind)
	}

	end := len(b) - 2
	if end < 0 {
		return 0, 0, errors.New("decode byte: invalid data")
	}
	return b[end], 2, nil
}

func DecodeBool(b []byte) (bool, int, error) {
	if len(b) == 0 {
		return false, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		return false, 0, errors.New("decode bool: invalid data")
	}

	v := kind == format.KindTrue
	size := n
	return v, size, nil
}
