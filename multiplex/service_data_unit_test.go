package multiplex

import (
	"testing"

	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

func TestParseServiceDataUnit(t *testing.T) {

	var validHeader []byte = []byte{
		0x54, 0x95, 0x8a, 0x41, //time
		0x00, 0x00, // protocol ID
		0x00, 0x19, // payload length
	}

	var validPayload []byte = []byte{
		0x82, 0x00, 0xa3, 0x01, 0x1a, 0x2d, 0x96, 0x4a, 0x09, 0x19, 0x80, 0x02, 0x1a, 0x2d, 0x96, 0x4a, 0x09, 0x19, 0x80, 0x03, 0x1a, 0x2d, 0x96, 0x4a, 0x09,
	}

	// Scenario: Valid header/payload
	buf := append(validHeader, validPayload...)
	sdus, err := ParseServiceDataUnits(buf)
	assert.Nil(t, err)
	assert.NotNil(t, sdus)
	assert.Equal(t, 1, len(sdus[0].dataItems))

	// Scenario: Missing
	buf = append(validHeader, 0x01)
	_, err = ParseServiceDataUnits(buf)
	assert.NotNil(t, err)
	log.Error(err)
}
