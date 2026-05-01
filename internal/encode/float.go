// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeFloat32(b buffer.Buffer, v float32) (int, error) {
	p := b.Grow(5)
	binary.BigEndian.PutUint32(p, math.Float32bits(v))
	p[4] = byte(format.KindFloat32)
	return 5, nil
}

func EncodeFloat64(b buffer.Buffer, v float64) (int, error) {
	p := b.Grow(9)
	binary.BigEndian.PutUint64(p, math.Float64bits(v))
	p[8] = byte(format.KindFloat64)
	return 9, nil
}

// Decode

func DecodeFloat32(b []byte) (float32, int, error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	v, n := decodeFloat64(b)
	switch {
	case n < 0:
		return 0, 0, errors.New("decode float32: invalid data")
	case v < -math.MaxFloat32:
		return 0, 0, errors.New("decode float32: overflow, value too small")
	case v > math.MaxFloat32:
		return 0, 0, errors.New("decode float32: overflow, value too large")
	}

	size := n
	return float32(v), size, nil
}

func DecodeFloat64(b []byte) (float64, int, error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	v, n := decodeFloat64(b)
	if n < 0 {
		return 0, n, errors.New("decode float64: invalid data")
	}

	size := n
	return v, size, nil
}

// private

// decodeFloat64 reads and returns any float as float64 and the number of decode bytes n,
// or -n on error.
func decodeFloat64(b []byte) (float64, int) {
	t, n := decodeKind(b)
	if n < 0 {
		return 0, n
	}

	switch t {
	case format.KindFloat32:
		start := len(b) - 5
		if start < 0 {
			return 0, -1
		}

		v := binary.BigEndian.Uint32(b[start:])
		f := math.Float32frombits(v)
		return float64(f), 5

	case format.KindFloat64:
		start := len(b) - 9
		if start < 0 {
			return 0, -1
		}

		v := binary.BigEndian.Uint64(b[start:])
		f := math.Float64frombits(v)
		return f, 9
	}

	return 0, -1
}
