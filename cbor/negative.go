package cbor

/**
 * NegativeInteger value is the -1 minus the encoded unsigned integer.
 *
 * When decoding: decodedValue  = -1 - encodedValue
 * When encoding: encodedValue = -1 - decodedValue
 */
import (
	"fmt"
	"math"
)

// NegativeInteger8 wraps a negative integer with 8 bits (range: -1 to -127)
type NegativeInteger8 struct {
	baseNegativeInteger
}

// NegativeInteger16 wraps a negative integer with 16 bits (range: -128 to -32768)
type NegativeInteger16 struct {
	baseNegativeInteger
}

// NegativeInteger32 wraps a negative integer with 32 bits (range: -32769 to -2147483648)
type NegativeInteger32 struct {
	baseNegativeInteger
}

// NegativeInteger64 wraps a negative integer with 64 bits (range: -2147483649 to ...)
type NegativeInteger64 struct {
	baseNegativeInteger
}

// baseNegativeInteger wraps the base attributes of a negative integer
type baseNegativeInteger struct {
	baseDataItem
	V int64
}

// NewNegativeInteger8 returns a negative integer (8 bits) instance
func NewNegativeInteger8(value int64) *NegativeInteger8 {
	if value >= 0 {
		return nil
	}
	return &NegativeInteger8{
		baseNegativeInteger: newBaseNegativeInteger(int64(value), additionalType8Bits),
	}
}

// NewNegativeInteger16 returns a negative integer (16 bits) instance (range: -128 to -32768)
func NewNegativeInteger16(value int64) *NegativeInteger16 {
	if value >= 0 {
		return nil
	}
	return &NegativeInteger16{
		baseNegativeInteger: newBaseNegativeInteger(int64(value), additionalType16Bits),
	}
}

// NewNegativeInteger32 returns a negative integer (32 bits) instance
func NewNegativeInteger32(value int64) *NegativeInteger32 {
	if value >= 0 {
		return nil
	}
	return &NegativeInteger32{
		baseNegativeInteger: newBaseNegativeInteger(int64(value), additionalType32Bits),
	}
}

// NewNegativeInteger64 returns a negative integer (64 bits) instance
func NewNegativeInteger64(value int64) *NegativeInteger64 {
	if value >= 0 {
		return nil
	}
	return &NegativeInteger64{
		baseNegativeInteger: newBaseNegativeInteger(int64(value), additionalType64Bits),
	}
}

// newBaseNegativeInteger returns new base negative integer instance
func newBaseNegativeInteger(value int64, additionalType uint8) baseNegativeInteger {
	return baseNegativeInteger{
		baseDataItem: baseDataItem{
			majorType:      MajorTypeNegativeInt,
			additionalType: additionalType,
		},
		V: value,
	}
}

// AdditionalTypeValue returns the encoded uint value
func (b *baseNegativeInteger) AdditionalTypeValue() uint64 {
	return b.encodedValue()
}

// Value of this negative integer
func (b *baseNegativeInteger) Value() interface{} {
	switch {
	case b.V < math.MinInt32:
		return b.V
	case b.V < math.MinInt16:
		return int32(b.V)
	case b.V < math.MinInt8:
		return int16(b.V)
	default:
		return int8(b.V)
	}
}

// EncodeCBOR returns the CBOR binary representation of this negative integer
func (b *baseNegativeInteger) EncodeCBOR() []byte {
	return dataItemPrefix(MajorTypeNegativeInt, b.encodedValue())
}

// String representation of this negative integer
func (b *baseNegativeInteger) String() string {
	return fmt.Sprintf("NegativeInteger(%d)", b.V)
}

// ValueAsInt64 of this negative integer as int64
func (b *baseNegativeInteger) ValueAsInt64() int64 {
	return b.V
}

// encodedValue is the value represented as uint when transmitted via CBOR (ie. -1 - value)
func (b *baseNegativeInteger) encodedValue() uint64 {
	return uint64(int64(-1) - b.V)
}
