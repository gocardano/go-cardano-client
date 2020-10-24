package multiplex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiniProtocolValues(t *testing.T) {
	for _, m := range miniProtocols {
		assert.True(t, len(m.String()) > 0)
		assert.True(t, m.Value() >= 0)
	}

	// Test Unknown MiniProtocol
	testMiniProtocol := uint16(9999)
	assert.Equal(t, "unknown", MiniProtocol(testMiniProtocol).String())
}
