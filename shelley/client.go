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

	client := &Client{
		socket: socket,
	}

	if err := client.handshake(); err != nil {
		return nil, err
	}

	return client, nil
}

// Disconnect client from socket
func (c *Client) Disconnect() error {
	return c.socket.Close()
}

// Reset the socket by disconnecting and reconnecting
func (c *Client) Reset() error {

	if err := c.Disconnect(); err != nil {
		return fmt.Errorf("Error trying to reset the shelley client %s", err)
	}

	socket, err := NewUnixSocket(c.socket.filename, defaultReadTimeoutMs, defaultWriteTimeoutMs)
	if err == nil {
		c.socket = socket
	}

	if err := c.handshake(); err != nil {
		return err
	}

	return err
}

// Handshake negotiation with protocol version
func (c *Client) handshake() error {

	messageResponse, err := c.queryNode(multiplex.MiniProtocolIDMuxControl, handshakeRequest())
	if err != nil {
		log.WithError(err).Error("Error querying node")
	}

	response, err := parseHandshakeResponse(messageResponse)
	if err != nil {
		log.WithError(err).Error("Error parsing handshake response from node")
		return err
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

// StakePools returns list of stake pools
func (c *Client) StakePools(slotNumber uint32, hash []byte) (*multiplex.ServiceDataUnit, error) {

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> R E Q U E S T   #     1 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// MiniProtocol: 7   /   MessageMode: 0
	// Array: [2]
	//   PositiveInteger8(0)
	//   Array: [2]
	//     PositiveInteger32(14592398)
	//     ByteString - Length: [32]; Value: [85575f65630abaab2bab1ab1171ff7baf3554986529877c168cb88cc36d6a945];
	// ============================================================================================
	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< R E S P O N S E   #     1 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	// MiniProtocol: 7 / MessageMode: 1
	// Array: [1]
	//   PositiveInteger8(1)
	// ============================================================================================

	setBlockRequest := cbor.NewArrayWithItems([]cbor.DataItem{
		cbor.NewPositiveInteger8(0),
		cbor.NewArrayWithItems([]cbor.DataItem{
			cbor.NewPositiveInteger32(slotNumber),
			cbor.NewByteString(hash),
		}),
	})
	setBlockResponse, err := c.queryNode(multiplex.MiniProtocolIDLocalStateQuery, []cbor.DataItem{setBlockRequest})
	if err != nil {
		log.WithError(err).Error("Error parsing block fetch response from node")
		return nil, err
	}
	if setBlockResponse == nil {
		return nil, fmt.Errorf("setBlockResponse was nil")
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> R E Q U E S T   #     2 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// MiniProtocol: 7   /   MessageMode: 0
	// Array: [2]
	//   PositiveInteger8(3)
	//   Array: [2]
	// 	   PositiveInteger8(0)
	// 	   Array: [2]
	// 	     PositiveInteger8(1)
	// 	     Array: [1]
	// 		   PositiveInteger8(5)
	// ============================================================================================
	stakePoolRequest := cbor.NewArrayWithItems([]cbor.DataItem{
		cbor.NewPositiveInteger8(3),
		cbor.NewArrayWithItems([]cbor.DataItem{
			cbor.NewPositiveInteger8(0),
			cbor.NewArrayWithItems([]cbor.DataItem{
				cbor.NewPositiveInteger8(1),
				cbor.NewArrayWithItems([]cbor.DataItem{
					cbor.NewPositiveInteger8(5),
				}),
			}),
		}),
	})
	sduStakePools, err := c.queryNode(multiplex.MiniProtocolIDLocalStateQuery, []cbor.DataItem{stakePoolRequest})
	if err != nil {
		log.WithError(err).Error("Error parsing block fetch response from node")
		return nil, err
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> R E Q U E S T   #     3 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// MiniProtocol: 7   /   MessageMode: 0
	// Array: [1]
	//   PositiveInteger8(5)
	// Array: [1]
	//   PositiveInteger8(7)
	// ============================================================================================

	terminateRequest := []cbor.DataItem{
		cbor.NewArrayWithItems([]cbor.DataItem{cbor.NewPositiveInteger8(5)}),
		cbor.NewArrayWithItems([]cbor.DataItem{cbor.NewPositiveInteger8(7)}),
	}
	if _, err = c.queryNode(multiplex.MiniProtocolIDLocalStateQuery, terminateRequest); err != nil {
		return nil, err
	}

	return sduStakePools, nil
}

// QueryTip returns the block header hash (slotNumber, string, blockNumber, error)
func (c *Client) QueryTip() (uint32, []byte, uint32, error) {

	// Step 1: Send the chain sync request object
	log.Debug("Sending command: msgRequestNext")
	chainSyncRequest := cbor.NewArray()
	chainSyncRequest.Add(cbor.NewPositiveInteger(0)) // msgRequestNext
	messageResponse, err := c.queryNode(multiplex.MiniProtocolIDChainSyncBlocks, []cbor.DataItem{chainSyncRequest})
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
	log.Debug("Sending command: chainSyncMessageDone")
	chainSyncDone := cbor.NewArray()
	chainSyncDone.Add(cbor.NewPositiveInteger(7)) // chainSyncMessageDone
	_, err = c.queryNode(multiplex.MiniProtocolIDChainSyncBlocks, []cbor.DataItem{chainSyncDone})
	if err != nil {
		log.WithError(err).Error("Unexpected error received while terminating with chainSyncMessageDone")
	}

	return slotNumber, hash, blockNumber, nil
}

// queryNode query the node given the input parameters
func (c *Client) queryNode(miniProtocol multiplex.MiniProtocol, dataItems []cbor.DataItem) (*multiplex.ServiceDataUnit, error) {

	// Step 1: Create message for the request
	sdu := multiplex.NewServiceDataUnit(miniProtocol, multiplex.MessageModeInitiator, dataItems)
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Multiplexed Request:")
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> R E Q U E S T >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println(sdu.Debug())
	}

	// Step 2: Send the request
	messageResponse, err := c.socket.Write(sdu.Bytes())
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Error writing to socket %w", err)
	}
	if log.IsLevelEnabled(log.DebugLevel) && messageResponse != nil {
		log.Debug("Multiplexed Response:")
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< R E S P O N S E <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		fmt.Println(messageResponse.Debug())
	}

	return messageResponse, nil
}
