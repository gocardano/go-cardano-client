package cbor

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {

	log.SetLevel(log.TraceLevel)

	testCases := []struct {
		scenario     string
		input        []byte
		expectValues []interface{}
		fixedLength  bool
	}{
		{
			scenario:     "1: simple one item",
			input:        []byte{0x81, 0x01},
			expectValues: []interface{}{uint8(0x01)},
			fixedLength:  true,
		},
		{
			scenario: "2: 4 items of various types (int, neg, byte, string)",
			input: []byte{0x84,
				0x01,       // 1
				0x20,       // -1
				0x41, 0x61, // []byte("a")
				0x61, 0x61, // "a"
			},
			expectValues: []interface{}{
				uint8(1),
				int8(-1),
				[]byte{0x61},
				"a",
			},
			fixedLength: true,
		},
		{
			scenario: "3: 4 items of various types (primitives, floats)",
			input: []byte{0x86,
				0x19, 0xff, 0xff, // 65535
				0x39, 0xff, 0xfe, // -65535
				0xf4,             // false
				0xf5,             // true
				0xf9, 0x3c, 0x00, // 1.0
				0xfb, 0x40, 0xC1, 0x5C, 0x70, 0xA3, 0xD7, 0x0A, 0x3D, // 8888.88
			},
			expectValues: []interface{}{
				uint16(65535),
				int32(-65535),
				false,
				true,
				float32(1.0),
				8888.88,
			},
			fixedLength: true,
		},
		{
			scenario:     "4: indefinite length",
			input:        []byte{0x9f, 0x01, 0x02, 0x03, 0xff},
			expectValues: []interface{}{uint8(1), uint8(2), uint8(3)},
			fixedLength:  false,
		},
		{
			scenario:     "5: empty array",
			input:        []byte{0x80},
			expectValues: []interface{}{},
			fixedLength:  true,
		},
	}

	for _, testCase := range testCases {

		t.Log("Testing " + testCase.scenario)
		c, err := Decode(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, MajorTypeArray, c[0].MajorType())

		expectArray := testCase.expectValues
		actualArray := c[0].Value().([]DataItem)
		assert.Equal(t, len(expectArray), len(actualArray))

		for i := 0; i < len(expectArray); i++ {
			assert.Equal(t, expectArray[i], actualArray[i].Value())
		}

		assert.Equal(t, testCase.input, c[0].(*Array).doEncodeCBOR(testCase.fixedLength))
	}
}

func TestRecursiveArray(t *testing.T) {

	// Test scenario from cbor.me:
	// [1, [-1, ["b"]]]
	//
	// Expected:
	// 82             # array(2)
	//    01          # unsigned(1)
	//    82          # array(2)
	//       20       # negative(0)
	//       81       # array(1)
	//          61    # text(1)
	//             62 # "b"

	input := []byte{0x82, 0x01, 0x82, 0x20, 0x81, 0x61, 0x62}

	c, err := Decode(input)
	assert.Nil(t, err)

	// First Level Array
	arr1 := c[0].(*Array)
	assert.Equal(t, MajorTypeArray, arr1.MajorType())
	assert.Equal(t, 2, arr1.Length())
	assert.Equal(t, uint8(1), arr1.List()[0].Value())

	// Second Level Array
	arr2 := arr1.List()[1].(*Array)
	assert.Equal(t, MajorTypeArray, arr2.MajorType())
	assert.Equal(t, 2, arr2.Length())
	assert.Equal(t, int8(-1), arr2.List()[0].Value())

	// Third Level Array
	arr3 := arr2.List()[1].(*Array)
	assert.Equal(t, MajorTypeArray, arr3.MajorType())
	assert.Equal(t, 1, arr3.Length())
	assert.Equal(t, MajorTypeTextString, arr3.List()[0].MajorType())
	assert.Equal(t, "b", arr3.List()[0].Value())

	assert.Equal(t, input, c[0].(*Array).EncodeCBOR())

}
