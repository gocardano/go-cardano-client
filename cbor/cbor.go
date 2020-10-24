package cbor

import (
	"github.com/gocardano/go-cardano-client/errors"

	log "github.com/sirupsen/logrus"
)

// DataItem represent the interface for all CBOR encoded data item
type DataItem interface {
	MajorType() MajorType
	AdditionalType() uint8
	AdditionalTypeValue() uint64
	Value() interface{}
	EncodeCBOR() []byte
	String() string
}

// baseDataItem includes attributes for all base data item structs
type baseDataItem struct {
	majorType      MajorType
	additionalType uint8
}

// MajorType of this data item
func (b *baseDataItem) MajorType() MajorType {
	return b.majorType
}

// AdditionalType of this data item
func (b *baseDataItem) AdditionalType() uint8 {
	return b.additionalType
}

// Decode binary data and return list of CBOR encoded data items
func Decode(data []byte) ([]DataItem, error) {
	r := NewBitstreamReader(data)
	result := []DataItem{}
	for r.HasMoreBits() {
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error while parsing for more bits, giving up.")
			return result, err
		}
		result = append(result, obj)
	}

	return result, nil
}

// doGetAdditionalType is called after reading the 3 majorType bits.  It then read
// the next 5 bits (which represent the additional type).  And depending on the value:
//
//  - 0-23 (use the same value as the additionalTypeValue)
//  - 24   (use the next uint8 as the additionalTypeValue)
//  - 25   (use the next uint16 as the additionalTypeValue)
//  - 26   (use the next uint32 as the additionalTypeValue)
//  - 27   (use the next uint64 as the additionalTypeValue)
//  - 31   (indefinite tag, return 0 for additionalTypeValue)
//
// For positiveInt/negativeInt, the additionalType returned is the actual uint_xx value
//
// Returns format: additionalType, additionalTypeValue, error
func doGetAdditionalType(r *BitstreamReader) (uint8, uint64, error) {

	additionalType, err := r.ReadBitsAsUint8(5)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case additionalType <= additionalTypeDirectValue23:
		// if additionalType is within the range 0-23, use this as the additionalTypeValue
		return additionalType, uint64(additionalType), nil
	case additionalType == additionalType8Bits:
		additionalTypeValue, err := r.ReadBitsAsUint64(8)
		return additionalType, additionalTypeValue, err
	case additionalType == additionalType16Bits:
		additionalTypeValue, err := r.ReadBitsAsUint64(16)
		return additionalType, additionalTypeValue, err
	case additionalType == additionalType32Bits:
		additionalTypeValue, err := r.ReadBitsAsUint64(32)
		return additionalType, additionalTypeValue, err
	case additionalType == additionalType64Bits:
		additionalTypeValue, err := r.ReadBitsAsUint64(64)
		return additionalType, additionalTypeValue, err
	case additionalType == additionalTypeIndefinite:
		return additionalType, 0, err
	default:
		log.WithField("additionalType", additionalType).Error("Unhandled additional type")
		return additionalType, 0, errors.NewMessageErrorf(errors.ErrCborAdditionalTypeUnhandled, "Unhandled additional type [%d];", additionalType)
	}
}

// doGetNextDataItem parse the next CBOR encoded data item.  Assumes that the reader
// is positioned to read the next major type (3 bits).
func doGetNextDataItem(r *BitstreamReader) (DataItem, error) {

	majorValue, err := r.ReadBitsAsUint64(3)
	if err != nil {
		log.Error("Error reading major value of object", err)
		return nil, err
	}

	switch MajorType(majorValue) {
	case MajorTypePositiveInt:
		return r.decodePositiveUnsignedInt()
	case MajorTypeNegativeInt:
		return r.decodeNegativeInt()
	case MajorTypeByteString:
		return r.decodeByteString()
	case MajorTypeTextString:
		return r.decodeTextString()
	case MajorTypeArray:
		return r.decodeArray()
	case MajorTypeMap:
		return r.decodeMap()
	case MajorTypeSemantic:
		return r.decodeSemantic()
	case MajorTypePrimitive:
		return r.decodePrimitive()
	}

	return nil, errors.NewError(errors.ErrCborMajorTypeUnhandled)
}
