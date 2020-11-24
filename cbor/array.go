package cbor

import (
	"fmt"
	"math"
)

const (
	indefiniteBreakCodeMajorType      uint8 = 0x07
	indefiniteBreakCodeAdditionalType uint8 = 0x1f
)

// Array represents a CBOR array
type Array struct {
	baseDataItem
	V []DataItem
}

// NewArray returns instance of array data items
func NewArray() *Array {
	return NewArrayWithItems([]DataItem{})
}

// NewArrayWithItems returns instance of array data items
func NewArrayWithItems(items []DataItem) *Array {
	return &Array{
		baseDataItem: baseDataItem{
			majorType: MajorTypeArray,
		},
		V: items,
	}
}

// Value returns the array
func (a *Array) Value() interface{} {
	return a.V
}

// ValuesAsString returns the values as a string (for debugging purposes)
func (a *Array) ValuesAsString() string {
	result := ""
	for _, item := range a.V {
		result += item.String() + "; "
	}
	return result
}

// EncodeCBOR returns CBOR representation for this item
func (a *Array) EncodeCBOR() []byte {
	if uint64(a.Length()) > math.MaxUint64 {
		return a.doEncodeCBOR(false)
	}
	return a.doEncodeCBOR(true)
}

// String returns description of this item
func (a *Array) String() string {
	return fmt.Sprintf("Array: [%d]", len(a.V))
}

// Add an item to the list
func (a *Array) Add(item DataItem) {
	a.V = append(a.V, item)
}

// AdditionalTypeValue returns the length of the array
func (a *Array) AdditionalTypeValue() uint64 {
	return uint64(len(a.V))
}

// Length returns the item count in the array
func (a *Array) Length() int {
	return int(a.AdditionalTypeValue())
}

// Get returns the item from the array given the index position
func (a *Array) Get(idx int) DataItem {
	return a.V[idx]
}

// List returns the items in the array
func (a *Array) List() []DataItem {
	return a.V
}

// doEncodeCBOR returns CBOR representation for this array.
// If fixedLength is true, indicate the length of the array as the additional type.
// If fixedLength is false, indicate with 0x1f as the additional type, and suffix with 0xff to indicate break code.
func (a *Array) doEncodeCBOR(fixedLength bool) []byte {

	var result []byte

	if fixedLength {
		// use length of map as additional type
		result = dataItemPrefix(MajorTypeArray, uint64(a.Length()))
	} else {
		// use indefinite as additional type
		result = []byte{MajorTypeArray.EncodeCBOR() | additionalTypeIndefinite}
	}

	for _, item := range a.V {
		result = append(result, item.EncodeCBOR()...)
	}

	if !fixedLength {
		// for indefinite format, terminate with the break code
		result = append(result, indefiniteBreakCode)
	}

	return result
}
