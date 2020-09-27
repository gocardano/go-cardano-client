package multiplex

import (
	"math"
)

// MiniProtocol identifies the protocol of the message transmission
type MiniProtocol uint16

const (
	// MiniProtocolIDMuxControl used for handshake
	MiniProtocolIDMuxControl MiniProtocol = 0

	// MiniProtocolIDDeltaQ available for both NtN and NtC
	MiniProtocolIDDeltaQ MiniProtocol = 1

	// MiniProtocolIDChainSyncHeaders available only for NtN (node to node)
	MiniProtocolIDChainSyncHeaders MiniProtocol = 2

	// MiniProtocolIDBlockFetch available only for NtN (node to node)
	MiniProtocolIDBlockFetch MiniProtocol = 3

	// MiniProtocolIDTransactionSubmission available only for NtC (node to client)
	MiniProtocolIDTransactionSubmission MiniProtocol = 4

	// MiniProtocolIDChainSyncBlocks available only for NtC (node to client)
	MiniProtocolIDChainSyncBlocks MiniProtocol = 5

	// MiniProtocolUnknown unknown protocol
	MiniProtocolUnknown MiniProtocol = math.MaxUint16
)

var miniProtocols = map[uint16]MiniProtocol{
	uint16(MiniProtocolIDMuxControl):            MiniProtocolIDMuxControl,
	uint16(MiniProtocolIDDeltaQ):                MiniProtocolIDDeltaQ,
	uint16(MiniProtocolIDChainSyncHeaders):      MiniProtocolIDChainSyncHeaders,
	uint16(MiniProtocolIDBlockFetch):            MiniProtocolIDBlockFetch,
	uint16(MiniProtocolIDTransactionSubmission): MiniProtocolIDTransactionSubmission,
	uint16(MiniProtocolIDChainSyncBlocks):       MiniProtocolIDChainSyncBlocks,
}

var miniProtocolNames = map[MiniProtocol]string{
	MiniProtocolIDMuxControl:            "muxControl",
	MiniProtocolIDDeltaQ:                "deltaQ",
	MiniProtocolIDChainSyncHeaders:      "chainSyncHeaders",
	MiniProtocolIDBlockFetch:            "blockFetch",
	MiniProtocolIDTransactionSubmission: "transactionSubmission",
	MiniProtocolIDChainSyncBlocks:       "chainSyncBlocks",
}

// miniProtocolFromBytes return the mini protocol given the value
func miniProtocolFromBytes(value uint16) MiniProtocol {
	result, ok := miniProtocols[value]
	if !ok {
		return MiniProtocolUnknown
	}
	return result
}

// Value of this mini protocol
func (m MiniProtocol) Value() uint16 {
	return uint16(m)
}

// String representation of this mini protocol
func (m MiniProtocol) String() string {
	result, ok := miniProtocolNames[m]
	if !ok {
		return "unknown"
	}
	return result
}
