package cbor

import (
	"math"
	"math/big"

	"github.com/gocardano/go-cardano-client/errors"
	log "github.com/sirupsen/logrus"
	"github.com/x448/float16"
)

// decodeArray parses the next array object.
// Only called after the majorType array has been determined.
func (r *BitstreamReader) decodeArray() (*Array, error) {

	array := NewArray()

	// additionalTypeValue (second parameter) in this case indicates the size of the array
	additionalType, arrayLength, err := doGetAdditionalType(r)
	if err != nil {
		return nil, err
	}

	log.Tracef("Starting to iterate on array with length: %d", arrayLength)

	hasMoreItems := false
	counter := uint64(0)

	if arrayLength > 0 || additionalType == additionalTypeIndefinite {
		hasMoreItems = true
	}

	for hasMoreItems {

		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error reading next array item")
			return array, err
		}

		log.Tracef("Found another item in the array %+v", obj)

		if additionalType == additionalTypeIndefinite &&
			byte(obj.MajorType()) == indefiniteBreakCodeMajorType &&
			obj.AdditionalType() == indefiniteBreakCodeAdditionalType {
			log.Tracef("Array of indefinite length reached the break stop code, found [%d] items", array.Length())
			hasMoreItems = false
		} else {
			array.Add(obj)
			counter++
			if counter == arrayLength {
				log.Tracef("Array of [%d] length reached [%d] items", arrayLength, array.Length())
				hasMoreItems = false
			}
		}
	}

	return array, nil
}

// decodeByteString parses the next byte string object.
// Only called after the majorType byte string has been determined.
func (r *BitstreamReader) decodeByteString() (*ByteString, error) {
	result, err := r.doDecodeByteString()
	if err != nil {
		return nil, err
	}
	return NewByteString(result), nil
}

// decodeMap parses the next map object.  Only called after the majorType map has been determined.
func (r *BitstreamReader) decodeMap() (DataItem, error) {

	m := NewMap()

	// additionalTypeValue (second parameter) in this case indicates the size of the map
	additionalType, mapLength, err := doGetAdditionalType(r)
	if err != nil {
		return nil, err
	}

	log.Tracef("Starting to iterate on map with length: %d", mapLength)

	hasMoreItems := true
	counter := uint64(0)

	for hasMoreItems {
		key, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error reading map key item")
			return m, err
		}
		if additionalType == additionalTypeIndefinite &&
			key.MajorType() == MajorTypePrimitive &&
			key.AdditionalType() == primitiveBreakStopCode {
			// Found break stop code
			hasMoreItems = false
		} else {

			value, err := doGetNextDataItem(r)
			if err != nil {
				log.Error("Error reading map value item")
				return m, err
			}

			log.Tracef("Adding map key: [%+v] with value: [%+v]", key, value)
			m.Add(key, value)

			counter++
			if counter == mapLength {
				log.Tracef("Map of [%d] length reached [%d] items", mapLength, m.Length())
				hasMoreItems = false
			}
		}
	}

	return m, nil
}

// decodeNegativeInt parses the next negative integer object.
// Only called after the majorType negative integer has been determined.
func (r *BitstreamReader) decodeNegativeInt() (DataItem, error) {

	var result DataItem

	// encodedValue (2nd return parameter) for positiveInt/negativeInt is the actual encoded value
	additionalType, encodedValue, err := doGetAdditionalType(r)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"encodedValue": encodedValue,
		}).Error("Unable to handle negative int")
		return nil, err
	}

	// compute the actual value
	actualValue := int64(-1) - int64(encodedValue)

	log.WithFields(log.Fields{
		"encodedValue": encodedValue,
		"actualValue":  actualValue,
	}).Debug("Decoded CBOR item from array")

	switch additionalType {
	case additionalType64Bits:
		result = NewNegativeInteger64(actualValue)
		break
	case additionalType32Bits:
		result = NewNegativeInteger32(actualValue)
		break
	case additionalType16Bits:
		result = NewNegativeInteger16(actualValue)
		break
	default:
		result = NewNegativeInteger8(actualValue)
		break
	}

	if result == nil {
		log.Error("Error creating negative integer instance")
	}

	return result, nil
}

// decodePositiveUnsignedInt parses the next positive unsigned integer.
// Only called after the majorType negative integer has been determined.
func (r *BitstreamReader) decodePositiveUnsignedInt() (DataItem, error) {

	var result DataItem

	// actualValue (2nd return parameter) for positiveInt/negativeInt is the unsigned positive value
	additionalType, actualValue, err := doGetAdditionalType(r)
	if err != nil {
		return nil, err
	}

	switch {
	case additionalType == additionalType64Bits:
		result = NewPositiveInteger64(actualValue)
		break
	case additionalType == additionalType32Bits:
		result = NewPositiveInteger32(uint32(actualValue))
		break
	case additionalType == additionalType16Bits:
		result = NewPositiveInteger16(uint16(actualValue))
		break
	case additionalType == additionalType8Bits:
		result = NewPositiveInteger8(uint8(actualValue))
		break
	case actualValue <= uint64(additionalTypeDirectValue23):
		result = NewPositiveInteger8(uint8(actualValue))
		break
	default:
		log.WithFields(log.Fields{
			"additionalType": additionalType,
			"actualValue":    actualValue,
		}).Error("Unhandled additional type while parsing positive unsigned int")
		break
	}

	return result, nil
}

// decodePrimitive parses the next primitive object.
// Only called after the majorType primitive has been determined.
func (r *BitstreamReader) decodePrimitive() (DataItem, error) {

	var obj DataItem

	additionalType, err := r.ReadBitsAsUint8(5)
	if err != nil {
		log.Error("Error parsing additional type for primitive item")
		return nil, err
	}

	switch additionalType {
	case primitiveFalse:
		obj = NewPrimitiveFalse()
		break
	case primitiveTrue:
		obj = NewPrimitiveTrue()
		break
	case primitiveNull:
		obj = NewPrimitiveNull()
		break
	case primitiveUndefined:
		obj = NewPrimitiveUndefined()
		break
	case primitiveSimpleValue:
		val, err := r.ReadBitsAsUint8(8)
		if err != nil {
			return nil, err
		}
		obj = NewPrimitiveSimpleValue(val)
		break
	case primitiveHalfPrecisionFloat:
		val, err := r.ReadBitsAsUint16(16)
		if err != nil {
			return nil, err
		}
		float16.Frombits(val)
		obj = NewPrimitiveHalfPrecisionFloat(val)
		break
	case primitiveSinglePrecisionFloat:
		val, err := r.ReadBitsAsUint32(32)
		if err != nil {
			return nil, err
		}
		obj = NewPrimitiveSinglePrecisionFloat(math.Float32frombits(val))
		break
	case primitiveDoublePrecisionFloat:
		val, err := r.ReadBitsAsUint64(64)
		if err != nil {
			return nil, err
		}
		obj = NewPrimitiveDoublePrecisionFloat(math.Float64frombits(val))
		break
	case primitiveBreakStopCode:
		obj = NewPrimitiveBreakStopCode()
		break
	}

	return obj, nil
}

// decodeSemantic parses the next semantic tagged object.
// Only called after the majorType semantic has been determined.
func (r *BitstreamReader) decodeSemantic() (DataItem, error) {

	var result DataItem

	// semanticTagID (second parameter) in this case indicates the type of the semantic tag
	_, semanticTagID, err := doGetAdditionalType(r)
	if err != nil {
		return nil, err
	}

	switch semanticTagID {
	case semanticDateTimeString:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next text string for semantic date time string", err)
			return nil, err
		}
		if obj.MajorType() != MajorTypeTextString {
			log.Error("Unhandled exception, expected to see text as date time string object")
			return nil, errors.NewError(errors.ErrCborMajorTypeUnhandled)
		}
		result = NewDateTimeString(obj.Value().(string))
		break
	case semanticDateTimeEpoch:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next positive/negative int for semantic date time epoch", err)
			return nil, err
		}
		epoch := int64(0)
		switch obj.MajorType() {
		case MajorTypePositiveInt:
			switch obj.AdditionalType() {
			case additionalType64Bits:
				epoch = int64(obj.Value().(uint64))
				break
			case additionalType32Bits:
				epoch = int64(obj.Value().(uint32))
				break
			case additionalType16Bits:
				epoch = int64(obj.Value().(uint16))
				break
			default:
				// additional type value below 24 (use same value)
				epoch = int64(obj.Value().(uint8))
				break
			}
			result = NewDateTimeEpoch(epoch)
			break
		case MajorTypeNegativeInt:
			switch obj.AdditionalType() {
			case additionalType64Bits:
				epoch = int64(obj.(*NegativeInteger64).ValueAsInt64())
				break
			case additionalType32Bits:
				epoch = int64(obj.(*NegativeInteger32).ValueAsInt64())
				break
			case additionalType16Bits:
				epoch = int64(obj.(*NegativeInteger16).ValueAsInt64())
				break
			default:
				// additional type value below 24 (use same value)
				epoch = int64(obj.(*NegativeInteger8).ValueAsInt64())
				break
			}
			result = NewDateTimeEpoch(epoch)
			break
		default:
			log.Errorf("Unhandled major type [%d] while parsing the value for date time epoch due to [%s]", obj.MajorType(), err)
			result = nil
			break
		}
		break
	case semanticPositiveBignum:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next string for positive bignum", err)
			return nil, err
		}
		if obj.MajorType() != MajorTypeByteString {
			log.Errorf("Expected bytestring payload in positive bignum, unhandled major type: %d", obj.MajorType())
			return nil, err
		}
		n := new(big.Int)
		n = n.SetBytes(obj.(*ByteString).ValueAsBytes())
		result = NewPositiveBignumber(n)
		break
	case semanticNegativeBignum:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next string for negative bignum", err)
			return nil, err
		}
		if obj.MajorType() != MajorTypeByteString {
			log.Errorf("Expected bytestring payload in negative bignum, unhandled major type: %d", obj.MajorType())
			return nil, err
		}
		n := new(big.Int)
		n = n.SetBytes(obj.(*ByteString).ValueAsBytes())
		n = n.Sub(big.NewInt(-1), n)
		result = NewNegativeBignumber(n)
		break
	case semanticURI:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next string for uri", err)
			return nil, err
		}
		result = NewURI(obj.Value().(string))
		break
	case semanticBase64URL:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next text string for base64 string", err)
			return nil, err
		}
		result = NewBase64URL(obj.Value().(string))
		break
	case semanticBase64:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next text string for base64 string", err)
			return nil, err
		}
		result = NewBase64String(obj.Value().(string))
		break
	case semanticRegularExpression:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next text string for regular expression string", err)
			return nil, err
		}
		result = NewRegularExpression(obj.Value().(string))
		break
	case semanticMimeMessage:
		obj, err := doGetNextDataItem(r)
		if err != nil {
			log.Error("Error trying to parse the next text string for mime message string", err)
			return nil, err
		}
		result = NewMimeMessage(obj.Value().(string))
		break
	case semanticDecimalFraction,
		semanticBigFloat,
		semanticExpectedConversionToBase64URL,
		semanticExpectedConversionToBase64,
		semanticExpectedConversionToBase16,
		semanticEncodedCBORDataItems,
		semanticSelfDescribeCBOR:
		log.Infof("Semantic unhandled tag %d", semanticTagID)
		// TBD
		return nil, errors.NewErrorf(errors.ErrCborAdditionalTypeUnhandled, "Unable to parse due to unhandled semantic tag ID encountered: %d", semanticTagID)
	}

	return result, nil
}

// decodeTextString parses the next text string object.
// Only called after the majorType text string has been determined.
func (r *BitstreamReader) decodeTextString() (*TextString, error) {
	result, err := r.doDecodeByteString()
	if err != nil {
		return nil, err
	}
	return NewTextString(string(result)), nil
}

// decodeByteString handles parsing the next byte string or text string.
// Only called after the majorType byteString/textString has been determined.
func (r *BitstreamReader) doDecodeByteString() ([]byte, error) {

	// byteLength (second parameter) in this case indicates the length of the byte/text
	additionalType, byteLength, err := doGetAdditionalType(r)
	if err != nil {
		return nil, err
	}

	payload := []byte{}

	if additionalType != additionalTypeIndefinite {

		log.Tracef("Reading bytes of payload length: %d", byteLength)
		payload, err = r.ReadBytes(byteLength)
		if err != nil {
			return nil, err
		}

	} else {

		chunkToken, err := r.ReadBitsAsUint64(8)
		log.Tracef("ChunkLength: 0x%02x", chunkToken)
		if err != nil {
			return nil, err
		}
		for chunkToken != indefiniteBreakCode {

			// ignore first 3 bits, only the 5 bits matters for length
			chunkLength := chunkToken & 0x1f
			tmp, err := r.ReadBytes(chunkLength)
			if err != nil {
				return nil, err
			}
			log.Tracef("Read [%d] chunk payload", len(tmp))
			payload = append(payload, tmp...)

			chunkToken, err = r.ReadBitsAsUint64(8)
			log.Tracef("ChunkLength: 0x%02x", chunkToken)
			if err != nil {
				return nil, err
			}
		}
	}

	return payload, nil
}
