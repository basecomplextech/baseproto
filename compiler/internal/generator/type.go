// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
)

// typeName returns a type name.
func typeName(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindAny:
		return "baseproto.Value"

	case model.KindBool:
		return "bool"
	case model.KindByte:
		return "byte"

	case model.KindInt16:
		return "int16"
	case model.KindInt32:
		return "int32"
	case model.KindInt64:
		return "int64"

	case model.KindUint16:
		return "uint16"
	case model.KindUint32:
		return "uint32"
	case model.KindUint64:
		return "uint64"

	case model.KindFloat32:
		return "float32"
	case model.KindFloat64:
		return "float64"

	case model.KindBin64:
		return "bin.Bin64"
	case model.KindBin128:
		return "bin.Bin128"
	case model.KindBin192:
		return "bin.Bin192"
	case model.KindBin256:
		return "bin.Bin256"

	case model.KindBytes:
		return "[]byte"
	case model.KindString:
		return "string"
	case model.KindAnyMessage:
		return "baseproto.Message"

	case model.KindList:
		elem := typeName(typ.Element)
		if typ.Element.Kind == model.KindMessage {
			return fmt.Sprintf("baseproto.MessageList[%v]", elem)
		}
		return fmt.Sprintf("baseproto.ValueList[%v]", elem)

	case model.KindEnum,
		model.KindMessage,
		model.KindStruct:
		if typ.Import != nil {
			return fmt.Sprintf("%v.%v", typ.ImportName, typ.Name)
		}
		return typ.Name

	case model.KindService:
		if typ.Import != nil {
			return fmt.Sprintf("%v.%v", typ.ImportName, typ.Name)
		}
		return typ.Name
	}

	panic(fmt.Sprintf("unsupported type kind %v", typ.Kind))
}

func typeRefName(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindBytes:
		return "baseproto.Bytes"
	case model.KindString:
		return "baseproto.String"

	case model.KindList:
		elem := typeRefName(typ.Element)
		if typ.Element.Kind == model.KindMessage {
			return fmt.Sprintf("baseproto.MessageList[%v]", elem)
		}
		return fmt.Sprintf("baseproto.ValueList[%v]", elem)
	}

	return typeName(typ)
}

func inTypeName(typ *model.Type) string {
	kind := typ.Kind
	switch kind {
	case model.KindBytes:
		return "[]byte"
	case model.KindString:
		return "string"
	}
	return typeName(typ)
}

func typeNewFunc(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindList:
		elem := typeName(typ.Element)
		return "baseproto.List[]" + elem

	case model.KindEnum,
		model.KindMessage,
		model.KindStruct:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Open%v", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Open%v", typ.Name)
	}
	return ""
}

func typeMakeMessageFunc(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindMessage:
		if typ.Import != nil {
			return fmt.Sprintf("%v.New%v", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("New%v", typ.Name)
	}

	panic(fmt.Sprintf("unsupported type kind %v", typ.Kind))
}

func typeParseFunc(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindList:
		elem := typeName(typ.Element)
		return "baseproto.List[]" + elem

	case model.KindEnum,
		model.KindStruct:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Decode%v", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Decode%v", typ.Name)

	case model.KindMessage:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Parse%v", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Parse%v", typ.Name)

	default:
		return typeDecodeFunc(typ)
	}
}

func typeDecodeFunc(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindAny:
		return "baseproto.ParseValue"

	case model.KindBool:
		return "baseproto.DecodeBool"
	case model.KindByte:
		return "baseproto.DecodeByte"

	case model.KindInt16:
		return "baseproto.DecodeInt16"
	case model.KindInt32:
		return "baseproto.DecodeInt32"
	case model.KindInt64:
		return "baseproto.DecodeInt64"

	case model.KindUint16:
		return "baseproto.DecodeUint16"
	case model.KindUint32:
		return "baseproto.DecodeUint32"
	case model.KindUint64:
		return "baseproto.DecodeUint64"

	case model.KindFloat32:
		return "baseproto.DecodeFloat32"
	case model.KindFloat64:
		return "baseproto.DecodeFloat64"

	case model.KindBin64:
		return "baseproto.DecodeBin64"
	case model.KindBin128:
		return "baseproto.DecodeBin128"
	case model.KindBin192:
		return "baseproto.DecodeBin192"
	case model.KindBin256:
		return "baseproto.DecodeBin256"

	case model.KindBytes:
		return "baseproto.DecodeBytes"
	case model.KindString:
		return "baseproto.DecodeString"
	case model.KindAnyMessage:
		return "baseproto.ParseMessage"

	case model.KindList:
		elem := typ.Element
		name := typeName(typ.Element)
		if elem.Kind == model.KindMessage {
			return fmt.Sprintf("baseproto.OpenMessageListErr[%v]", name)
		}
		return fmt.Sprintf("baseproto.OpenValueListErr[%v]", name)

	case model.KindEnum,
		model.KindStruct:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Decode%v", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Decode%v", typ.Name)

	case model.KindMessage:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Open%vErr", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Open%vErr", typ.Name)
	}

	return ""
}

func typeDecodeRefFunc(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindBytes:
		return "baseproto.DecodeBytes"
	case model.KindString:
		return "baseproto.DecodeString"

	case model.KindList:
		elem := typeRefName(typ.Element)
		return fmt.Sprintf("baseproto.ParseTypedList[%v]", elem)
	}

	return typeDecodeFunc(typ)
}

func typeWriteFunc(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindAny:
		return "baseproto.WriteValue"

	case model.KindBool:
		return "baseproto.EncodeBool"
	case model.KindByte:
		return "baseproto.EncodeByte"

	case model.KindInt16:
		return "baseproto.EncodeInt16"
	case model.KindInt32:
		return "baseproto.EncodeInt32"
	case model.KindInt64:
		return "baseproto.EncodeInt64"

	case model.KindUint16:
		return "baseproto.EncodeUint16"
	case model.KindUint32:
		return "baseproto.EncodeUint32"
	case model.KindUint64:
		return "baseproto.EncodeUint64"

	case model.KindBin64:
		return "baseproto.EncodeBin64"
	case model.KindBin128:
		return "baseproto.EncodeBin128"
	case model.KindBin192:
		return "baseproto.EncodeBin192"
	case model.KindBin256:
		return "baseproto.EncodeBin256"

	case model.KindFloat32:
		return "baseproto.EncodeFloat32"
	case model.KindFloat64:
		return "baseproto.EncodeFloat64"

	case model.KindBytes:
		return "baseproto.EncodeBytes"
	case model.KindString:
		return "baseproto.EncodeString"
	case model.KindAnyMessage:
		return "baseproto.WriteMessage"

	case model.KindEnum:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Encode%vTo", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Encode%vTo", typ.Name)

	case model.KindList:
		elem := typ.Element
		if elem.Kind == model.KindMessage {
			return fmt.Sprintf("baseproto.NewMessageListWriter")
		}
		return fmt.Sprintf("baseproto.NewValueListWriter")

	case model.KindMessage:
		if typ.Import != nil {
			return fmt.Sprintf("%v.New%vWriterTo", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("New%vWriterTo", typ.Name)

	case model.KindStruct:
		if typ.Import != nil {
			return fmt.Sprintf("%v.Encode%vTo", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("Encode%vTo", typ.Name)
	}

	return ""
}

func typeWriter(typ *model.Type) string {
	kind := typ.Kind

	switch kind {
	case model.KindList:
		elem := typ.Element
		if elem.Kind == model.KindMessage {
			encoder := typeWriter(elem)
			return fmt.Sprintf("baseproto.MessageListWriter[%v]", encoder)
		}

		elemName := inTypeName(elem)
		return fmt.Sprintf("baseproto.ValueListWriter[%v]", elemName)

	case model.KindMessage:
		if typ.Import != nil {
			return fmt.Sprintf("%v.%vWriter", typ.ImportName, typ.Name)
		}
		return fmt.Sprintf("%vWriter", typ.Name)
	}

	return ""
}
