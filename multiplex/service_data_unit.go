package multiplex

// Fixed-size ServiceDataUnits (SDUs) can group multiple messages in a single SDU.
// Reference: https://roadmap.cardano.org/en/status-updates/update/2020-01-24/

import (
	"github.com/gocardano/go-cardano-client/cbor"
	"github.com/gocardano/go-cardano-client/errors"

	log "github.com/sirupsen/logrus"
)

const (
	// MaxSDUSize is the maximum number of bytes for an SDU
	MaxSDUSize = 12288
)

// ServiceDataUnit can group multiple messages
type ServiceDataUnit struct {
	miniProtocol MiniProtocol
	messageMode  MessageMode
	dataItems    []cbor.DataItem
}

// NewServiceDataUnit returns a new message
func NewServiceDataUnit(miniProtocol MiniProtocol, messageMode MessageMode, dataItems []cbor.DataItem) *ServiceDataUnit {
	return &ServiceDataUnit{
		miniProtocol: miniProtocol,
		messageMode:  messageMode,
		dataItems:    dataItems,
	}
}

// ParseServiceDataUnits returns a list of parsed SDU.  Assumes first 8 bytes are the header bits.
func ParseServiceDataUnits(data []byte) ([]*ServiceDataUnit, error) {

	sdus := []*ServiceDataUnit{}

	if len(data) < HeaderSize {
		log.WithField("messageLength", len(data)).Error("Data length below header minimum size")
		return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
	}

	counter := 0
	tmpCborBytes := []byte{}

	for counter < len(data) {

		// Parse header
		header, err := ParseHeader(data[counter : counter+HeaderSize])
		if err != nil {
			log.WithError(err).Error("Error parsing header")
			return sdus, err
		}

		if len(data) < counter+HeaderSize+header.PayloadLengthAsInt32() {
			log.WithFields(log.Fields{
				"actualLength":   len(data),
				"expectedLength": counter + HeaderSize + header.PayloadLengthAsInt32(),
			}).Error("Data length below expected size")
			return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
		}
		tmpCborBytes = append(tmpCborBytes, data[counter+HeaderSize:counter+HeaderSize+header.PayloadLengthAsInt32()]...)

		if header.PayloadLength() == MaxSDUSize {
			log.Tracef("Found payload of MaxSDUSize, aggregating the payload")
		} else {
			dataItems, err := cbor.Decode(tmpCborBytes)
			if err != nil {
				log.WithError(err).Error("Error decoding cbor bytes")
				return nil, err
			}
			sdu := NewServiceDataUnit(header.MiniProtocol(), header.MessageMode(), dataItems)
			sdus = append(sdus, sdu)
			tmpCborBytes = tmpCborBytes[:0]
		}

		counter += HeaderSize + header.PayloadLengthAsInt32()
	}
	return sdus, nil
}

// DataItems returns the dataItems associated with this message
func (s *ServiceDataUnit) DataItems() []cbor.DataItem {
	return s.dataItems
}

// Bytes returns the byte array for this SDU wrapped in multiplexed message format
func (s *ServiceDataUnit) Bytes() []byte {

	buf := []byte{}

	cborPayload := cbor.EncodeList(s.dataItems)
	cborPayloadLength := len(cborPayload)

	for counter := 0; counter < MaxSDUSize; counter += MaxSDUSize {
		if counter+MaxSDUSize > cborPayloadLength {
			buf = append(buf, NewHeader(s.miniProtocol, s.messageMode, uint16(cborPayloadLength-counter)).Bytes()...)
			buf = append(buf, cborPayload[counter:cborPayloadLength]...)
		} else {
			buf = append(buf, NewHeader(s.miniProtocol, s.messageMode, MaxSDUSize).Bytes()...)
			buf = append(buf, cborPayload[counter:counter+MaxSDUSize]...)
		}
	}

	return buf
}

// Debug returns a string representation of the message
func (s *ServiceDataUnit) Debug() string {

	r := cbor.DebugList(s.dataItems)
	r += "============================================================================================"

	return r
}
