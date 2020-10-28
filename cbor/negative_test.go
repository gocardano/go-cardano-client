package cbor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegativeIntErrorCaseZeroOrPositive(t *testing.T) {

	array := []int64{0, 1}

	for _, i := range array {
		assert.Nil(t, NewNegativeInteger8(i))
		assert.Nil(t, NewNegativeInteger16(i))
		assert.Nil(t, NewNegativeInteger32(i))
		assert.Nil(t, NewNegativeInteger64(i))
	}
}

func TestNegativeInt(t *testing.T) {

	testCases := []struct {
		input  []byte
		expect interface{}
	}{
		// Test the int8 boundaries [-128 to 127]:
		{
			input:  []byte{0x20},
			expect: int8(-1),
		},
		{
			input:  []byte{0x38, 0x7f},
			expect: int8(-128),
		},

		// Test the int16 boundaries [-32,768 to 32,767]:
		{
			input:  []byte{0x38, 0x80},
			expect: int16(-129),
		},
		{
			input:  []byte{0x39, 0x7f, 0xff},
			expect: int16(-32768),
		},

		// Test the int32 boundaries [-2,147,483,648 to 2,147,483,647]:
		{
			input:  []byte{0x39, 0x80, 0x00},
			expect: int32(-32769),
		},
		{
			input:  []byte{0x3A, 0x7f, 0xff, 0xff, 0xff},
			expect: int32(-2147483648),
		},

		// Test the int64 boundaries [-9,223,372,036,854,775,808 to 9,223,372,036,854,775,807]
		{
			input:  []byte{0x3a, 0x80, 0x00, 0x00, 0x00},
			expect: int64(-2147483649),
		},

		// Test the uint8 (encoded value) boundaries [0-255]:
		{
			input:  []byte{0x38, 0xff},
			expect: int16(-256),
		},
		{
			input:  []byte{0x39, 0x01, 0x00},
			expect: int16(-257),
		},

		// Test the uint16 (encoded value) boundaries [0-65535]:
		{
			input:  []byte{0x39, 0xff, 0xff},
			expect: int32(-65536),
		},
		{
			input:  []byte{0x3a, 0x00, 0x01, 0x00, 0x00},
			expect: int32(-65537),
		},

		// Test the uint32 (encoded value) boundaries [0-4294967295]:
		{
			input:  []byte{0x3a, 0xff, 0xff, 0xff, 0xff},
			expect: int64(-4294967296),
		},
		{
			input:  []byte{0x3b, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			expect: int64(-4294967297),
		},
	}

	for _, testCase := range testCases {
		t.Log("Testing with input: ", fmt.Sprintf("%x", testCase.input))
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, MajorTypeNegativeInt, c[0].MajorType(), "incorrect major type")
		assert.Equal(t, testCase.expect, c[0].Value(), "incorrect value")
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}
