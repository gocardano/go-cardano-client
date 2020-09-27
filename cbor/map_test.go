package cbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSimple(t *testing.T) {

	input := []byte{0xa1, 0x01, 0x02}

	c, err := Decode(input)
	assert.Nil(t, err)
	assert.Equal(t, MajorTypeMap, c[0].MajorType())

	// Test: Check size of map
	m := c[0].(*Map)
	assert.Equal(t, 1, m.Length())

	// Test: Check generated CBOR
	assert.Equal(t, input, c[0].EncodeCBOR())
}

func TestMapIndefiniteLengthMultipleEntries(t *testing.T) {

	input := []byte{
		0xbf,                   // indefinite length map
		0x63, 0x46, 0x75, 0x6e, // Key: Text("Fun")
		0xf5,                   // Value: true
		0x63, 0x41, 0x6d, 0x74, // Key: Text("Amt")
		0x21, // Value: -2
		0xff, // break
	}

	c, err := Decode(input)
	assert.Nil(t, err)
	assert.Equal(t, MajorTypeMap, c[0].MajorType())

	m := c[0].(*Map)
	assert.Equal(t, 2, m.Length())
}

func TestMapIndefiniteLengthEncoding(t *testing.T) {

	// Test with only one entry in map, so encoding will be consistent

	input := []byte{
		0xbf,                   // indefinite length map
		0x63, 0x46, 0x75, 0x6e, // Key: Text("Fun")
		0xf5, // Value: true
		0xff, // break
	}

	c, err := Decode(input)
	assert.Nil(t, err)
	assert.Equal(t, MajorTypeMap, c[0].MajorType())

	// Parse the generated encoded CBOR
	assert.Equal(t, input, c[0].(*Map).doEncodeCBOR(false))
}
