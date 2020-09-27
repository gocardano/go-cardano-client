package cbor

import (
	"encoding/base64"

	log "github.com/sirupsen/logrus"
)

type baseSemantic struct {
	baseDataItem
	additionalTypeValue uint64
}

// baseSemanticUTF8String wraps base attributes for semantic UTF8 string
type baseSemanticUTF8String struct {
	baseSemantic
	V string
}

// baseSemanticUTF8Base64 wraps base attributes for semantic utf8 string encoded in base64
type baseSemanticUTF8Base64 struct {
	baseSemantic
	encoded string
	decoded string
}

////////////////////////////////////////////////////////////////////////////////

// AdditionalTypeValue returns the semantic tag ID as uint64
func (b baseSemantic) AdditionalTypeValue() uint64 {
	return b.additionalTypeValue
}

// newBaseSemanticUTF8String returns a base semantic UTF8 string instance
func newBaseSemanticUTF8String(str string, additionalTypeValue uint64) *baseSemanticUTF8String {
	return &baseSemanticUTF8String{
		baseSemantic: baseSemantic{
			baseDataItem: baseDataItem{
				majorType: MajorTypeSemantic,
			},
			additionalTypeValue: additionalTypeValue,
		},
		V: str,
	}
}

// Value returns the UTF8 string
func (b *baseSemanticUTF8String) Value() interface{} {
	return b.V
}

// EncodeCBOR returns CBOR representation for this item
func (b *baseSemanticUTF8String) EncodeCBOR() []byte {
	return append(
		dataItemPrefix(MajorTypeSemantic, b.additionalTypeValue),
		NewTextString(b.V).EncodeCBOR()...)
}

////////////////////////////////////////////////////////////////////////////////

// newBaseSemanticUTF8String returns a base semantic UTF8 encoded in base64 instance
func newBaseSemanticUTF8Base64(encoded string, additionalTypeValue uint64) *baseSemanticUTF8Base64 {
	decodedBuf, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Errorf("Error decoding base64 encoded string [%s] due to %s", encoded, err)
		return nil
	}
	return &baseSemanticUTF8Base64{
		baseSemantic: baseSemantic{
			baseDataItem: baseDataItem{
				majorType: MajorTypeSemantic,
			},
			additionalTypeValue: additionalTypeValue,
		},
		encoded: encoded,
		decoded: string(decodedBuf),
	}
}

// Value returns the UTF8 string encoded in base64
func (b *baseSemanticUTF8Base64) Value() interface{} {
	return b.encoded
}

// EncodeCBOR returns CBOR representation for this item
func (b *baseSemanticUTF8Base64) EncodeCBOR() []byte {
	return append(
		dataItemPrefix(MajorTypeSemantic, b.additionalTypeValue),
		NewTextString(b.encoded).EncodeCBOR()...)
}

// Decode returns the actual value base64-decoded
func (b *baseSemanticUTF8Base64) Decode() string {
	return b.decoded
}
