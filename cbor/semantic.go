package cbor

import (
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	rfc3339             = "2006-01-02T15:04:05Z" // diff from time.RFC3339
	rfc3339StringLength = len(rfc3339)
)

var bigZero *big.Int = big.NewInt(0)

// DateTimeString wraps a date/time string (semantic additional type value: 0)
type DateTimeString struct {
	baseSemantic
	V time.Time
}

// DateTimeEpoch wraps a date/time epoch value (semantic additional type value: 1)
type DateTimeEpoch struct {
	baseSemantic
	V time.Time
}

// PositiveBignum wraps a big integer - section 2.4.2 of RFC7049 (semantic additional type value: 2)
type PositiveBignum struct {
	baseSemantic
	V *big.Int
}

// NegativeBignum wraps a big negative integer - section 2.4.2 of RFC7049 (semantic additional type value: 3)
type NegativeBignum struct {
	baseSemantic
	V *big.Int
}

// URI wraps a URI string (semantic additional type value: 32)
type URI struct {
	*baseSemanticUTF8String
}

// Base64URL wraps a URI string encoded in base64 (semantic additional type value: 33)
type Base64URL struct {
	*baseSemanticUTF8Base64
}

// Base64String wraps a string encoded in base64 (semantic additional type value: 34)
type Base64String struct {
	*baseSemanticUTF8Base64
}

// RegularExpression wraps a regular expression (semantic additional type value: 35)
type RegularExpression struct {
	*baseSemanticUTF8String
}

// MimeMessage wraps a mime message (semantic additional type value: 36)
type MimeMessage struct {
	*baseSemanticUTF8String
}

////////////////////////////////////////////////////////////////////////////////

// NewDateTimeString returns a date time string instance (additional type: 0)
func NewDateTimeString(str string) *DateTimeString {
	t, err := time.Parse(rfc3339, str)
	if err != nil {
		log.Errorf("Error parsing RFC3339 for semantic tag date string: [%s] due to [%s]", str, err)
		return nil
	}
	return &DateTimeString{
		baseSemantic: baseSemantic{
			baseDataItem: baseDataItem{
				majorType: MajorTypeSemantic,
			},
			additionalTypeValue: semanticDateTimeString,
		},
		V: t,
	}
}

// Value returns the date/time string
func (d *DateTimeString) Value() interface{} {
	return d.V
}

// EncodeCBOR returns CBOR representation for this item
func (d *DateTimeString) EncodeCBOR() []byte {
	textString := d.V.Format(time.RFC3339)

	// semantic tag prefix
	buf := dataItemPrefix(MajorTypeSemantic, semanticDateTimeString)

	// text tag prefix
	buf = append(buf, dataItemPrefix(MajorTypeTextString, uint64(len(textString)))...)

	// rfc 3339 text
	buf = append(buf, []byte(textString)...)

	return buf
}

// String returns description of this item
func (d *DateTimeString) String() string {
	return fmt.Sprintf("DateTimeString - Value: [%s]", d.V.Format(time.RFC3339))
}

////////////////////////////////////////////////////////////////////////////////

// NewDateTimeEpoch returns a date time epoch instance (additional type: 1)
func NewDateTimeEpoch(epoch int64) *DateTimeEpoch {
	t := time.Unix(epoch, 0)
	return &DateTimeEpoch{
		baseSemantic: baseSemantic{
			baseDataItem: baseDataItem{
				majorType: MajorTypeSemantic,
			},
			additionalTypeValue: semanticDateTimeEpoch,
		},
		V: t,
	}
}

// Value returns the date/time epoch value
func (d *DateTimeEpoch) Value() interface{} {
	return d.V
}

// EncodeCBOR returns CBOR representation for this item
func (d *DateTimeEpoch) EncodeCBOR() []byte {

	// semantic tag prefix
	buf := dataItemPrefix(MajorTypeSemantic, semanticDateTimeEpoch)

	// text tag prefix
	epoch := d.V.Unix()
	if epoch > 0 {
		buf = append(buf, NewPositiveInteger64(uint64(epoch)).EncodeCBOR()...)
	} else {
		buf = append(buf, NewNegativeInteger64(epoch).EncodeCBOR()...)
	}

	return buf
}

// String returns description of this item
func (d *DateTimeEpoch) String() string {
	return fmt.Sprintf("DateTimeEpoch - Value: [%d]", d.V.Unix())
}

////////////////////////////////////////////////////////////////////////////////

// NewPositiveBignumber returns a new positive bignum (additional type: 2)
func NewPositiveBignumber(n *big.Int) *PositiveBignum {
	if n.Cmp(bigZero) == -1 {
		log.Error("Positive bignum should be positive number")
		return nil
	}
	return &PositiveBignum{
		baseSemantic: baseSemantic{
			baseDataItem: baseDataItem{
				majorType: MajorTypeSemantic,
			},
			additionalTypeValue: semanticPositiveBignum,
		},
		V: n,
	}
}

// Value returns the positive big integer
func (n *PositiveBignum) Value() interface{} {
	return n.V
}

// EncodeCBOR returns CBOR representation for this item
func (n *PositiveBignum) EncodeCBOR() []byte {
	result := dataItemPrefix(MajorTypeSemantic, semanticPositiveBignum)
	result = append(result, NewByteString(n.V.Bytes()).EncodeCBOR()...)
	return result
}

// String returns description of this item
func (n *PositiveBignum) String() string {
	return fmt.Sprintf("PositiveBigNum - Value: [%s]", n.V.String())
}

////////////////////////////////////////////////////////////////////////////////

// NewNegativeBignumber returns a new negative bignum (additional type: 3)
func NewNegativeBignumber(n *big.Int) *NegativeBignum {
	if n.Cmp(bigZero) >= 0 {
		log.Error("Negative bignum should be negative number")
		return nil
	}
	return &NegativeBignum{
		baseSemantic: baseSemantic{
			baseDataItem: baseDataItem{
				majorType: MajorTypeSemantic,
			},
			additionalTypeValue: semanticNegativeBignum,
		},
		V: n,
	}
}

// Value returns the negative big integer
func (n *NegativeBignum) Value() interface{} {
	return n.V
}

// EncodeCBOR returns CBOR representation for this item
func (n *NegativeBignum) EncodeCBOR() []byte {
	result := dataItemPrefix(MajorTypeSemantic, semanticNegativeBignum)
	result = append(result, NewByteString(n.encodedValue().Bytes()).EncodeCBOR()...)
	return result
}

// encodedValue is the value represented as big.Int when transmitted via CBOR (ie. -1 - value)
func (n *NegativeBignum) encodedValue() *big.Int {
	t := big.NewInt(-1)
	t = t.Sub(t, n.V)
	return t
}

// String returns description of this item
func (n *NegativeBignum) String() string {
	return fmt.Sprintf("NegativeBignum - Value: [%s]", n.V.String())
}

////////////////////////////////////////////////////////////////////////////////

// NewURI returns a URI instance (additional type: 32)
func NewURI(str string) *URI {
	return &URI{
		baseSemanticUTF8String: newBaseSemanticUTF8String(str, semanticURI),
	}
}

// String returns description of this item
func (u *URI) String() string {
	return fmt.Sprintf("URI - Value: [%s]", u.V)
}

////////////////////////////////////////////////////////////////////////////////

// NewBase64URL returns a Base64 URL instance (additional type: 33)
func NewBase64URL(encoded string) *Base64URL {
	return &Base64URL{
		baseSemanticUTF8Base64: newBaseSemanticUTF8Base64(encoded, semanticBase64URL),
	}
}

// String returns description of this item
func (b *Base64URL) String() string {
	return fmt.Sprintf("Base64URL - Value: [%s]; Decoded: [%s];",
		b.baseSemanticUTF8Base64.encoded,
		b.baseSemanticUTF8Base64.decoded)
}

////////////////////////////////////////////////////////////////////////////////

// NewBase64String returns a Base64 encoded string (additional type: 34)
func NewBase64String(encoded string) *Base64String {
	return &Base64String{
		baseSemanticUTF8Base64: newBaseSemanticUTF8Base64(encoded, semanticBase64),
	}
}

// String returns description of this item
func (b *Base64String) String() string {
	return fmt.Sprintf("Base64String - Value: [%s]; Decoded: [%s];",
		b.baseSemanticUTF8Base64.encoded,
		b.baseSemanticUTF8Base64.decoded)
}

////////////////////////////////////////////////////////////////////////////////

// NewRegularExpression returns a regular expression instance (additional type: 35)
func NewRegularExpression(str string) *RegularExpression {
	return &RegularExpression{
		baseSemanticUTF8String: newBaseSemanticUTF8String(str, semanticRegularExpression),
	}
}

// String returns description of this item
func (r *RegularExpression) String() string {
	return fmt.Sprintf("RegularExpression - Value: [%s]", r.V)
}

////////////////////////////////////////////////////////////////////////////////

// NewMimeMessage returns a mime message instance (additional type: 36)
func NewMimeMessage(str string) *MimeMessage {
	return &MimeMessage{
		baseSemanticUTF8String: newBaseSemanticUTF8String(str, semanticMimeMessage),
	}
}

// String returns description of this item
func (m *MimeMessage) String() string {
	return fmt.Sprintf("MimeMessage - Value: [%s]", m.V)
}
