package multiplex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {

	h := NewHeader(miniProtocolFromBytes(0x01), ContainerModeInitiator, 2)
	assert.True(t, h.transmissionTime > 0)
	assert.Equal(t, uint16(1), h.miniProtocol.Value())
	assert.Equal(t, ContainerModeInitiator, h.ContainerMode())
	assert.Equal(t, uint16(2), h.payloadLength)
	assert.Equal(t, h.IsFromInitiator(), true)
	assert.Equal(t, h.IsFromResponder(), false)
	assert.True(t, len(h.String()) > 0)

	h = NewHeader(miniProtocolFromBytes(0x01), ContainerModeResponder, 2)
	assert.True(t, h.transmissionTime > 0)
	assert.Equal(t, uint16(1), h.miniProtocol.Value())
	assert.Equal(t, ContainerModeResponder, h.ContainerMode())
	assert.Equal(t, uint16(2), h.payloadLength)
	assert.Equal(t, h.IsFromInitiator(), false)
	assert.Equal(t, h.IsFromResponder(), true)
	assert.True(t, len(h.String()) > 0)

	// Test update payload length
	h.update(uint16(88))
	assert.Equal(t, uint16(88), h.payloadLength)
}

func TestNewMessageHeaderFromBytes(t *testing.T) {

	buf := []byte{
		0x54, 0x95, 0x8a, 0x41, // lower 32 bits of the sender's monotonic clock
		0x00, 0x00, // messageModeInitiator && protocol ID
		0x00, 0x19, // payload length
	}

	// Test 1: Verify Parsing
	header, err := ParseHeader(buf)
	assert.Nil(t, err)
	assert.NotNil(t, header)
	assert.Equal(t, uint32(1419086401), header.transmissionTime)
	assert.Equal(t, ContainerModeInitiator, header.ContainerMode())
	assert.Equal(t, uint16(0), header.miniProtocol.Value())
	assert.Equal(t, uint16(25), header.payloadLength)

	buf2 := header.Bytes()
	assert.Equal(t, buf, buf2)
}

func TestMessageMode(t *testing.T) {

	bufInitiator := []byte{
		0x00, 0x00, 0x00, 0x00, // timestamp
		0x00, 0x00, // messageModeInitiator && protocol ID
		0x00, 0x00, // payload length
	}

	header, err := ParseHeader(bufInitiator)
	assert.Nil(t, err)
	assert.NotNil(t, header)
	assert.Equal(t, ContainerModeInitiator, header.ContainerMode())

	bufResponder := []byte{
		0x00, 0x00, 0x00, 0x00, // timestamp
		0x80, 0x00, // messageModeResponder && protocol ID
		0x00, 0x00, // payload length
	}

	header, err = ParseHeader(bufResponder)
	assert.Nil(t, err)
	assert.NotNil(t, header)
	assert.Equal(t, ContainerModeResponder, header.ContainerMode())
}
