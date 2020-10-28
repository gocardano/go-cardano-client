package multiplex

import (
	"encoding/binary"
	"fmt"

	"github.com/gocardano/go-cardano-client/errors"
	"github.com/gocardano/go-cardano-client/utils"
	log "github.com/sirupsen/logrus"
)

// MessageMode set to 0 if from initiator, 1 from the responder.
type MessageMode uint8

const (
	// MessageModeInitiator indicates that this is from the initiator
	MessageModeInitiator MessageMode = 0

	// MessageModeResponder indicates that this is from the responder
	MessageModeResponder MessageMode = 1

	// HeaderSize of a multiplex message
	HeaderSize = 8
)

// Header wraps the ouroboros mux header
type Header struct {
	transmissionTime uint32
	mode             MessageMode
	miniProtocol     MiniProtocol
	payloadLength    uint16
}

// NewHeader returns an instance of a shelley message header
func NewHeader(miniProtocol MiniProtocol, messageMode MessageMode,
	payloadLength uint16) *Header {
	return &Header{
		transmissionTime: utils.TimeNowLower32(),
		mode:             messageMode,
		miniProtocol:     miniProtocol,
		payloadLength:    payloadLength,
	}
}

// update the header with payload length
func (h *Header) update(payloadLength uint16) {
	h.payloadLength = payloadLength
}

// ParseHeader returns the shelley message header from an 8-byte array
func ParseHeader(buf []byte) (*Header, error) {

	if len(buf) != HeaderSize {
		log.WithFields(log.Fields{
			"expectedLength": HeaderSize,
			"actualLength":   len(buf),
		}).Error("Invalid message header length")
		return nil, errors.NewError(errors.ErrMuxHeaderInvalidSize)
	}

	messageMode := MessageMode(buf[4] & 0x80 >> 7)
	miniProtocol := MiniProtocol(binary.BigEndian.Uint16(buf[4:6]) & 0x7fff)

	return &Header{
		transmissionTime: binary.BigEndian.Uint32(buf[0:4]),
		mode:             messageMode,
		miniProtocol:     miniProtocol,
		payloadLength:    binary.BigEndian.Uint16(buf[6:8]),
	}, nil
}

// TransmissionTime returns the transmission time
func (h *Header) TransmissionTime() uint32 {
	return h.transmissionTime
}

// MessageMode returns the mode of this message (0 from initiator, 1 from responder)
func (h *Header) MessageMode() MessageMode {
	return h.mode
}

// MiniProtocol returns the mini protocol of this header
func (h *Header) MiniProtocol() MiniProtocol {
	return h.miniProtocol
}

// MiniProtocolID returns the mini protocol ID
func (h *Header) MiniProtocolID() uint16 {
	return uint16(h.miniProtocol)
}

// PayloadLength returns the payload length
func (h *Header) PayloadLength() uint16 {
	return h.payloadLength
}

// PayloadLengthAsInt32 returns the payload length as int
func (h *Header) PayloadLengthAsInt32() int {
	return int(h.payloadLength)
}

// IsFromInitiator return boolean indicating if this message is from initiator
func (h *Header) IsFromInitiator() bool {
	return h.mode == MessageModeInitiator
}

// IsFromResponder return boolean indicating if this message is from initiator
func (h *Header) IsFromResponder() bool {
	return h.mode == MessageModeResponder
}

// String description of this message header
func (h *Header) String() string {
	return fmt.Sprintf("Transmission Time: [%d], Mode: [%d], Protocol ID: [%d], Payload Length: [%d]",
		h.TransmissionTime(),
		h.MessageMode(),
		h.MiniProtocolID(),
		h.PayloadLength())
}

// Bytes returns a byte array representation
func (h *Header) Bytes() []byte {

	result := make([]byte, 4)
	tmp := make([]byte, 2)

	// transmission time
	binary.BigEndian.PutUint32(result, h.transmissionTime)

	// protocol ID
	protocolID := h.miniProtocol
	if h.mode == MessageModeResponder {
		// in responder mode, flip the MSB to 1
		protocolID = protocolID | 0x8000
	}
	binary.BigEndian.PutUint16(tmp, uint16(protocolID))
	result = append(result, tmp...)

	// payload length
	binary.BigEndian.PutUint16(tmp, h.payloadLength)
	result = append(result, tmp...)

	return result
}
