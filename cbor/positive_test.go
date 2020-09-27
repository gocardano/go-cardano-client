package cbor

import (
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
