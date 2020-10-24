package multiplex

// Container represent a message envelope that is sent and receive
// to/from the Shelley node.  It contains one or more segments as payload.
// The following is the wire format of a container:
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
// Container header:
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

// Container represents an envelope message to/from a Shelley node
type Container struct {
	header    *Header
	dataItems []cbor.DataItem
}

// NewContainer returns a new container
func NewContainer(miniProtocol MiniProtocol, containerMode ContainerMode, array *cbor.Array) *Container {
	return &Container{
		header:    NewHeader(miniProtocol, containerMode, 0),
		dataItems: []cbor.DataItem{array},
	}
}

// DataItems returns the dataItems associated with this container
func (c *Container) DataItems() []cbor.DataItem {
	return c.dataItems
}

// Header of this container
func (c *Container) Header() *Header {
	return c.header
}

// Bytes returns the byte array for this container
func (c *Container) Bytes() []byte {
	payload := cbor.EncodeList(c.dataItems)
	payloadLength := len(payload)
	if payloadLength > math.MaxUint16 {
		log.WithFields(log.Fields{
			"payloadLength": payloadLength,
		}).Error("Payload length exceeded maximum limit in the container of math.MaxUint16")
		return nil
	}
	c.header.update(uint16(len(payload)))
	return append(c.header.Bytes(), payload...)
}

// ParseContainerWithHeader uses the header and parses the byte array and return the container
func ParseContainerWithHeader(header *Header, payload []byte) (*Container, error) {

	if int(header.PayloadLength()) != len(payload) {
		log.WithFields(log.Fields{
			"expectPayloadLength": header.PayloadLength(),
			"actualPayloadLength": len(payload),
		}).Error("Container header's payload length does not match actual payload length")
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

	log.Infof("Shelley container has [%d] CBOR encoded data items", len(dataItems))

	return &Container{
		header:    header,
		dataItems: dataItems,
	}, nil
}

// ParseContainer parses the byte array and return the container
func ParseContainer(buf []byte) (*Container, error) {
	if len(buf) < headerSize {
		log.WithFields(log.Fields{
			"messageLength": len(buf),
		}).Error("Message length below header minimum size")
		return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
	}

	header, err := ParseHeader(buf[:headerSize])
	if err != nil {
		log.WithFields(log.Fields{
			"header": utils.DebugBytes(buf[:headerSize]),
			"error":  err.Error(),
		}).Error("Error parsing shelley payload header")
		return nil, err
	}

	return ParseContainerWithHeader(header, buf[8:])
}

// Debug returns a string representation of the messsage
func (c *Container) Debug() string {

	r := "==========================================================================================" + newline
	r += "Header: " + c.header.String() + newline
	r += "------------------------------------------------------------------------------------------" + newline
	r += cbor.DebugList(c.dataItems)
	r += "=========================================================================================="

	return r
}
