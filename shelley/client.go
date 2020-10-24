package shelley

import (
	"fmt"

	"github.com/gocardano/go-cardano-client/cbor"
	"github.com/gocardano/go-cardano-client/errors"
	"github.com/gocardano/go-cardano-client/multiplex"
	log "github.com/sirupsen/logrus"
)

const (
	defaultReadTimeoutMs  = 3000
	defaultWriteTimeoutMs = 3000
)

// Client wraps interaction with the shelley node
type Client struct {
	socket *UnixSocket
}

// NewClient returns a new shelley client instance
func NewClient(socketFilename string) (*Client, error) {
	socket, err := NewUnixSocket(socketFilename, defaultReadTimeoutMs, defaultWriteTimeoutMs)
	if err != nil {
		return nil, err
	}
	return &Client{
		socket: socket,
	}, nil
}

// Handshake negotiation with protocol version
func (c *Client) Handshake() error {

	containerResponse, err := c.queryNode(multiplex.MiniProtocolIDMuxControl, handshakeRequest())
	if err != nil {
		log.WithError(err).Error("Error querying node")
	}

	response, err := parseHandshakeResponse(containerResponse)
	if err != nil {
		log.WithError(err).Error("Error parsing handshake response from node")
	}

	if response.accepted == false {
		log.WithFields(log.Fields{
			"versionNumber": response.versionNumber,
			"extraParams":   response.extraParams,
			"refuseReason":  response.refuseReason,
		}).Debug("Handshake failed")
		return errors.NewMessageErrorf(errors.ErrShellyHandshakeFailed, "Handshake failed due to %s", response.refuseReason)
	}

	log.WithFields(log.Fields{
		"versionNumber": response.versionNumber,
		"extraParams":   response.extraParams,
	}).Debug("Handshake was successful")

	return nil
}

// QueryTip returns the block header hash (slotNumber, string, blockNumber, error)
func (c *Client) QueryTip() (uint32, []byte, uint32, error) {

	chainSyncRequest := cbor.NewArray()
	chainSyncRequest.Add(cbor.NewPositiveInteger(0))

	// Response Format:
	// Array: [3]
	//   PositiveInteger8(3)
	//   Array: [0]
	//   Array: [2]
	//     Array: [2]
	// 	     PositiveInteger32(11918355)  // slot
	// 	     ByteString - Length: [32]; Value: [95a417047d3660f2dbd0d70f21b46d7348e9dd0b0e0156ca368cca2d54bcb61b];
	//     PositiveInteger32(4857537)     // blockNumber

	containerResponse, err := c.queryNode(multiplex.MiniProtocolIDChainSyncBlocks, chainSyncRequest)
	if err != nil {
		log.WithError(err).Error("Error parsing block fetch response from node")
	}

	arr := containerResponse.DataItems()[0].(*cbor.Array)
	// arr.Get(0).(*cbor.PositiveInteger8).ValueAsUint8()) : blockFetch "3"
	// arr.Get(1).(*cbor.Array))                           // Empty Array
	slotNumber := arr.Get(2).(*cbor.Array).Get(0).(*cbor.Array).Get(0).(*cbor.PositiveInteger32).ValueAsUint32()
	hash := arr.Get(2).(*cbor.Array).Get(0).(*cbor.Array).Get(1).(*cbor.ByteString).ValueAsBytes()
	blockNumber := arr.Get(2).(*cbor.Array).Get(1).(*cbor.PositiveInteger32).ValueAsUint32()

	return slotNumber, hash, blockNumber, nil
}

// queryNode query the node given the input parameters
func (c *Client) queryNode(miniProtocol multiplex.MiniProtocol, input *cbor.Array) (*multiplex.Container, error) {

	// Step 1: Create container for the request
	containerRequest := multiplex.NewContainer(miniProtocol, multiplex.ContainerModeInitiator, input)
	log.Trace("Multiplexed Request")
	if log.IsLevelEnabled(log.TraceLevel) {
		fmt.Println(containerRequest.Debug())
	}

	// Step 2: Transmit the request via socket
	containerResponse, err := c.socket.Write(containerRequest.Bytes())
	if err != nil {
		log.WithError(err).Error("Error writing to socket")
		return nil, err
	}
	log.Trace("Multiplexed Response")
	if log.IsLevelEnabled(log.TraceLevel) {
		fmt.Println(containerResponse.Debug())
	}

	return containerResponse, nil
}

// Disconnect client from socket
func (c *Client) Disconnect() error {
	return c.socket.Close()
}
