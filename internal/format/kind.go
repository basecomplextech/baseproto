// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package format

import (
	"fmt"
	"math"
	"strconv"
)

const (
	MaxSize = math.MaxInt32
)

// Kind specifies a value kind.
type Kind byte

const (
	KindUndefined Kind = 0

	KindTrue  Kind = 01
	KindFalse Kind = 02
	KindByte  Kind = 03

	KindInt16 Kind = 10
	KindInt32 Kind = 11
	KindInt64 Kind = 12

	KindUint16 Kind = 20
	KindUint32 Kind = 21
	KindUint64 Kind = 22

	KindFloat32 Kind = 30
	KindFloat64 Kind = 31

	KindBin64  Kind = 40
	KindBin128 Kind = 41
	KindBin192 Kind = 42
	KindBin256 Kind = 43

	KindBytes  Kind = 50
	KindString Kind = 60

	KindList    Kind = 70
	KindListBig Kind = 71

	KindMessage    Kind = 80
	KindMessageBig Kind = 81

	KindStruct Kind = 90
)

func (k Kind) Check() error {
	switch k {
	case
		KindTrue,
		KindFalse,
		KindByte,

		KindInt16,
		KindInt32,
		KindInt64,

		KindUint16,
		KindUint32,
		KindUint64,

		KindFloat32,
		KindFloat64,

		KindBin64,
		KindBin128,
		KindBin192,
		KindBin256,

		KindBytes,
		KindString,

		KindList,
		KindListBig,

		KindMessage,
		KindMessageBig,

		KindStruct:
		return nil
	}

	return fmt.Errorf("unsupported type %d", k)
}

func (k Kind) String() string {
	switch k {
	case KindTrue:
		return "true"
	case KindFalse:
		return "false"
	case KindByte:
		return "int8"

	case KindInt16:
		return "int16"
	case KindInt32:
		return "int32"
	case KindInt64:
		return "int64"

	case KindUint16:
		return "uint16"
	case KindUint32:
		return "uint32"
	case KindUint64:
		return "uint64"

	case KindFloat32:
		return "float32"
	case KindFloat64:
		return "float64"

	case KindBin64:
		return "bin64"
	case KindBin128:
		return "bin128"
	case KindBin192:
		return "bin192"
	case KindBin256:
		return "bin256"

	case KindBytes:
		return "bytes"
	case KindString:
		return "string"

	case KindList:
		return "list"
	case KindListBig:
		return "big_list"

	case KindMessage:
		return "message"
	case KindMessageBig:
		return "big_message"

	case KindStruct:
		return "struct"
	}

	return strconv.Itoa(int(k))
}
