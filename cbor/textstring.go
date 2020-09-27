package cbor

// TextString represents a text string item
type TextString struct {
	baseByteString
}

// NewTextString returns a text string instance
func NewTextString(value string) *TextString {
	return &TextString{
		baseByteString: newBaseByteString(MajorTypeTextString, []byte(value)),
	}
}
