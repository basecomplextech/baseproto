// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"errors"
	"fmt"
	"math"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baselibrary/encoding/compactint"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeInt16(b buffer.Buffer, v int16) (int, error) {
	p := [compactint.MaxLen32]byte{}
	n := compactint.PutReverseInt32(p[:], int32(v))
	off := compactint.MaxLen32 - n

	buf := b.Grow(n + 1)
	copy(buf[:n], p[off:])
	buf[n] = byte(format.KindInt16)

	return n + 1, nil
}

func EncodeInt32(b buffer.Buffer, v int32) (int, error) {
	p := [compactint.MaxLen32]byte{}
	n := compactint.PutReverseInt32(p[:], v)
	off := compactint.MaxLen32 - n

	buf := b.Grow(n + 1)
	copy(buf[:n], p[off:])
	buf[n] = byte(format.KindInt32)

	return n + 1, nil
}

func EncodeInt64(b buffer.Buffer, v int64) (int, error) {
	p := [compactint.MaxLen64]byte{}
	n := compactint.PutReverseInt64(p[:], v)
	off := compactint.MaxLen64 - n

	buf := b.Grow(n + 1)
	copy(buf[:n], p[off:])
	buf[n] = byte(format.KindInt64)

	return n + 1, nil
}

// Decode

func DecodeInt16(b []byte) (int16, int, error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		return 0, 0, errors.New("decode int16: invalid data")
	}
	end := len(b) - n

	switch kind {
	case format.KindInt16, format.KindInt32:
		v, m := compactint.ReverseInt32(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode int16: invalid data")
		}

		switch {
		case v < math.MinInt16:
			return 0, 0, errors.New("decode int16: overflow, value too small")
		case v > math.MaxInt16:
			return 0, 0, errors.New("decode int16: overflow, value too large")
		}

		n += m
		return int16(v), n, nil

	case format.KindInt64:
		v, m := compactint.ReverseInt64(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode int16: invalid data")
		}

		switch {
		case v < math.MinInt16:
			return 0, 0, errors.New("decode int16: overflow, value too small")
		case v > math.MaxInt16:
			return 0, 0, errors.New("decode int16: overflow, value too large")
		}

		n += m
		return int16(v), n, nil
	}

	return 0, 0, fmt.Errorf("decode int16: invalid kind, kind=%v", kind)
}

func DecodeInt32(b []byte) (int32, int, error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		return 0, 0, errors.New("decode int32: invalid data")
	}
	end := len(b) - n

	switch kind {
	case format.KindInt16, format.KindInt32:
		v, m := compactint.ReverseInt32(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode int32: invalid data")
		}

		n += m
		return v, n, nil

	case format.KindInt64:
		v, m := compactint.ReverseInt64(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode int32: invalid data")
		}

		switch {
		case v < math.MinInt32:
			return 0, 0, errors.New("decode int32: overflow, value too small")
		case v > math.MaxInt32:
			return 0, 0, errors.New("decode int32: overflow, value too large")
		}

		n += m
		return int32(v), n, nil
	}

	return 0, 0, fmt.Errorf("decode int32: invalid kind, kind=%v", kind)
}

func DecodeInt64(b []byte) (int64, int, error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	kind, n := decodeKind(b)
	if n < 0 {
		return 0, 0, errors.New("decode int64: invalid data")
	}
	end := len(b) - n

	switch kind {
	case format.KindInt16, format.KindInt32:
		v, m := compactint.ReverseInt32(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode int64: invalid data")
		}
		n += m
		return int64(v), n, nil

	case format.KindInt64:
		v, m := compactint.ReverseInt64(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode int64: invalid data")
		}
		n += m
		return int64(v), n, nil
	}

	return 0, 0, fmt.Errorf("decode int64: invalid kind, kind=%v", kind)
}
