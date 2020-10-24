package multiplex

import (
	"encoding/binary"
	"fmt"

	"github.com/gocardano/go-cardano-client/errors"
	"github.com/gocardano/go-cardano-client/utils"
	log "github.com/sirupsen/logrus"
)

// ContainerMode set to 0 if from initiator, 1 from the responder.
type ContainerMode uint8

const (
	// ContainerModeInitiator indicates that this is from the initiator
	ContainerModeInitiator ContainerMode = 0

	// ContainerModeResponder indicates that this is from the responder
	ContainerModeResponder ContainerMode = 1

	headerSize = 8
)

// Header wraps the ouroboros mux header
type Header struct {
	transmissionTime uint32
	mode             ContainerMode
	miniProtocol     MiniProtocol
	payloadLength    uint16
}

// NewHeader returns an instance of a shelley message header
func NewHeader(miniProtocol MiniProtocol, containerMode ContainerMode,
	payloadLength uint16) *Header {
	return &Header{
		transmissionTime: utils.TimeNowLower32(),
		mode:             containerMode,
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

	if len(buf) != headerSize {
		log.WithFields(log.Fields{
			"expectedLength": headerSize,
			"actualLength":   len(buf),
		}).Error("Invalid message header length")
		return nil, errors.NewError(errors.ErrMuxHeaderInvalidSize)
	}

	containerMode := ContainerMode(buf[4] & 0x80 >> 7)
	miniProtocol := MiniProtocol(binary.BigEndian.Uint16(buf[4:6]) & 0x7fff)

	return &Header{
		transmissionTime: binary.BigEndian.Uint32(buf[0:4]),
		mode:             containerMode,
		miniProtocol:     miniProtocol,
		payloadLength:    binary.BigEndian.Uint16(buf[6:8]),
	}, nil
}

// TransmissionTime returns the transmission time
func (h *Header) TransmissionTime() uint32 {
	return h.transmissionTime
}

// ContainerMode returns the mode of this container (0 from initiator, 1 from responder)
func (h *Header) ContainerMode() ContainerMode {
	return h.mode
}

// MiniProtocolID returns the mini protocol ID
func (h *Header) MiniProtocolID() uint16 {
	return uint16(h.miniProtocol)
}

// PayloadLength returns the payload length
func (h *Header) PayloadLength() uint16 {
	return h.payloadLength
}

// IsFromInitiator return boolean indicating if this container is from initiator
func (h *Header) IsFromInitiator() bool {
	return h.mode == ContainerModeInitiator
}

// IsFromResponder return boolean indicating if this container is from initiator
func (h *Header) IsFromResponder() bool {
	return h.mode == ContainerModeResponder
}

// String description of this message header
func (h *Header) String() string {
	return fmt.Sprintf("Transmission Time: [%d], Mode: [%d], Protocol ID: [%d], Payload Length: [%d]",
		h.TransmissionTime(),
		h.ContainerMode(),
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
	if h.mode == ContainerModeResponder {
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
