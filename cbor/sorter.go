package cbor

import (
	"bytes"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// DataItemSorter provides functions to sort list of dataItems
type DataItemSorter struct {
	dataItems []DataItem
}

// NewDataItemSorter returns data item sorter
func NewDataItemSorter(dataItems []DataItem) *DataItemSorter {
	return &DataItemSorter{
		dataItems: dataItems,
	}
}

// Len returns the size of the array
func (sorter *DataItemSorter) Len() int {
	return len(sorter.dataItems)
}

// Swap exchanges the position of two items in the array
func (sorter *DataItemSorter) Swap(i, j int) {
	sorter.dataItems[i], sorter.dataItems[j] = sorter.dataItems[j], sorter.dataItems[i]
}

// Less returns if the item is "less" than the other
func (sorter *DataItemSorter) Less(i, j int) bool {

	if sorter.dataItems[i].MajorType() != sorter.dataItems[j].MajorType() {
		log.WithFields(log.Fields{
			"i": sorter.dataItems[i].MajorType(),
			"j": sorter.dataItems[j].MajorType(),
		}).Debug("Unable to compare items of different majorType for sorting")
		return false
	}

	if sorter.dataItems[i].MajorType() == MajorTypePositiveInt || sorter.dataItems[j].MajorType() == MajorTypePositiveInt {
		return reflect.ValueOf(sorter.dataItems[i].Value()).Uint() < reflect.ValueOf(sorter.dataItems[j].Value()).Uint()
	} else if sorter.dataItems[i].MajorType() == MajorTypeNegativeInt || sorter.dataItems[j].MajorType() == MajorTypeNegativeInt {
		return reflect.ValueOf(sorter.dataItems[i].Value()).Int() < reflect.ValueOf(sorter.dataItems[j].Value()).Int()
	} else if sorter.dataItems[i].MajorType() == MajorTypeByteString || sorter.dataItems[j].MajorType() == MajorTypeByteString {
		return bytes.Compare(sorter.dataItems[i].(*ByteString).ValueAsBytes(), sorter.dataItems[j].(*ByteString).ValueAsBytes()) == -1
	} else if sorter.dataItems[i].MajorType() == MajorTypeTextString || sorter.dataItems[j].MajorType() == MajorTypeTextString {
		return strings.Compare(sorter.dataItems[i].(*TextString).ValueAsString(), sorter.dataItems[j].(*TextString).ValueAsString()) == -1
	} else {
		log.WithFields(log.Fields{
			"i": sorter.dataItems[i].String(),
			"j": sorter.dataItems[j].String(),
		}).Debug("Unexpected major types not handled for sorting")
	}

	// TBD: need to handle for other major types
	return false
}
