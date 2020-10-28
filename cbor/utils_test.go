package cbor

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataItemPrefix(t *testing.T) {

	var testCases = []struct {
		majorType MajorType
		length    uint64
		expect    []byte
	}{
		{
			majorType: MajorTypeTextString,
			length:    1,
			expect:    []byte{0x61},
		},
		{
			majorType: MajorTypeTextString,
			length:    23,
			expect:    []byte{0x77},
		},
		{
			majorType: MajorTypeTextString,
			length:    24,
			expect:    []byte{0x78, 0x18},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint8,
			expect:    []byte{0x78, 0xff},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint8 + 1,
			expect:    []byte{0x79, 0x01, 0x00},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint16,
			expect:    []byte{0x79, 0xff, 0xff},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint16 + 1,
			expect:    []byte{0x7a, 0x00, 0x01, 0x00, 0x00},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint32,
			expect:    []byte{0x7a, 0xff, 0xff, 0xff, 0xff},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint32 + 1,
			expect:    []byte{0x7b, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
		},
		{
			majorType: MajorTypeTextString,
			length:    math.MaxUint64,
			expect:    []byte{0x7b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expect, dataItemPrefix(testCase.majorType, testCase.length))
	}
}

func TestDebug(t *testing.T) {

	m := NewMap()
	m.Add(NewPositiveInteger8(1), NewNegativeInteger8(-1))
	m.Add(NewURI("http://a.com"), NewTextString("abc"))

	arr := NewArray()
	arr.Add(m)
}
