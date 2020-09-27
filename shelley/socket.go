package shelley

import (
	"net"
	"time"

	"github.com/masterjk/go-cardano-client/errors"
	"github.com/masterjk/go-cardano-client/multiplex"
	"github.com/masterjk/go-cardano-client/utils"

	log "github.com/sirupsen/logrus"
)

const (
	receivePacketSize = 512
)

// UnixSocket wraps a unix socket connection
type UnixSocket struct {
	filename       string
	connection     net.Conn
	readTimeoutMs  int
	writeTimeoutMs int
}

// NewUnixSocket returns a new instance of socket
func NewUnixSocket(filename string, readTimeoutMs, writeTimeoutMs int) (*UnixSocket, error) {

	if !utils.FileExists(filename) {
		return nil, errors.NewMessageErrorf(errors.ErrSocketNotExists,
			"Socket [%s] not found", filename)
	}

	connection, err := net.Dial("unix", filename)
	if err != nil {
		return nil, err
	}

	s := &UnixSocket{
		connection:     connection,
		filename:       filename,
		readTimeoutMs:  readTimeoutMs,
		writeTimeoutMs: writeTimeoutMs,
	}
	return s, nil
}

// Close the socket connection
func (s *UnixSocket) Close() error {
	return s.connection.Close()
}

// Write the payload to the socket and return the result
func (s *UnixSocket) Write(payload []byte) (*multiplex.Container, error) {

	////////////////////////////////////////////////////////////
	// Step 1: Write to socket
	////////////////////////////////////////////////////////////
	s.connection.SetWriteDeadline(time.Now().Add(time.Duration(s.writeTimeoutMs) * time.Millisecond))
	log.Debugf("Attempting to write %d bytes to socket", len(payload))
	written, err := s.connection.Write(payload)
	if err != nil {
		log.WithError(err).Error("Error writing to socket")
		return nil, errors.NewMessageErrorf(errors.ErrSocketWritingToSocket, "Error writing to socket [%s]", s.filename)
	}
	log.Debugf("Successfully written [%d] bytes to socket", written)

	s.connection.SetReadDeadline(time.Now().Add(time.Duration(s.readTimeoutMs) * time.Millisecond))

	////////////////////////////////////////////////////////////
	// Step 2: Read header: transmission time (4 bytes) + Mini Protocol ID (2 bytes) + Payload Length (2 bytes)
	////////////////////////////////////////////////////////////
	header := make([]byte, 8)
	readCount, err := s.connection.Read(header)
	if err != nil {
		log.WithError(err).Error("Error reading packet header of size 8 bytes")
		return nil, err
	}
	if readCount != 8 {
		log.WithError(errors.NewMessageErrorf(errors.ErrSocketReceivedInvalidHeaderSize, "Expected: [%d] Actual: [%d];", 8, readCount))
		return nil, err
	}

	msgHeader, err := multiplex.ParseHeader(header)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"transmissionTime": msgHeader.TransmissionTime(),
		"protocolID":       msgHeader.MiniProtocolID(),
		"payloadLength":    msgHeader.PayloadLength(),
		"header":           utils.DebugBytes(header),
	}).Debugf("Received response header")

	////////////////////////////////////////////////////////////
	// Step 3: Read packet of size payload length
	////////////////////////////////////////////////////////////
	response := []byte{}
	tmp := make([]byte, receivePacketSize)
	totalReadCount := 0
	for totalReadCount < int(msgHeader.PayloadLength()) {
		readCount, err := s.connection.Read(tmp)
		if err != nil {
			log.WithError(err).Error("Error reading from socket")
			return nil, errors.NewMessageErrorf(errors.ErrSocketReadingFromSocket, "Error reading from socket [%s]", s.filename)
		}
		totalReadCount += int(readCount)
		if readCount > 0 {
			response = append(response, tmp[:readCount]...)
		}
		log.WithFields(log.Fields{
			"readCount":      readCount,
			"totalReadCount": totalReadCount,
		}).Debug("Read packet")
	}

	log.WithFields(log.Fields{
		"response": utils.DebugBytes(response),
	}).Debugf("Successfully read %d bytes from socket", totalReadCount)

	return multiplex.ParseContainerWithHeader(msgHeader, response)
}
