package cbor

// MajorType represents the CBOR major type (first 3 bits)
type MajorType int

const (
	// MajorTypePositiveInt is the type for unsigned positive integer
	MajorTypePositiveInt MajorType = 0

	// MajorTypeNegativeInt is the type for negative integer
	MajorTypeNegativeInt MajorType = 1

	// MajorTypeByteString is the type for a byte array
	MajorTypeByteString MajorType = 2

	// MajorTypeTextString is the type for a text string
	MajorTypeTextString MajorType = 3

	// MajorTypeArray is the type for an array
	MajorTypeArray MajorType = 4

	// MajorTypeMap is the type for a map
	MajorTypeMap MajorType = 5

	// MajorTypeSemantic is the type for a semantic object
	MajorTypeSemantic MajorType = 6

	// MajorTypePrimitive is the type for a primitive
	MajorTypePrimitive MajorType = 7
)

const (
	additionalTypeDirectValue23 uint8 = 23
	additionalType8Bits         uint8 = 24
	additionalType16Bits        uint8 = 25
	additionalType32Bits        uint8 = 26
	additionalType64Bits        uint8 = 27
	additionalTypeIndefinite    uint8 = 31
)

const (
	primitiveFalse                uint8 = 20
	primitiveTrue                 uint8 = 21
	primitiveNull                 uint8 = 22
	primitiveUndefined            uint8 = 23
	primitiveSimpleValue          uint8 = 24
	primitiveHalfPrecisionFloat   uint8 = 25
	primitiveSinglePrecisionFloat uint8 = 26
	primitiveDoublePrecisionFloat uint8 = 27
	primitiveBreakStopCode        uint8 = 31

	indefiniteBreakCode = 0xff
)

const (
	semanticDateTimeString                uint64 = 0
	semanticDateTimeEpoch                 uint64 = 1
	semanticPositiveBignum                uint64 = 2
	semanticNegativeBignum                uint64 = 3
	semanticDecimalFraction               uint64 = 4
	semanticBigFloat                      uint64 = 5
	semanticExpectedConversionToBase64URL uint64 = 21
	semanticExpectedConversionToBase64    uint64 = 22
	semanticExpectedConversionToBase16    uint64 = 23
	semanticEncodedCBORDataItems          uint64 = 24
	semanticURI                           uint64 = 32
	semanticBase64URL                     uint64 = 33
	semanticBase64                        uint64 = 34
	semanticRegularExpression             uint64 = 35
	semanticMimeMessage                   uint64 = 36
	semanticSelfDescribeCBOR              uint64 = 55799
)

var majorTypes = map[MajorType]string{
	MajorTypePositiveInt: "PositiveUnsignedInt",
	MajorTypeNegativeInt: "NegativeInteger",
	MajorTypeByteString:  "ByteString",
	MajorTypeTextString:  "TextString",
	MajorTypeArray:       "Array",
	MajorTypeMap:         "MapPairsDataItems",
	MajorTypeSemantic:    "SemanticTag",
	MajorTypePrimitive:   "Primitive",
}

// EncodeCBOR returns this as CBOR (first 3 bits of the byte, big endian)
func (m MajorType) EncodeCBOR() byte {
	return uint8(m) << 5
}

// String description of this major type
func (m MajorType) String() string {
	return majorTypes[m]
}
