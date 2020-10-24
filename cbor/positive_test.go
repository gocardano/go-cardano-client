package cbor

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositiveUnsignedInt(t *testing.T) {

	testCases := []struct {
		input       []byte
		expectValue interface{}
	}{
		{
			input:       []byte{0x01},
			expectValue: uint8(1),
		},
		{
			input:       []byte{0x17},
			expectValue: uint8(23),
		},
		{
			input:       []byte{0x18, 0x18},
			expectValue: uint8(24),
		},
		{
			input:       []byte{0x18, 0xff},
			expectValue: uint8(255),
		},
		{
			input:       []byte{0x19, 0x01, 0x00},
			expectValue: uint16(256),
		},
		{
			input:       []byte{0x19, 0xff, 0xff},
			expectValue: uint16(65535),
		},
		{
			input:       []byte{0x1A, 0x00, 0x01, 0x00, 0x00},
			expectValue: uint32(65536),
		},
		{
			input:       []byte{0x1A, 0xff, 0xff, 0xff, 0xff},
			expectValue: uint32(4294967295),
		},
		{
			input:       []byte{0x1B, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			expectValue: uint64(4294967296),
		},
		{
			input:       []byte{0x1B, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			expectValue: uint64(18446744073709551615),
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, MajorTypePositiveInt, c[0].MajorType())
		assert.Equal(t, testCase.expectValue, c[0].Value())
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}

func TestNewPositiveUnsignedInt(t *testing.T) {

	// PositiveInteger8
	value8 := NewPositiveInteger(1)
	assert.Equal(t, uint8(1), value8.(*PositiveInteger8).ValueAsUint8())
	assert.Equal(t, uint16(1), value8.(*PositiveInteger8).ValueAsUint16())
	assert.Equal(t, uint32(1), value8.(*PositiveInteger8).ValueAsUint32())
	assert.Equal(t, uint64(1), value8.(*PositiveInteger8).ValueAsUint64())

	// PositiveInteger16
	value16 := NewPositiveInteger(math.MaxUint8 + 1)
	assert.Equal(t, uint16(math.MaxUint8+1), value16.(*PositiveInteger16).ValueAsUint16())
	assert.Equal(t, uint32(math.MaxUint8+1), value16.(*PositiveInteger16).ValueAsUint32())
	assert.Equal(t, uint64(math.MaxUint8+1), value16.(*PositiveInteger16).ValueAsUint64())

	// PositiveInteger32
	value32 := NewPositiveInteger(math.MaxUint16 + 1)
	assert.Equal(t, uint32(math.MaxUint16+1), value32.(*PositiveInteger32).ValueAsUint32())
	assert.Equal(t, uint64(math.MaxUint16+1), value32.(*PositiveInteger32).ValueAsUint64())

	// PositiveInteger64
	value64 := NewPositiveInteger(math.MaxUint32 + 1)
	assert.Equal(t, uint64(math.MaxUint32+1), value64.(*PositiveInteger64).ValueAsUint64())
}
