// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import (
	"github.com/basecomplextech/baseproto/internal/encode"
)

var (
	EncodeBool = encode.EncodeBool
	EncodeByte = encode.EncodeByte

	EncodeBin64       = encode.EncodeBin64
	EncodeBin128      = encode.EncodeBin128
	EncodeBin128Bytes = encode.EncodeBin128Bytes
	EncodeBin192      = encode.EncodeBin192
	EncodeBin256      = encode.EncodeBin256

	EncodeBytes = encode.EncodeBytes

	EncodeFloat32 = encode.EncodeFloat32
	EncodeFloat64 = encode.EncodeFloat64

	EncodeInt16 = encode.EncodeInt16
	EncodeInt32 = encode.EncodeInt32
	EncodeInt64 = encode.EncodeInt64

	EncodeListTable    = encode.EncodeListTable
	EncodeMessageTable = encode.EncodeMessageTable

	EncodeString = encode.EncodeString
	EncodeStruct = encode.EncodeStruct

	EncodeUint16 = encode.EncodeUint16
	EncodeUint32 = encode.EncodeUint32
	EncodeUint64 = encode.EncodeUint64
)

var (
	DecodeKind     = encode.DecodeKind
	DecodeKindSize = encode.DecodeKindSize

	DecodeBool = encode.DecodeBool
	DecodeByte = encode.DecodeByte

	DecodeBin64  = encode.DecodeBin64
	DecodeBin128 = encode.DecodeBin128
	DecodeBin192 = encode.DecodeBin192
	DecodeBin256 = encode.DecodeBin256

	DecodeBytes = encode.DecodeBytes

	DecodeFloat32 = encode.DecodeFloat32
	DecodeFloat64 = encode.DecodeFloat64

	DecodeInt16 = encode.DecodeInt16
	DecodeInt32 = encode.DecodeInt32
	DecodeInt64 = encode.DecodeInt64

	DecodeListTable    = encode.DecodeListTable
	DecodeMessageTable = encode.DecodeMessageTable

	DecodeString      = encode.DecodeString
	DecodeStringClone = encode.DecodeStringClone

	DecodeStruct = encode.DecodeStruct

	DecodeUint16 = encode.DecodeUint16
	DecodeUint32 = encode.DecodeUint32
	DecodeUint64 = encode.DecodeUint64
)
