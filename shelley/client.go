package shelley

import (
	"fmt"
	"io"

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

	messageResponse, err := c.queryNode(multiplex.MiniProtocolIDMuxControl, handshakeRequest())
	if err != nil {
		log.WithError(err).Error("Error querying node")
	}

	response, err := parseHandshakeResponse(messageResponse)
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

	// Step 1: Send the chain sync request object
	chainSyncRequest := cbor.NewArray()
	chainSyncRequest.Add(cbor.NewPositiveInteger(0)) // msgRequestNext
	messageResponse, err := c.queryNode(multiplex.MiniProtocolIDChainSyncBlocks, chainSyncRequest)
	if err != nil {
		log.WithError(err).Error("Error parsing block fetch response from node")
	}

	// Step 2: Parse the response
	// Response Format:
	// Array: [3]
	//   PositiveInteger8(3)
	//   Array: [0]
	//   Array: [2]
	//     Array: [2]
	// 	     PositiveInteger32(11918355)  // slot
	// 	     ByteString - Length: [32]; Value: [95a417047d3660f2dbd0d70f21b46d7348e9dd0b0e0156ca368cca2d54bcb61b];
	//     PositiveInteger32(4857537)     // blockNumber

	arr := messageResponse.DataItems()[0].(*cbor.Array)
	// arr.Get(0).(*cbor.PositiveInteger8).ValueAsUint8()) // ignore: blockFetch "3"
	// arr.Get(1).(*cbor.Array))                           // ignore: Empty Array
	slotNumber := arr.Get(2).(*cbor.Array).Get(0).(*cbor.Array).Get(0).(*cbor.PositiveInteger32).ValueAsUint32()
	hash := arr.Get(2).(*cbor.Array).Get(0).(*cbor.Array).Get(1).(*cbor.ByteString).ValueAsBytes()
	blockNumber := arr.Get(2).(*cbor.Array).Get(1).(*cbor.PositiveInteger32).ValueAsUint32()

	// Step 3: Send the chainSyncMessageDone to terminate
	chainSyncDone := cbor.NewArray()
	chainSyncDone.Add(cbor.NewPositiveInteger(7)) // chainSyncMessageDone
	_, err = c.queryNode(multiplex.MiniProtocolIDChainSyncBlocks, chainSyncDone)
	if err != nil {
		log.WithError(err).Error("Unexpected error received while terminating with chainSyncMessageDone")
	}

	return slotNumber, hash, blockNumber, nil
}

// queryNode query the node given the input parameters
func (c *Client) queryNode(miniProtocol multiplex.MiniProtocol, input *cbor.Array) (*multiplex.Message, error) {

	// Step 1: Create message for the request
	messageRequest := multiplex.NewMessage(miniProtocol, multiplex.MessageModeInitiator, input)
	log.Trace("Multiplexed Request")
	if log.IsLevelEnabled(log.TraceLevel) {
		fmt.Println(messageRequest.Debug())
	}

	// Step 2: Transmit the request via socket
	messageResponse, err := c.socket.Write(messageRequest.Bytes())
	if err != nil && err != io.EOF {
		log.WithError(err).Error("Error writing to socket")
		return nil, err
	}
	log.Trace("Multiplexed Response")
	if log.IsLevelEnabled(log.TraceLevel) && messageResponse != nil {
		fmt.Println(messageResponse.Debug())
	}

	return messageResponse, nil
}

// Disconnect client from socket
func (c *Client) Disconnect() error {
	return c.socket.Close()
}
