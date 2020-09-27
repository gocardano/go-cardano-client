package cbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextString(t *testing.T) {

	testCases := []struct {
		scenario    string
		input       []byte
		expectValue string
	}{
		{
			scenario:    "additionalType == 1 (payloadLength)",
			input:       []byte{0x61, 0x61},
			expectValue: "a",
		},
		{
			scenario:    "additionalType == 2 (payloadLength)",
			input:       []byte{0x62, 0x61, 0x62},
			expectValue: "ab",
		},
		{
			scenario:    "additionalType == 23 (payloadLength)",
			input:       []byte{0x77, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63, 0x64, 0x61, 0x62, 0x63},
			expectValue: "abcdabcdabcdabcdabcdabc",
		},
		{
			scenario:    "additionalType == 24 (payloadLength), next uint8 is payload length",
			input:       []byte{0x78, 0x01, 0x61},
			expectValue: "a",
		},
		{
			scenario:    "additionalType == 25 (payloadLength), next uint16 is payload length",
			input:       []byte{0x79, 0x00, 0x01, 0x61},
			expectValue: "a",
		},
		{
			scenario:    "additionalType == 26 (payloadLength), next uint32 is payload length",
			input:       []byte{0x7A, 0x00, 0x00, 0x00, 0x01, 0x61},
			expectValue: "a",
		},
		{
			scenario:    "additionalType == 27 (payloadLength), next uint64 is payload length",
			input:       []byte{0x7B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x61},
			expectValue: "a",
		},
		{
			scenario:    "additionalType == 31 (indefinite length); pattern: [ignore 3 bits, 5 bits length, ... until 0xFF break code",
			input:       []byte{0x7f, 0x44, 0x61, 0x61, 0x61, 0x61, 0x43, 0x61, 0x61, 0x61, 0xff},
			expectValue: "aaaaaaa",
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(c), "incorrect number of CBOR items parsed")
		assert.Equal(t, MajorTypeTextString, c[0].MajorType(), "incorrect major type")
		assert.Equal(t, testCase.expectValue, c[0].Value(), "incorrect value")
	}
}
