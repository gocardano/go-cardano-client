package cbor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoReadBitHasMoreBits(t *testing.T) {

	// Input Binary: 00001111 00001111
	r := NewBitstreamReader([]byte{0x0f, 0x0f})

	// Bit 0-3
	val1, err := r.ReadBitsAsUint8(4)
	assert.Nil(t, err)
	assert.Equal(t, uint8(0x00), val1)
	assert.True(t, r.HasMoreBits())

	// Bit 4-7
	val2, err := r.ReadBitsAsUint16(4)
	assert.Nil(t, err)
	assert.Equal(t, uint16(0x0f), val2)
	assert.True(t, r.HasMoreBits())

	// Bit 8-11
	val3, err := r.ReadBitsAsUint32(4)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0x00), val3)
	assert.True(t, r.HasMoreBits())

	// Bit 12-15
	val4, err := r.ReadBitsAsUint64(4)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0x0f), val4)
	assert.False(t, r.HasMoreBits())

	// Expect error from reading more bits
	_, err = r.ReadBitsAsUint8(1)
	assert.NotNil(t, err)
}

func TestDoReadBitFromOneByte(t *testing.T) {

	// Test 1: Verify reading bit from one byte
	for i := 0; i < 8; i++ {

		// Generate input (0: msb, 7: lsb)
		input := byte(1 << (7 - i))

		r := NewBitstreamReader([]byte{input})
		for j := 0; j < 8; j++ {
			actual, err := r.doReadBit(0, uint64(j))
			assert.Nil(t, err)
			if i == j {
				assert.Equal(t, uint8(1), actual)
			} else {
				assert.Equal(t, uint8(0), actual)
			}
		}
	}
}

func TestDoReadBitWithBitPositionGreaterThan8(t *testing.T) {

	// Input Binary: 11111111 11111111
	r := NewBitstreamReader([]byte{0xff, 0xff})

	for i := 0; i < 16; i++ {
		actual, err := r.doReadBit(0, uint64(i))
		assert.Nil(t, err)
		assert.Equal(t, uint8(1), actual)
	}
}

func TestDoReadBitErrorHandling(t *testing.T) {

	// Input is 1 byte, verify if error is expected on gets
	r := NewBitstreamReader([]byte{0x00})

	var testCases = []struct {
		bytePosition uint32
		bitPosition  uint64
		expectError  bool
	}{
		{bytePosition: 0, bitPosition: 0, expectError: false},
		{bytePosition: 0, bitPosition: 1, expectError: false},
		{bytePosition: 0, bitPosition: 2, expectError: false},
		{bytePosition: 0, bitPosition: 3, expectError: false},
		{bytePosition: 0, bitPosition: 4, expectError: false},
		{bytePosition: 0, bitPosition: 5, expectError: false},
		{bytePosition: 0, bitPosition: 6, expectError: false},
		{bytePosition: 0, bitPosition: 7, expectError: false},
		{bytePosition: 0, bitPosition: 8, expectError: true}, // equivalent to read byte: [1] bit: [0]
		{bytePosition: 0, bitPosition: 9, expectError: true}, // equivalent to read byte: [1] bit: [1]
		{bytePosition: 1, bitPosition: 0, expectError: true},
	}

	for _, testCase := range testCases {
		_, err := r.doReadBit(testCase.bytePosition, testCase.bitPosition)
		assert.Equal(t, testCase.expectError, err != nil)
	}
}

func TestDoReadBitsAsUint(t *testing.T) {

	// Input: 01010101 01010101 01010101 01010101 (32 bits)
	r := NewBitstreamReader([]byte{0x55, 0x55, 0x55, 0x55})

	var testCases = []struct {
		bytePosition uint32
		bitPosition  uint64
		bitCount     uint8
		expect       string
	}{
		{bytePosition: 0, bitPosition: 0, bitCount: 32, expect: "01010101010101010101010101010101"},
		{bytePosition: 0, bitPosition: 0, bitCount: 3, expect: "010"},
		{bytePosition: 0, bitPosition: 1, bitCount: 3, expect: "101"},
		{bytePosition: 3, bitPosition: 0, bitCount: 8, expect: "01010101"},
		{bytePosition: 3, bitPosition: 5, bitCount: 3, expect: "101"},
		{bytePosition: 2, bitPosition: 9, bitCount: 2, expect: "10"},
		{bytePosition: 2, bitPosition: 7, bitCount: 4, expect: "1010"},
	}

	for _, testCase := range testCases {
		actual, err := r.doReadBitsAsUint32(testCase.bytePosition, testCase.bitPosition, testCase.bitCount)
		assert.Nil(t, err)
		assert.Equal(t, testCase.expect, fmt.Sprintf(fmt.Sprintf("%%0%db", testCase.bitCount), actual))
	}

	// Expect error reading greater than input bits
	var errorTestCases = []struct {
		bytePosition uint32
		bitPosition  uint64
		bitCount     uint8
	}{
		{bytePosition: 4, bitPosition: 0, bitCount: 1},
		{bytePosition: 0, bitPosition: 0, bitCount: 33},
		{bytePosition: 3, bitPosition: 0, bitCount: 9},
		{bytePosition: 3, bitPosition: 4, bitCount: 5},
	}

	for _, testCase := range errorTestCases {
		_, err := r.doReadBitsAsUint32(testCase.bytePosition, testCase.bitPosition, testCase.bitCount)
		assert.NotNil(t, err)
	}
}

func TestDoReadBitsAsUintRequestingMoreBitsThanCapacity(t *testing.T) {
	r := NewBitstreamReader([]byte{0x55, 0x55, 0x55, 0x55})
	_, err := r.doReadBitsAsUint8(0, 0, 33)
	assert.NotNil(t, err)
	_, err = r.doReadBitsAsUint16(0, 0, 17)
	assert.NotNil(t, err)
	_, err = r.doReadBitsAsUint32(0, 0, 33)
	assert.NotNil(t, err)
	_, err = r.doReadBitsAsUint64(0, 0, 65)
	assert.NotNil(t, err)
}

func TestLength(t *testing.T) {
	r := NewBitstreamReader([]byte{0x55, 0x55, 0x55, 0x55})
	assert.Equal(t, 4, int(r.lengthInBytes))
	assert.Equal(t, 32, int(r.lengthInBits))
}

func TestReadBytes(t *testing.T) {

	// Input: 01010101 01010101 01010101 01010101 (32 bits)
	r := NewBitstreamReader([]byte{0x55, 0x55, 0x55, 0x55})

	// Test 1: Read one byte
	buf, err := r.ReadBytes(1)
	assert.Equal(t, []byte{0x55}, buf)
	assert.Nil(t, err)

	// Test 2: Read three bytes
	buf, err = r.ReadBytes(3)
	assert.Equal(t, []byte{0x55, 0x55, 0x55}, buf)
	assert.Nil(t, err)

	// Test 3: EOF
	buf, err = r.ReadBytes(1)
	assert.Nil(t, buf)
	assert.NotNil(t, err)
}

func TestReadBytesErrorBitOffset(t *testing.T) {

	r := NewBitstreamReader([]byte{0x55, 0x55})

	// Move the bit
	r.ReadBitsAsUint16(1)

	// Test: Reading bytes outside of byte boundary cause error
	buf, err := r.ReadBytes(1)
	assert.Nil(t, buf)
	assert.NotNil(t, err)
}
