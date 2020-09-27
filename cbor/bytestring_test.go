package cbor

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteString(t *testing.T) {

	testCases := []struct {
		scenario            string
		input               []byte
		expectValue         []byte
		expectPayloadLength int
	}{
		{
			scenario:            "additionalType == 1 (payloadLength)",
			input:               []byte{0x41, 0x61},
			expectValue:         []byte{0x61},
			expectPayloadLength: 1,
		},
		{
			scenario:            "additionalType == 2 (payloadLength)",
			input:               []byte{0x42, 0x61, 0x62},
			expectValue:         []byte{0x61, 0x62},
			expectPayloadLength: 2,
		},
		{
			scenario:            "additionalType == 23 (payloadLength)",
			input:               []byte{0x57, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63},
			expectValue:         []byte{0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63},
			expectPayloadLength: 23,
		},
		{
			scenario:            "additionalType == 24 (payloadLength), next uint8 is payload length (24)",
			input:               []byte{0x58, 0x18, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64},
			expectValue:         []byte{0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64},
			expectPayloadLength: 24,
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(c), "incorrect number of CBOR items parsed")
		assert.Equal(t, MajorTypeByteString, c[0].MajorType(), "incorrect major type")
		assert.Equal(t, testCase.expectValue, c[0].Value(), "incorrect value")
		assert.Equal(t, testCase.expectPayloadLength, len(c[0].Value().([]byte)), "incorrect payload length")
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}

func TestByteStringDecodingBoundaries(t *testing.T) {

	var testCases = []struct {
		additionalType      uint8
		additionalTypeValue uint64
	}{
		{additionalType: additionalType8Bits, additionalTypeValue: 24},
		{additionalType: additionalType8Bits, additionalTypeValue: math.MaxUint8},
		{additionalType: additionalType16Bits, additionalTypeValue: math.MaxUint8 + 1},
		{additionalType: additionalType32Bits, additionalTypeValue: math.MaxUint16},
		{additionalType: additionalType64Bits, additionalTypeValue: math.MaxUint16 + 1},
		// {additionalType: additionalType32Bits, additionalTypeValue: math.MaxUint32},
		// {additionalType: additionalType64Bits, additionalTypeValue: math.MaxUint32 + 1},
		// {additionalType: additionalType64Bits, additionalTypeValue: math.MaxUint64},
	}

	for _, testCase := range testCases {

		input := dataItemPrefix(MajorTypeByteString, testCase.additionalTypeValue)
		input = append(input, make([]byte, testCase.additionalTypeValue)...)

		c, err := Decode(input)
		assert.Nil(t, err)
		assert.Equal(t, input, c[0].EncodeCBOR())
		assert.Equal(t, testCase.additionalTypeValue, uint64(len(c[0].Value().([]byte))))
	}
}

func TestByteStringEncodeAsChunks(t *testing.T) {

	input := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	b := NewByteString(input)

	// Test: Chunk by 2
	// 	5F         # bytes(*)
	//    42      # bytes(2)
	//       0102 # "\x01\x02"
	//    42      # bytes(2)
	//       0304 # "\x03\x04"
	//    41      # bytes(1)
	//       05   # "\x05"
	//    FF      # primitive(*)
	assert.Equal(t, []byte{0x5f, 0x42, 0x01, 0x02, 0x42, 0x03, 0x04, 0x41, 0x05, 0xff}, b.encodeAsCBORChunks(2))

	// Test: Chunk by 3
	// 	5F          # bytes(*)
	//    43        # bytes(2)
	//       010203 # "\x01\x02\x03"
	//    42        # bytes(2)
	//       0405   # "\x04\x05"
	//    FF        # primitive(*)
	assert.Equal(t, []byte{0x5f, 0x43, 0x01, 0x02, 0x03, 0x42, 0x04, 0x05, 0xff}, b.encodeAsCBORChunks(3))

	// Test: With no chunking
	// 	45           # bytes(*)
	//    0102030405 # "\x01\x02\x03\x04\x05"
	assert.Equal(t, []byte{0x45, 0x01, 0x02, 0x03, 0x04, 0x05}, b.EncodeCBOR())
}
