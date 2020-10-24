package cbor

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
)

// PositiveInteger8 wraps a positive integer with 8 bits
type PositiveInteger8 struct {
	basePositiveInteger
}

// PositiveInteger16 wraps a positive integer with 16 bits
type PositiveInteger16 struct {
	basePositiveInteger
}

// PositiveInteger32 wraps a positive integer with 32 bits
type PositiveInteger32 struct {
	basePositiveInteger
}

// PositiveInteger64 wraps a positive integer with 64 bits
type PositiveInteger64 struct {
	basePositiveInteger
}

type basePositiveInteger struct {
	baseDataItem
	V uint64
}

// NewPositiveInteger8 returns a positive integer (8 bits) instance
func NewPositiveInteger8(value uint8) *PositiveInteger8 {
	return &PositiveInteger8{
		basePositiveInteger: newBasePositiveInteger(uint64(value), additionalType8Bits),
	}
}

// ValueAsUint8 return value as uint8
func (p *PositiveInteger8) ValueAsUint8() uint8 {
	return uint8(p.V)
}

// ValueAsUint16 return value as uint16
func (p *PositiveInteger8) ValueAsUint16() uint16 {
	return uint16(p.V)
}

// ValueAsUint32 return value as uint16
func (p *PositiveInteger8) ValueAsUint32() uint32 {
	return uint32(p.V)
}

// ValueAsUint64 return value as uint16
func (p *PositiveInteger8) ValueAsUint64() uint64 {
	return p.V
}

// NewPositiveInteger16 returns a positive integer (16 bits) instance
func NewPositiveInteger16(value uint16) *PositiveInteger16 {
	return &PositiveInteger16{
		basePositiveInteger: newBasePositiveInteger(uint64(value), additionalType16Bits),
	}
}

// ValueAsUint16 return value as uint16
func (p *PositiveInteger16) ValueAsUint16() uint16 {
	return uint16(p.V)
}

// ValueAsUint32 return value as uint32
func (p *PositiveInteger16) ValueAsUint32() uint32 {
	return uint32(p.V)
}

// ValueAsUint64 return value as uint64
func (p *PositiveInteger16) ValueAsUint64() uint64 {
	return uint64(p.V)
}

// NewPositiveInteger32 returns a positive integer (32 bits) instance
func NewPositiveInteger32(value uint32) *PositiveInteger32 {
	return &PositiveInteger32{
		basePositiveInteger: newBasePositiveInteger(uint64(value), additionalType32Bits),
	}
}

// ValueAsUint32 return value as uint32
func (p *PositiveInteger32) ValueAsUint32() uint32 {
	return uint32(p.V)
}

// ValueAsUint64 return value as uint64
func (p *PositiveInteger32) ValueAsUint64() uint64 {
	return uint64(p.V)
}

// NewPositiveInteger64 returns a positive integer (64 bits) instance
func NewPositiveInteger64(value uint64) *PositiveInteger64 {
	return &PositiveInteger64{
		basePositiveInteger: newBasePositiveInteger(uint64(value), additionalType64Bits),
	}
}

// ValueAsUint64 return value as uint64
func (p *PositiveInteger64) ValueAsUint64() uint64 {
	return uint64(p.V)
}

// NewPositiveInteger returns a positive integer using the most compact
// struct that would fit the value
func NewPositiveInteger(value uint64) DataItem {
	switch {
	case value > math.MaxUint32:
		return NewPositiveInteger64(value)
	case value > math.MaxUint16:
		return NewPositiveInteger32(uint32(value))
	case value > math.MaxUint8:
		return NewPositiveInteger16(uint16(value))
	default:
		return NewPositiveInteger8(uint8(value))
	}
}

// newBasePositiveInteger returns new base positive integer instance
func newBasePositiveInteger(value uint64, additionalType uint8) basePositiveInteger {
	return basePositiveInteger{
		baseDataItem: baseDataItem{
			majorType:      MajorTypePositiveInt,
			additionalType: additionalType,
		},
		V: value,
	}
}

// AdditionalType of this positive integer
func (b *basePositiveInteger) AdditionalType() uint8 {
	switch {
	case b.V > math.MaxUint32:
		return additionalType64Bits
	case b.V > math.MaxUint16:
		return additionalType32Bits
	case b.V > math.MaxUint8:
		return additionalType16Bits
	case b.V > uint64(additionalTypeDirectValue23):
		return additionalType8Bits
	default:
		// value is between 0-23, use this as the additional type
		return uint8(b.V)
	}
}

// AdditionalTypeValue returns the uint value as uint64
func (b *basePositiveInteger) AdditionalTypeValue() uint64 {
	return b.V
}

// Value of this positive integer
func (b *basePositiveInteger) Value() interface{} {
	switch {
	case b.V > math.MaxUint32:
		return b.V
	case b.V > math.MaxUint16:
		return uint32(b.V)
	case b.V > math.MaxUint8:
		return uint16(b.V)
	default:
		return uint8(b.V)
	}
}

// ValueAsInt64 of this positive integer as int64
func (b *basePositiveInteger) ValueAsUInt64() uint64 {
	return b.V
}

// EncodeCBOR returns the CBOR binary representation of this positive integer
func (b *basePositiveInteger) EncodeCBOR() []byte {
	return dataItemPrefix(MajorTypePositiveInt, b.V)
}

// String representation of this positive integer
func (b *basePositiveInteger) String() string {

	typeString := ""

	switch {
	case b.additionalType == additionalType64Bits:
		typeString = "PositiveInteger64"
		break
	case b.additionalType == additionalType32Bits:
		typeString = "PositiveInteger32"
		break
	case b.additionalType == additionalType16Bits:
		typeString = "PositiveInteger16"
		break
	case b.additionalType == additionalType8Bits,
		b.additionalType <= additionalTypeDirectValue23:
		typeString = "PositiveInteger8"
		break
	default:
		log.WithField("additionalType", b.additionalType).Error("Unknown additional type")
	}

	return fmt.Sprintf("%s(%d)", typeString, b.V)
}
