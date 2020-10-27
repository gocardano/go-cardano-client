package multiplex

// Message represent an envelope that is sent and receive
// to/from the Shelley node.  It contains one or more segments as payload.
// The following is the wire format of a message:
//
// +---------------------------------------------------------------+
// |0|1|2|3|4|5|6|7|8|9|0|1|2|3|4|5|6|7|8|9|0|1|2|3|4|5|6|7|8|9|0|1|
// +---------------------------------------------------------------+
// |                       TRANSMISSION TIME                       |
// +---------------------------------------------------------------+
// |M|     MINI PROTOCOL ID      |         PAYLOAD LENGTH          |
// +---------------------------------------------------------------+
// |                                                               |
// |                       PAYLOAD of n BYTES                      |
// |                                                               |
// +---------------------------------------------------------------+
//
// Message header:
// - Transmission Time The transmission time is a time stamp based the wall clock
//   of the peer with a resolution of one microsecond.
// - Mini Protocol ID The unique ID of the mini protocol as in Table 3.2.
// - Payload Length The payload length is the size of the segment payload in Bytes.
//   The maximum payload length that is supported by the multiplexing wire format
//   is 2^16 − 1. Note, that an instance of the protocol can choose a smaller
//   limit for the size of segments it transmits.
// - Mode The single bit M (the mode) is used to distinct the dual instances of a
//   mini protocol. The mode is set to 0 in segments from the initiator, i.e. the
//   side that initially has agency and 1 in segments from the responder.
//
// Reference: https://hydra.iohk.io/build/4110312/download/2/network-spec.pdf

import (
	"math"

	"github.com/gocardano/go-cardano-client/cbor"
	"github.com/gocardano/go-cardano-client/errors"
	"github.com/gocardano/go-cardano-client/utils"

	log "github.com/sirupsen/logrus"
)

const (
	newline = "\n"
)

// Message represents an envelope message to/from a Shelley node
type Message struct {
	header    *Header
	dataItems []cbor.DataItem
}

// NewMessage returns a new message
func NewMessage(miniProtocol MiniProtocol, messageMode MessageMode, array *cbor.Array) *Message {
	return &Message{
		header:    NewHeader(miniProtocol, messageMode, 0),
		dataItems: []cbor.DataItem{array},
	}
}

// DataItems returns the dataItems associated with this message
func (m *Message) DataItems() []cbor.DataItem {
	return m.dataItems
}

// Header of this message
func (m *Message) Header() *Header {
	return m.header
}

// Bytes returns the byte array for this message
func (m *Message) Bytes() []byte {
	payload := cbor.EncodeList(m.dataItems)
	payloadLength := len(payload)
	if payloadLength > math.MaxUint16 {
		log.WithFields(log.Fields{
			"payloadLength": payloadLength,
		}).Error("Payload length exceeded maximum limit in the message of math.MaxUint16")
		return nil
	}
	m.header.update(uint16(len(payload)))
	return append(m.header.Bytes(), payload...)
}

// ParseMessageWithHeader uses the header and parses the byte array and return the message
func ParseMessageWithHeader(header *Header, payload []byte) (*Message, error) {

	if int(header.PayloadLength()) != len(payload) {
		log.WithFields(log.Fields{
			"expectPayloadLength": header.PayloadLength(),
			"actualPayloadLength": len(payload),
		}).Error("Message header's payload length does not match actual payload length")
		return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
	}

	dataItems, err := cbor.Decode(payload)
	if err != nil {
		log.WithFields(log.Fields{
			"payload": utils.DebugBytes(payload),
			"error":   err.Error(),
		}).Error("Error parsing shelley payload")
		return nil, err
	}

	log.Infof("Shelley message has [%d] CBOR encoded data items", len(dataItems))

	return &Message{
		header:    header,
		dataItems: dataItems,
	}, nil
}

// ParseMessage parses the byte array and return the message (uses the first 8 bytes as the header)
func ParseMessage(buf []byte) (*Message, error) {
	if len(buf) < HeaderSize {
		log.WithFields(log.Fields{
			"messageLength": len(buf),
		}).Error("Message length below header minimum size")
		return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
	}

	header, err := ParseHeader(buf[:HeaderSize])
	if err != nil {
		log.WithFields(log.Fields{
			"header": utils.DebugBytes(buf[:HeaderSize]),
			"error":  err.Error(),
		}).Error("Error parsing shelley payload header")
		return nil, err
	}

	return ParseMessageWithHeader(header, buf[8:])
}

// Debug returns a string representation of the message
func (m *Message) Debug() string {

	r := "==========================================================================================" + newline
	r += "Header: " + m.header.String() + newline
	r += "------------------------------------------------------------------------------------------" + newline
	r += cbor.DebugList(m.dataItems)
	r += "=========================================================================================="

	return r
}
