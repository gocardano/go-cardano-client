package cbor

import (
	"fmt"
	"math"
	"sort"
)

// Map wraps a CBOR map
type Map struct {
	baseDataItem
	m map[DataItem]DataItem
}

// NewMap returns a new map instance
func NewMap() *Map {
	return &Map{
		baseDataItem: baseDataItem{
			majorType: MajorTypeMap,
		},
		m: map[DataItem]DataItem{},
	}
}

// AdditionalTypeValue returns the length of the byte string
func (m *Map) AdditionalTypeValue() uint64 {
	return uint64(len(m.m))
}

// Add a key/vaklue pair to the map
func (m *Map) Add(key, value DataItem) {
	m.m[key] = value
}

// Get returns the value for the given key
func (m *Map) Get(key DataItem) (DataItem, bool) {
	value, ok := m.m[key]
	return value, ok
}

// Length returns the number of map entries
func (m *Map) Length() int {
	return int(m.AdditionalTypeValue())
}

// Value returns the map
func (m *Map) Value() interface{} {
	return m.m
}

// EncodeCBOR returns CBOR representation for this item
func (m *Map) EncodeCBOR() []byte {
	if uint64(m.Length()) > math.MaxUint64 {
		return m.doEncodeCBOR(false)
	}
	return m.doEncodeCBOR(true)
}

// ValueAsMap returns the value as a map
func (m *Map) ValueAsMap() map[DataItem]DataItem {
	return m.m
}

// String returns description of this item
func (m *Map) String() string {
	return fmt.Sprintf("Map - Items: [%d]", len(m.m))
}

// doEncodeCBOR returns CBOR representation for this map.
// If fixedLength is true, indicate the length of the map as the additional type.
// If fixedLength is false, indicate with 0x1f as the additional type, and suffix with 0xff to indicate break code.
func (m *Map) doEncodeCBOR(fixedLength bool) []byte {

	var result []byte

	if fixedLength {
		// use length of map as additional type
		result = dataItemPrefix(MajorTypeMap, uint64(m.Length()))
	} else {
		// use indefinite as additional type
		result = []byte{MajorTypeMap.EncodeCBOR() | additionalTypeIndefinite}
	}

	keys := []DataItem{}
	for key := range m.m {
		keys = append(keys, key)
	}
	keySorter := NewDataItemSorter(keys)
	sort.Sort(keySorter)

	for _, key := range keys {
		result = append(result, key.(DataItem).EncodeCBOR()...)
		result = append(result, m.m[key].(DataItem).EncodeCBOR()...)
	}

	if !fixedLength {
		// for indefinite format, terminate with the break code
		result = append(result, indefiniteBreakCode)
	}

	return result
}
