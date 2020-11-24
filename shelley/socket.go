package shelley

import (
	e "errors"
	"io"
	"net"
	"time"

	"github.com/gocardano/go-cardano-client/errors"
	"github.com/gocardano/go-cardano-client/multiplex"
	"github.com/gocardano/go-cardano-client/utils"

	log "github.com/sirupsen/logrus"
)

const (
	networkUnix = "unix"
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

	connection, err := net.Dial(networkUnix, filename)
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
func (s *UnixSocket) Write(payload []byte) (*multiplex.ServiceDataUnit, error) {

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
	// Step 2: Read till EOF
	////////////////////////////////////////////////////////////
	response := []byte{}
	for {

		////////////////////////////////////////////////////////////
		// Step 2a: Read 8 bytes header to determine payload length
		////////////////////////////////////////////////////////////
		bytesHeader := make([]byte, multiplex.HeaderSize)
		readBytes, err := s.connection.Read(bytesHeader)
		if err != nil {
			if e.Is(err, io.EOF) {
				// nothing to read, no-op
				log.Trace("EOF received on reading for response, nothing to read")
				return nil, nil
			}
			log.WithError(err).Error("Error reading header of 8 bytes")
			return nil, err
		} else if readBytes != multiplex.HeaderSize {
			log.Errorf("Expecting to have read 8 bytes for header, but only read [%d] bytes", readBytes)
			return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
		}

		header, err := multiplex.ParseHeader(bytesHeader)
		if err != nil {
			log.WithError(err).Error("Parsed header is invalid")
			return nil, errors.NewError(errors.ErrShelleyPayloadInvalid)
		}
		response = append(response, header.Bytes()...)

		////////////////////////////////////////////////////////////
		// Step 2b: Reading until entire payload has been read
		////////////////////////////////////////////////////////////

		totalReadBytes := 0
		for totalReadBytes < header.PayloadLengthAsInt32() {
			buf := make([]byte, header.PayloadLengthAsInt32()-totalReadBytes)
			readBytes, err = s.connection.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.WithError(err).Error("Error reading from socket")
				}
				break
			}

			log.WithFields(log.Fields{
				"expectedPayloadLength": header.PayloadLength(),
				"readBytes":             readBytes,
				"totalReadBytes":        totalReadBytes,
			}).Trace("Received CBOR data")

			response = append(response, buf[:readBytes]...)
			totalReadBytes += readBytes
		}

		if int(header.PayloadLength()) != multiplex.MaxSDUSize {
			log.WithField("headerPayloadLength", header.PayloadLength()).
				Trace("Breaking out of loop since read payload is not MaxSDUSize")
			break
		}
	}

	log.WithField("responseLength", len(response)).Debug("Total read response bytes from socket")

	sdus, err := multiplex.ParseServiceDataUnits(response)
	if err != nil {
		return nil, err
	}

	return sdus[0], nil
}
