// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"errors"
	"fmt"

	"github.com/basecomplextech/baselibrary/bin"
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeBin64(b buffer.Buffer, v bin.Bin64) (int, error) {
	p := b.Grow(9)
	v.MarshalTo(p)
	p[8] = byte(format.KindBin64)
	return 9, nil
}

func EncodeBin128(b buffer.Buffer, v bin.Bin128) (int, error) {
	p := b.Grow(17)
	v.MarshalTo(p)
	p[16] = byte(format.KindBin128)
	return 17, nil
}

func EncodeBin128Bytes(b buffer.Buffer, v bin.Bin128) ([]byte, int, error) {
	p := b.Grow(17)
	v.MarshalTo(p)
	p[16] = byte(format.KindBin128)
	return p, 17, nil
}

func EncodeBin192(b buffer.Buffer, v bin.Bin192) (int, error) {
	p := b.Grow(25)
	v.MarshalTo(p)
	p[24] = byte(format.KindBin192)
	return 25, nil
}

func EncodeBin256(b buffer.Buffer, v bin.Bin256) (int, error) {
	p := b.Grow(33)
	v.MarshalTo(p)
	p[32] = byte(format.KindBin256)
	return 33, nil
}

// Decode

func DecodeBin64(b []byte) (_ bin.Bin64, size int, err error) {
	if len(b) == 0 {
		return bin.Bin64{}, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode bin64: invalid data")
		return
	}
	if kind != format.KindBin64 {
		err = fmt.Errorf("decode bin64: invalid kind, kind=%v", kind)
		return
	}

	size = n
	start := len(b) - (n + 8)
	end := len(b) - n

	if start < 0 {
		err = errors.New("decode bin64: invalid data")
		return
	}

	v, err := bin.Parse64(b[start:end])
	if err != nil {
		return
	}

	size += 8
	return v, size, nil
}

func DecodeBin128(b []byte) (_ bin.Bin128, size int, err error) {
	if len(b) == 0 {
		return bin.Bin128{}, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode bin128: invalid data")
		return
	}
	if kind != format.KindBin128 {
		err = fmt.Errorf("decode bin128: invalid kind, kind=%v", kind)
		return
	}

	size = n
	start := len(b) - (n + 16)
	end := len(b) - n

	if start < 0 {
		err = errors.New("decode bin128: invalid data")
		return
	}

	v, err := bin.Parse128(b[start:end])
	if err != nil {
		return
	}

	size += 16
	return v, size, nil
}

func DecodeBin192(b []byte) (_ bin.Bin192, size int, err error) {
	if len(b) == 0 {
		return bin.Bin192{}, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode bin192: invalid data")
		return
	}
	if kind != format.KindBin192 {
		err = fmt.Errorf("decode bin192: invalid kind, kind=%v", kind)
		return
	}

	size = n
	start := len(b) - (n + 24)
	end := len(b) - n

	if start < 0 {
		err = fmt.Errorf("decode bin192: invalid data")
		return
	}

	v, err := bin.Parse192(b[start:end])
	if err != nil {
		return
	}

	size += 24
	return v, size, nil
}

func DecodeBin256(b []byte) (_ bin.Bin256, size int, err error) {
	if len(b) == 0 {
		return bin.Bin256{}, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode bin256: invalid data")
		return
	}
	if kind != format.KindBin256 {
		err = fmt.Errorf("decode bin256: invalid kind, kind=%v", kind)
		return
	}

	size = n
	start := len(b) - (n + 32)
	end := len(b) - n

	if start < 0 {
		err = fmt.Errorf("decode bin256: invalid data")
		return
	}

	v, err := bin.Parse256(b[start:end])
	if err != nil {
		return
	}

	size += 32
	return v, size, nil
}
