package cbor

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/x448/float16"
)

const (
	primitiveSimpleValueMin = uint8(32)
	primitiveSimpleValueMax = uint8(255)
)

type basePrimitive struct {
	baseDataItem
}

// AdditionalTypeValue returns the additionalTypeValue
func (b basePrimitive) AdditionalTypeValue() uint64 {
	return 0
}

// PrimitiveFalse represents a primitive false object
type PrimitiveFalse struct {
	basePrimitive
}

// PrimitiveTrue represents a primitive true object
type PrimitiveTrue struct {
	basePrimitive
}

// PrimitiveNull represents a primitive null object
type PrimitiveNull struct {
	basePrimitive
}

// PrimitiveUndefined represents a primitive undefined object
type PrimitiveUndefined struct {
	basePrimitive
}

// PrimitiveSimpleValue represents a simple value (uint8_t: 32..255)
type PrimitiveSimpleValue struct {
	baseDataItem
	V uint8
}

// PrimitiveHalfPrecisionFloat represents a half precision float
// Format: sign (1 bit); exponent (5 bits); fraction (10 bits)
// Range ±65,504, with precision up to 0.0000000596046.
type PrimitiveHalfPrecisionFloat struct {
	baseDataItem
	V float16.Float16
}

// PrimitiveSinglePrecisionFloat represents a single precision float
// Format: sign (1 bit); exponent (8 bits); fraction (23 bits)
type PrimitiveSinglePrecisionFloat struct {
	baseDataItem
	V float32
}

// PrimitiveDoublePrecisionFloat represents a double precision float
// Format: sign (1 bit); exponent (11 bits); fraction (52 bits)
type PrimitiveDoublePrecisionFloat struct {
	baseDataItem
	V float64
}

// PrimitiveBreakStopCode represent a break stop code
type PrimitiveBreakStopCode struct {
	basePrimitive
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveFalse returns instance of PrimitiveFalse struct
func NewPrimitiveFalse() *PrimitiveFalse {
	return &PrimitiveFalse{
		basePrimitive: basePrimitive{
			baseDataItem: baseDataItem{
				majorType:      MajorTypePrimitive,
				additionalType: primitiveFalse,
			},
		},
	}
}

// Value returns the boolean false
func (p *PrimitiveFalse) Value() interface{} {
	return false
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveFalse) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveFalse}
}

// String returns description of this item
func (p *PrimitiveFalse) String() string {
	return "False"
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveTrue returns instance of PrimitiveTrue struct
func NewPrimitiveTrue() *PrimitiveTrue {
	return &PrimitiveTrue{
		basePrimitive: basePrimitive{
			baseDataItem: baseDataItem{
				majorType:      MajorTypePrimitive,
				additionalType: primitiveTrue,
			},
		},
	}
}

// Value returns the boolean true
func (p *PrimitiveTrue) Value() interface{} {
	return true
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveTrue) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveTrue}
}

// String returns description of this item
func (p *PrimitiveTrue) String() string {
	return "False"
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveNull returns instance of PrimitiveNull struct
func NewPrimitiveNull() *PrimitiveNull {
	return &PrimitiveNull{
		basePrimitive: basePrimitive{
			baseDataItem: baseDataItem{
				majorType:      MajorTypePrimitive,
				additionalType: primitiveNull,
			},
		},
	}
}

// Value returns a nil item
func (p *PrimitiveNull) Value() interface{} {
	return nil
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveNull) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveNull}
}

// String returns description of this item
func (p *PrimitiveNull) String() string {
	return "Null"
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveUndefined returns instance of PrimitiveUndefined struct
func NewPrimitiveUndefined() *PrimitiveUndefined {
	return &PrimitiveUndefined{
		basePrimitive: basePrimitive{
			baseDataItem: baseDataItem{
				majorType:      MajorTypePrimitive,
				additionalType: primitiveUndefined,
			},
		},
	}
}

// Value returns a nil item to represent an undefined
func (p *PrimitiveUndefined) Value() interface{} {
	return nil
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveUndefined) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveUndefined}
}

// String returns description of this item
func (p *PrimitiveUndefined) String() string {
	return "Undefined"
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveSimpleValue returns instance of PrimitiveSimpleValue struct
func NewPrimitiveSimpleValue(value uint8) *PrimitiveSimpleValue {
	if value < primitiveSimpleValueMin || primitiveSimpleValueMax < value {
		log.Errorf("Primitive simple values should be within the range [%d-%d], received [%d]",
			primitiveSimpleValueMin, primitiveSimpleValueMax, value)
		return nil
	}
	return &PrimitiveSimpleValue{
		baseDataItem: baseDataItem{
			majorType:      MajorTypePrimitive,
			additionalType: primitiveSimpleValue,
		},
		V: value,
	}
}

// AdditionalTypeValue returns the value as int64
func (p *PrimitiveSimpleValue) AdditionalTypeValue() uint64 {
	return uint64(p.V)
}

// Value returns the simple value
func (p *PrimitiveSimpleValue) Value() interface{} {
	return p.V
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveSimpleValue) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveSimpleValue, byte(p.V)}
}

// String returns description of this item
func (p *PrimitiveSimpleValue) String() string {
	return fmt.Sprintf("Simple Value - Value: [%d]", p.V)
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveHalfPrecisionFloat returns instance of PrimitiveHalfPrecisionFloat struct
func NewPrimitiveHalfPrecisionFloat(value uint16) *PrimitiveHalfPrecisionFloat {
	return &PrimitiveHalfPrecisionFloat{
		baseDataItem: baseDataItem{
			majorType:      MajorTypePrimitive,
			additionalType: primitiveHalfPrecisionFloat,
		},
		V: float16.Frombits(value),
	}
}

// AdditionalTypeValue returns the value as uint64
func (p *PrimitiveHalfPrecisionFloat) AdditionalTypeValue() uint64 {
	return uint64(p.V.Bits())
}

// Value returns the half precision float value
func (p *PrimitiveHalfPrecisionFloat) Value() interface{} {
	return p.V.Float32()
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveHalfPrecisionFloat) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveHalfPrecisionFloat,
		byte(p.V.Bits() >> 8),
		byte(p.V.Bits()),
	}
}

// String returns description of this item
func (p *PrimitiveHalfPrecisionFloat) String() string {
	return fmt.Sprintf("Half Precision Fault - Value: [%f]", p.V.Float32())
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveSinglePrecisionFloat returns instance of PrimitiveSinglePrecisionFloat struct
func NewPrimitiveSinglePrecisionFloat(value float32) *PrimitiveSinglePrecisionFloat {
	return &PrimitiveSinglePrecisionFloat{
		baseDataItem: baseDataItem{
			majorType:      MajorTypePrimitive,
			additionalType: primitiveSinglePrecisionFloat,
		},
		V: value,
	}

}

// AdditionalTypeValue returns the value as uint64
func (p *PrimitiveSinglePrecisionFloat) AdditionalTypeValue() uint64 {
	return uint64(p.V)
}

// Value returns the single precision float value
func (p *PrimitiveSinglePrecisionFloat) Value() interface{} {
	return p.V
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveSinglePrecisionFloat) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveSinglePrecisionFloat,
		byte(math.Float32bits(p.V) >> 24),
		byte(math.Float32bits(p.V) >> 16),
		byte(math.Float32bits(p.V) >> 8),
		byte(math.Float32bits(p.V)),
	}
}

// String returns description of this item
func (p *PrimitiveSinglePrecisionFloat) String() string {
	return fmt.Sprintf("Single Precision Fault - Value: [%f]", p.V)
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveDoublePrecisionFloat returns instance of PrimitiveDoublePrecisionFloat struct
func NewPrimitiveDoublePrecisionFloat(value float64) *PrimitiveDoublePrecisionFloat {
	return &PrimitiveDoublePrecisionFloat{
		baseDataItem: baseDataItem{
			majorType:      MajorTypePrimitive,
			additionalType: primitiveDoublePrecisionFloat,
		},
		V: value,
	}
}

// AdditionalTypeValue returns the value as uint64
func (p *PrimitiveDoublePrecisionFloat) AdditionalTypeValue() uint64 {
	return math.Float64bits(p.V)
}

// Value returns the double precision float value
func (p *PrimitiveDoublePrecisionFloat) Value() interface{} {
	return p.V
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveDoublePrecisionFloat) EncodeCBOR() []byte {
	return dataItemPrefix(MajorTypePrimitive, math.Float64bits(p.V))
}

// String returns description of this item
func (p *PrimitiveDoublePrecisionFloat) String() string {
	return fmt.Sprintf("Double Precision Fault - Value: [%f]", p.V)
}

////////////////////////////////////////////////////////////////////////////////

// NewPrimitiveBreakStopCode returns instance of PrimitiveFalse struct
func NewPrimitiveBreakStopCode() *PrimitiveBreakStopCode {
	return &PrimitiveBreakStopCode{
		basePrimitive: basePrimitive{
			baseDataItem: baseDataItem{
				majorType:      MajorTypePrimitive,
				additionalType: primitiveBreakStopCode,
			},
		},
	}
}

// Value returns a nil since it is not used
func (p *PrimitiveBreakStopCode) Value() interface{} {
	return nil
}

// EncodeCBOR returns CBOR representation for this item
func (p *PrimitiveBreakStopCode) EncodeCBOR() []byte {
	return []byte{MajorTypePrimitive.EncodeCBOR() | primitiveBreakStopCode}
}

// String returns description of this item
func (p *PrimitiveBreakStopCode) String() string {
	return "Break stop code"
}
