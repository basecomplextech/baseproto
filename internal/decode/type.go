// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package decode

import (
	"errors"
	"fmt"

	"github.com/basecomplextech/baselibrary/encoding/compactint"
	"github.com/basecomplextech/baseproto/internal/format"
)

// DecodeType decodes a value type.
func DecodeType(b []byte) (format.Kind, int, error) {
	v, n := decodeType(b)
	if n < 0 {
		return 0, 0, fmt.Errorf("decode type: invalid data")
	}

	size := n
	return format.Kind(v), size, nil
}

// DecodeTypeSize decodes a value type and its total size, returns 0, 0 on error.
func DecodeTypeSize(b []byte) (format.Kind, int, error) {
	if len(b) == 0 {
		return format.KindUndefined, 0, nil
	}

	t, n := decodeType(b)
	if n < 0 {
		return 0, 0, fmt.Errorf("decode type: invalid data")
	}

	end := len(b) - n
	v := b[:end]

	switch t {
	case format.KindTrue, format.KindFalse:
		return t, n, nil

	case format.KindByte:
		if len(v) < 1 {
			return 0, 0, fmt.Errorf("decode byte: invalid data")
		}
		return t, n + 1, nil

	// Int

	case format.KindInt16, format.KindInt32, format.KindInt64:
		m := compactint.ReverseSize(v)
		if m <= 0 {
			return 0, 0, fmt.Errorf("decode int: invalid data")
		}
		return t, n + m, nil

	// Uint

	case format.KindUint16, format.KindUint32, format.KindUint64:
		m := compactint.ReverseSize(v)
		if m <= 0 {
			return 0, 0, fmt.Errorf("decode uint: invalid data")
		}
		return t, n + m, nil

	// Float

	case format.KindFloat32:
		m := 4
		if len(v) < m {
			return 0, 0, fmt.Errorf("decode float32: invalid data")
		}
		return t, n + m, nil

	case format.KindFloat64:
		m := 8
		if len(v) < m {
			return 0, 0, fmt.Errorf("decode float64: invalid data")
		}
		return t, n + m, nil

	// Bin

	case format.KindBin64:
		m := 8
		if len(v) < m {
			return 0, 0, fmt.Errorf("decode bin64: invalid data")
		}
		return t, n + m, nil

	case format.KindBin128:
		m := 16
		if len(v) < m {
			return 0, 0, fmt.Errorf("decode bin128: invalid data")
		}
		return t, n + m, nil

	case format.KindBin192:
		m := 24
		if len(v) < m {
			return 0, 0, fmt.Errorf("decode bin192: invalid data")
		}
		return t, n + m, nil

	case format.KindBin256:
		m := 32
		if len(v) < m {
			return 0, 0, fmt.Errorf("decode bin256: invalid data")
		}
		return t, n + m, nil

	// Bytes/string

	case format.KindBytes:
		dataSize, m := decodeSize(v)
		if m < 0 {
			return 0, 0, errors.New("decode bytes: invalid data size")
		}
		size := n + m + int(dataSize)
		if len(b) < size {
			return 0, 0, errors.New("decode bytes: invalid data")
		}
		return t, size, nil

	case format.KindString:
		dataSize, m := decodeSize(v)
		if m < 0 {
			return 0, 0, errors.New("decode string: invalid data size")
		}
		size := n + m + int(dataSize) + 1 // +1 for null terminator
		if len(b) < size {
			return 0, 0, errors.New("decode string: invalid data")
		}
		return t, size, nil // +1 for null terminator

	// List

	case format.KindList, format.KindListBig:
		size := n

		// Table size
		tableSize, m := decodeSize(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode list: invalid table size")
		}
		end -= m
		size += m + int(tableSize)

		// Data size
		dataSize, m := decodeSize(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode list: invalid data size")
		}
		end -= m
		size += m + int(dataSize)

		if len(b) < size {
			return 0, 0, errors.New("decode list: invalid data")
		}
		return t, size, nil

	// Message

	case format.KindMessage, format.KindMessageBig:
		size := n

		// Table size
		tableSize, m := decodeSize(b[:end])
		if m < 0 {
			return 0, 0, errors.New("decode message: invalid table size")
		}
		end -= m
		size += m + int(tableSize)

		// Data size
		dataSize, m := decodeSize(b[:end])
		if m < 0 {
			return 0, 0, fmt.Errorf("decode message: invalid data size")
		}
		end -= m
		size += m + int(dataSize)

		if len(b) < size {
			return 0, 0, errors.New("decode message: invalid data")
		}
		return t, size, nil

	// Struct

	case format.KindStruct:
		size := n

		// Data size
		dataSize, m := decodeSize(b[:end])
		if n < 0 {
			return 0, 0, errors.New("decode struct: invalid data size")
		}

		size += m + int(dataSize)
		if len(b) < size {
			return 0, 0, errors.New("decode struct: invalid data")
		}
		return t, size, nil
	}

	return 0, 0, fmt.Errorf("decode: invalid kind, kind=%d", t)
}

// internal

func decodeType(b []byte) (format.Kind, int) {
	if len(b) == 0 {
		return format.KindUndefined, 0
	}

	v := b[len(b)-1]
	return format.Kind(v), 1
}

func decodeSize(b []byte) (uint32, int) {
	return compactint.ReverseUint32(b)
}
