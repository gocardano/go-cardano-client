package cbor

import (
	"encoding/base64"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSemanticConstants(t *testing.T) {

	var testCases = []struct {
		dataItem             DataItem
		expectMajorType      MajorType
		expectAdditionalType uint64
	}{
		{
			dataItem:             NewDateTimeString(rfc3339),
			expectMajorType:      MajorTypeSemantic,
			expectAdditionalType: semanticDateTimeString,
		},
		{
			dataItem:             NewDateTimeEpoch(time.Now().Unix()),
			expectMajorType:      MajorTypeSemantic,
			expectAdditionalType: semanticDateTimeEpoch,
		},
		{
			dataItem:             NewBase64URL("MQ=="),
			expectMajorType:      MajorTypeSemantic,
			expectAdditionalType: semanticBase64URL,
		},
		{
			dataItem:             NewBase64String("aHR0cDovL2EuY29t"), // http://a.com
			expectMajorType:      MajorTypeSemantic,
			expectAdditionalType: semanticBase64,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectMajorType, testCase.dataItem.MajorType())
		assert.Equal(t, testCase.expectAdditionalType, testCase.dataItem.AdditionalTypeValue())
	}
}

func TestSemanticBignum(t *testing.T) {

	var testCases = []struct {
		input                     []byte
		expectAdditionalTypeValue uint64
		expectValue               string
	}{
		{
			input:                     []byte{0xc2, 0x49, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectAdditionalTypeValue: semanticPositiveBignum,
			expectValue:               "18446744073709551616",
		},
		{
			input:                     []byte{0xc3, 0x49, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectAdditionalTypeValue: semanticNegativeBignum,
			expectValue:               "-18446744073709551617",
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.NotNil(t, c[0])
		assert.Equal(t, MajorTypeSemantic, c[0].MajorType())
		assert.Equal(t, testCase.expectAdditionalTypeValue, c[0].AdditionalTypeValue())
		assert.Equal(t, testCase.expectValue, c[0].Value().(*big.Int).String())
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}

	// Test Edge Case: positive or negative
	assert.Nil(t, NewPositiveBignumber(big.NewInt(-1)))
	assert.Nil(t, NewNegativeBignumber(big.NewInt(1)))
}

func TestSemanticParsingDateTime(t *testing.T) {

	testCases := []struct {
		input                []byte
		expectAdditionalType uint8
		expectValue          string
	}{
		{
			// Input: 0("2006-01-02T15:04:05Z")
			// C0                                      # tag(0)
			//    74                                   # text(20)
			//       323030362D30312D30325431353A30343A30355A # "2006-01-02T15:04:05Z"
			input:       []byte{0xc0, 0x74, 0x32, 0x30, 0x30, 0x36, 0x2d, 0x30, 0x31, 0x2d, 0x30, 0x32, 0x54, 0x31, 0x35, 0x3a, 0x30, 0x34, 0x3a, 0x30, 0x35, 0x5a},
			expectValue: rfc3339,
		},
		{
			// Input: 1(200) int8
			// C1             # tag(1)
			// 18 C8          # unsigned(200)
			input:       []byte{0xc1, 0x18, 0xc8},
			expectValue: "1970-01-01T00:03:20Z",
		},
		{
			// Input: 1(30000) int16
			// C1             # tag(1)
			// 19 75 30       # unsigned(30000)
			input:       []byte{0xc1, 0x19, 0x75, 0x30},
			expectValue: "1970-01-01T08:20:00Z",
		},
		{
			// Input: 1(1136214245) int32
			// C1             # tag(1)
			// 1A 43B940E5    # unsigned(1136214245)
			input:       []byte{0xc1, 0x1a, 0x43, 0xb9, 0x40, 0xe5},
			expectValue: rfc3339,
		},
		{
			// Input: 1(4294967296) int64
			// C1                      # tag(1)
			// 1B 00000001 00000000    # unsigned(4294967296)
			input:       []byte{0xc1, 0x1b, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			expectValue: "2106-02-07T06:28:16Z",
		},
		{
			// Input: 1(-1) negative int8
			// C1  			 # tag(1)
			// 21  			 # negative(-1)
			input:       []byte{0xc1, 0x21},
			expectValue: "1969-12-31T23:59:58Z",
		},
		{
			// Input: 1(-2147483648) negative int32
			// C1  	         # tag(1)
			// 3A 7fffffff	 # negative(-2147483648)
			input:       []byte{0xc1, 0x3a, 0x7f, 0xff, 0xff, 0xff},
			expectValue: "1901-12-13T20:45:52Z",
		},
		{
			// Input: 1(-4294967297) negative int64
			// C1  	         # tag(1)
			// 3B 80000000	 # negative(-4294967297)
			input:       []byte{0xc1, 0x3b, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			expectValue: "1833-11-24T17:31:43Z",
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, MajorTypeSemantic, c[0].MajorType(), "incorrect major type")
		assert.Equal(t, testCase.expectValue, c[0].Value().(time.Time).UTC().Format(rfc3339), "incorrect value")
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}

func TestSemanticString(t *testing.T) {

	var testCases = []struct {
		input                     []byte
		expectAdditionalTypeValue uint64
		expectValue               string
	}{
		{
			// D8 20                       # tag(32)
			// 6C                          # text(12)
			//    687474703A2F2F612E636F6D # "http://a.com"
			input:                     []byte{0xd8, 0x20, 0x6c, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x2E, 0x63, 0x6F, 0x6D},
			expectAdditionalTypeValue: semanticURI,
			expectValue:               "http://a.com",
		},
		{
			// D8 23       # tag(35)
			// 64          # text(4)
			//    612A622B # "a*b+"
			input:                     []byte{0xd8, 0x23, 0x64, 0x61, 0x2A, 0x62, 0x2B},
			expectAdditionalTypeValue: semanticRegularExpression,
			expectValue:               "a*b+",
		},
		{
			// D8 24                              # tag(36)
			// 71                                 # text(17)
			// 4D494D452D56657273696F6E3A20312E30 # "MIME-Version: 1.0"
			input:                     []byte{0xd8, 0x24, 0x71, 0x4D, 0x49, 0x4D, 0x45, 0x2D, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6F, 0x6E, 0x3A, 0x20, 0x31, 0x2E, 0x30},
			expectAdditionalTypeValue: semanticMimeMessage,
			expectValue:               "MIME-Version: 1.0",
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)

		assert.Equal(t, MajorTypeSemantic, c[0].MajorType())
		assert.Equal(t, testCase.expectAdditionalTypeValue, c[0].AdditionalTypeValue())
		assert.Equal(t, testCase.expectValue, c[0].Value())
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}

func TestSemanticBase64(t *testing.T) {

	var testCases = []struct {
		input                     []byte
		expectAdditionalTypeValue uint64
		expectValue               string
	}{
		{
			// D8 21                               # tag(33)
			// 70                                  # text(16)
			//    6148523063446F764C32457559323974 # "aHR0cDovL2EuY29t"  // Decoded: "http://a.com"
			input:                     []byte{0xd8, 0x21, 0x70, 0x61, 0x48, 0x52, 0x30, 0x63, 0x44, 0x6F, 0x76, 0x4C, 0x32, 0x45, 0x75, 0x59, 0x32, 0x39, 0x74},
			expectAdditionalTypeValue: semanticBase64URL,
			expectValue:               "http://a.com",
		},
		{
			// D8 22       # tag(34)
			// 64          # text(4)
			//    51554A44 # "QUJD"  // Decoded: "ABC"
			input:                     []byte{0xd8, 0x22, 0x64, 0x51, 0x55, 0x4a, 0x44},
			expectAdditionalTypeValue: semanticBase64,
			expectValue:               "ABC",
		},
	}

	for _, testCase := range testCases {
		c, err := Decode(testCase.input)
		assert.Nil(t, err)

		assert.Equal(t, MajorTypeSemantic, c[0].MajorType())
		assert.Equal(t, testCase.expectAdditionalTypeValue, c[0].AdditionalTypeValue())
		buf, err := base64.RawStdEncoding.DecodeString(c[0].Value().(string))
		assert.Equal(t, testCase.expectValue, string(buf))
		assert.Equal(t, testCase.input, c[0].EncodeCBOR())
	}
}
