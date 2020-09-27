package cbor

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
)

const (
	// maximum chunk length is 31 (last 5 bits = 0x1f)
	defaultMaxChunkLength uint8 = 31
)

// ByteString represents a CBOR byte string
type ByteString struct {
	baseByteString
}

// baseByteString wraps base attributes for byte and text
type baseByteString struct {
	baseDataItem
	V                   []byte
	length              uint64
	chunkEncodingLength uint64
}

// NewByteString returns a new byte string instance
func NewByteString(value []byte) *ByteString {
	return &ByteString{
		baseByteString: newBaseByteString(MajorTypeByteString, value),
	}
}

// ValueAsBytes returns the actual bytes
func (b *ByteString) ValueAsBytes() []byte {
	return b.V
}

// newBaseByteString return a baseByteString
func newBaseByteString(majorType MajorType, value []byte) baseByteString {
	return baseByteString{
		baseDataItem: baseDataItem{
			majorType: majorType,
		},
		V:      value,
		length: uint64(len(value)),
	}
}

// AdditionalTypeValue returns the length of the byte string
func (b baseByteString) AdditionalTypeValue() uint64 {
	return b.length
}

// EncodeCBOR returns CBOR representation for this item
func (b baseByteString) EncodeCBOR() []byte {

	var result []byte
	switch {
	case b.length > math.MaxUint64:
		result = b.encodeAsCBORChunks(defaultMaxChunkLength)
		break
	default:
		result = append(dataItemPrefix(b.majorType, b.length), b.V...)
	}

	return result
}

// Value returns the object
func (b baseByteString) Value() interface{} {
	if b.majorType == MajorTypeByteString {
		return b.V
	}
	return string(b.V)
}

// ValueAsString return the object as string
func (b baseByteString) ValueAsString() string {
	return string(b.V)
}

// String returns description of this item
func (b *baseByteString) String() string {
	if b.majorType == MajorTypeByteString {
		return fmt.Sprintf("ByteString - Length: [%d]; Value: [%x];", b.length, b.V)

	}
	return fmt.Sprintf("TextString - Length: [%d]; Value: [%s];", b.length, string(b.V))
}

// encodeAsCBORChunks returns a chunk payload (given maxChunkLength) using the indefinite additional tag
func (b baseByteString) encodeAsCBORChunks(maxChunkLength uint8) []byte {

	if maxChunkLength > defaultMaxChunkLength {
		log.Error("Chunking can only support 31 max bytes")
		return nil
	}

	result := []byte{b.majorType.EncodeCBOR() | additionalTypeIndefinite}
	counter := uint64(0)
	for counter < b.length {
		if counter+uint64(maxChunkLength) > b.length {
			chunkLength := b.length - counter
			result = append(result, b.majorType.EncodeCBOR()|byte(chunkLength))
			result = append(result, b.V[counter:]...)
			counter += chunkLength
		} else {
			result = append(result, b.majorType.EncodeCBOR()|maxChunkLength)
			result = append(result, b.V[counter:counter+uint64(maxChunkLength)]...)
			counter += uint64(maxChunkLength)
		}
	}

	result = append(result, indefiniteBreakCode)

	return result
}
