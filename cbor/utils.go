package cbor

import (
	"fmt"
	"math"
)

const (
	newline = "\n"
	indent  = 2
)

// Debug string representation for a CBOR encoded data item
func Debug(leadingSpace int, obj DataItem) string {

	switch obj.MajorType() {
	case MajorTypePositiveInt,
		MajorTypeNegativeInt,
		MajorTypeByteString,
		MajorTypeTextString,
		MajorTypeSemantic,
		MajorTypePrimitive:
		return whitespace(leadingSpace) + obj.String() + newline

	case MajorTypeArray:
		contents := ""
		for _, item := range obj.(*Array).V {
			contents += Debug(leadingSpace+indent, item)
		}
		return whitespace(leadingSpace) + obj.String() + newline + contents

	case MajorTypeMap:
		contents := ""
		count := 0
		m := obj.(*Map).ValueAsMap()
		for key, value := range m {
			contents += fmt.Sprintf("%s- key: %+v / value: %+v\n", whitespace(leadingSpace+indent), key, value)
			contents += Debug(leadingSpace+indent+indent, value)
			count++
		}
		return whitespace(leadingSpace) + obj.String() + newline + contents

	default:
		return fmt.Sprintf("ERROR, unhandled major type %+v", obj.MajorType())
	}
}

// DebugList return debug string for each item in the list
func DebugList(list []DataItem) string {
	result := ""
	for _, item := range list {
		result += Debug(0, item)
	}
	return result
}

// EncodeList return CBOR representation for each item in the list
func EncodeList(list []DataItem) []byte {
	result := []byte{}
	for _, item := range list {
		result = append(result, item.EncodeCBOR()...)
	}
	return result
}

func whitespace(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += " "
	}
	return result
}

// dataItemPrefix returns majorType (5 bits) | additionalType (3 bits) with the following rules:
//   Length: maxUint32+1 to Maxuint64   Result: [major_3_bits | 27, uint8[0], uint8[1], uint8[2], uint8[3], uint8[4], uint8[5], uint8[6], uint8[7]]
//   Length: maxUint16+1 to maxUint32   Result: [major_3_bits | 26, uint8[0], uint8[1], uint8[2], uint8[3]]
//   Length: maxUint8+1  to maxUint16   Result: [major_3_bits | 25, uint8[0], uint8[1]]
//   Length:          24 to maxUint8    Result: [major_3_bits | 24, uint8]
//   Length:           0 to 23          Result: [major_3_bits | additional_type_5_bits]
func dataItemPrefix(majorType MajorType, length uint64) []byte {

	var result []byte

	switch {
	case length > math.MaxUint32:
		result = []byte{majorType.EncodeCBOR() | additionalType64Bits,
			byte(length >> 56),
			byte(length >> 48),
			byte(length >> 40),
			byte(length >> 32),
			byte(length >> 24),
			byte(length >> 16),
			byte(length >> 8),
			byte(length)}
		break
	case length > math.MaxUint16:
		result = []byte{majorType.EncodeCBOR() | additionalType32Bits,
			byte(length >> 24),
			byte(length >> 16),
			byte(length >> 8),
			byte(length)}
		break
	case length > math.MaxUint8:
		result = []byte{majorType.EncodeCBOR() | additionalType16Bits,
			byte(length >> 8),
			byte(length)}
		break
	case length > uint64(additionalTypeDirectValue23):
		result = []byte{majorType.EncodeCBOR() | additionalType8Bits,
			byte(length)}
		break
	default:
		result = []byte{majorType.EncodeCBOR() | byte(length)}
		break
	}

	return result
}
