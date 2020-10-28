package multiplex

import "fmt"

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

const (
	newline = "\n"
)

// Message represents an envelope message to/from a Shelley node
type Message struct {
	header *Header
	data   []byte
}

// NewMessage returns a new message
func NewMessage(miniProtocol MiniProtocol, messageMode MessageMode, data []byte) *Message {
	return &Message{
		header: NewHeader(miniProtocol, messageMode, 0),
		data:   data,
	}
}

// Header of this message
func (m *Message) Header() *Header {
	return m.header
}

// Data returns the actual CBOR encoded data in this message
func (m *Message) Data() []byte {
	return m.data
}

// Bytes returns the byte array for this message to be sent on the wire.  It is prefixed with the message header.
func (m *Message) Bytes() []byte {
	return append(m.header.Bytes(), m.data...)
}

// Debug returns a string metadata for the message
func (m *Message) Debug() string {
	return fmt.Sprintf("%s, Actual Length: [%d]", m.header.String(), len(m.data))
}
