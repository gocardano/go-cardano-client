package cbor

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimitivesConstants(t *testing.T) {

	var testCases = []struct {
		dataItem             DataItem
		expectMajorType      MajorType
		expectAdditionalType uint8
		expectedValue        interface{}
	}{
		{
			dataItem:             NewPrimitiveFalse(),
			expectMajorType:      MajorTypePrimitive,
			expectAdditionalType: primitiveFalse,
			expectedValue:        false,
		},
		{
			dataItem:             NewPrimitiveTrue(),
			expectMajorType:      MajorTypePrimitive,
			expectAdditionalType: primitiveTrue,
			expectedValue:        true,
		},
		{
			dataItem:             NewPrimitiveNull(),
			expectMajorType:      MajorTypePrimitive,
			expectAdditionalType: primitiveNull,
			expectedValue:        nil,
		},
		{
			dataItem:             NewPrimitiveUndefined(),
			expectMajorType:      MajorTypePrimitive,
			expectAdditionalType: primitiveUndefined,
			expectedValue:        nil,
		},
		{
			dataItem:             NewPrimitiveBreakStopCode(),
			expectMajorType:      MajorTypePrimitive,
			expectAdditionalType: primitiveBreakStopCode,
			expectedValue:        nil,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectMajorType, testCase.dataItem.MajorType())
		assert.Equal(t, testCase.expectAdditionalType, testCase.dataItem.AdditionalType())
		assert.Equal(t, testCase.expectedValue, testCase.dataItem.Value())
	}
}

func TestPrimitivesParsing(t *testing.T) {

	testCases := []struct {
		input                []byte
		expectAdditionalType uint8
		expectValue          interface{}
	}{
		{
			input:                []byte{0xf4},
			expectAdditionalType: primitiveFalse,
			expectValue:          false,
		},
		{
			input:                []byte{0xf5},
			expectAdditionalType: primitiveTrue,
			expectValue:          true,
		},
		{
			input:                []byte{0xf6},
			expectAdditionalType: primitiveNull,
			expectValue:          nil,
		},
		{
			input:                []byte{0xf7},
			expectAdditionalType: primitiveUndefined,
			expectValue:          nil,
		},
		{
			input:                []byte{0xf8, 0x20},
			expectAdditionalType: primitiveSimpleValue,
			expectValue:          uint8(32),
		},
		{
			input:                []byte{0xf8, 0xff},
			expectAdditionalType: primitiveSimpleValue,
			expectValue:          uint8(255),
		},
		{
			input:                []byte{0xf9, 0x3c, 0x00},
			expectAdditionalType: primitiveHalfPrecisionFloat,
			expectValue:          float32(1),
		},
		{
			input:                []byte{0xf9, 0x56, 0x40},
			expectAdditionalType: primitiveHalfPrecisionFloat,
			expectValue:          float32(100),
		},
		{
			input:                []byte{0xf9, 0x70, 0x57},
			expectAdditionalType: primitiveHalfPrecisionFloat,
			expectValue:          float32(8888),
		},
		{
			input:                []byte{0xfa, 0x40, 0xa0, 0x00, 0x00},
			expectAdditionalType: primitiveSinglePrecisionFloat,
			expectValue:          float32(5),
		},
		{
			input:                []byte{0xfb, 0x40, 0xC1, 0x5C, 0x70, 0xA3, 0xD7, 0x0A, 0x3D},
			expectAdditionalType: primitiveDoublePrecisionFloat,
			expectValue:          float64(8888.88),
		},
		{
			input:                []byte{0xfb, 0x41, 0x9d, 0x6f, 0x34, 0x57, 0xf3, 0x5b, 0xa8},
			expectAdditionalType: primitiveDoublePrecisionFloat,
			expectValue:          float64(123456789.987654321),
		},
		{
			input:                []byte{0xff},
			expectAdditionalType: primitiveBreakStopCode,
			expectValue:          nil,
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, MajorTypePrimitive, c[0].MajorType())
		assert.Equal(t, testCase.expectAdditionalType, c[0].AdditionalType())
		assert.Equal(t, testCase.expectValue, c[0].Value())
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}

func TestNewPrimitiveSimpleValue(t *testing.T) {
	assert.Nil(t, NewPrimitiveSimpleValue(primitiveSimpleValueMin-1))
	assert.Equal(t, primitiveSimpleValueMin, NewPrimitiveSimpleValue(primitiveSimpleValueMin).Value())
	assert.Equal(t, primitiveSimpleValueMax, NewPrimitiveSimpleValue(primitiveSimpleValueMax).Value())
}

func TestNewPrimitiveFloatRange(t *testing.T) {
	assert.NotNil(t, NewPrimitiveHalfPrecisionFloat(1))

	assert.NotNil(t, NewPrimitiveSinglePrecisionFloat(math.MaxFloat32))
	assert.NotNil(t, NewPrimitiveSinglePrecisionFloat(math.SmallestNonzeroFloat32))

	assert.NotNil(t, NewPrimitiveDoublePrecisionFloat(math.MaxFloat64))
	assert.NotNil(t, NewPrimitiveDoublePrecisionFloat(math.SmallestNonzeroFloat64))
}
