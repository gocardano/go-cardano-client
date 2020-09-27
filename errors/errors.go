package errors

import (
	"fmt"
)

// Severity of the error
type Severity uint8

const (
	// ERROR level
	ERROR Severity = iota
	// WARNING level
	WARNING
	// INFO level
	INFO
	// DEBUG level
	DEBUG
)

// CLIError wraps an error encountered
type CLIError struct {
	severity Severity
	code     int
	desc     string
	message  string
}

const (
	ErrSocketNotExists                 = 100
	ErrSocketWritingToSocket           = 101
	ErrSocketReadingFromSocket         = 102
	ErrSocketReceivedInvalidHeaderSize = 103

	ErrMuxHeaderInvalidSize = 201

	ErrBitstreamReaderEOF               = 301
	ErrBitstreamVarInsufficientCapacity = 302

	ErrCborMajorTypeUnhandled              = 401
	ErrCborAdditionalTypeUnhandled         = 402
	ErrCborNegativeIntsOnly                = 403
	ErrCborUnhandledReadBytesInTermsOfBits = 404
	ErrCborBignumParsingFailed             = 405

	ErrShelleyPayloadInvalid       = 501
	ErrShelleyInvalidContainerMode = 502
	ErrShellyUnexpectedCborItem    = 503
	ErrShellyHandshakeFailed       = 504
)

var cliErrorMap = map[int]CLIError{
	ErrSocketNotExists: {
		severity: ERROR,
		code:     ErrSocketNotExists,
		desc:     "Unix socket not found",
	},
	ErrSocketWritingToSocket: {
		severity: ERROR,
		code:     ErrSocketWritingToSocket,
		desc:     "Error encountered while writing to socket",
	},
	ErrSocketReadingFromSocket: {
		severity: ERROR,
		code:     ErrSocketReadingFromSocket,
		desc:     "Error encountered while reading from socket",
	},
	ErrSocketReceivedInvalidHeaderSize: {
		severity: ERROR,
		code:     ErrSocketReceivedInvalidHeaderSize,
		desc:     "Received invalid header size from node socket",
	},
	ErrMuxHeaderInvalidSize: {
		severity: ERROR,
		code:     ErrMuxHeaderInvalidSize,
		desc:     "Invalid mux header size",
	},
	ErrBitstreamReaderEOF: {
		severity: ERROR,
		code:     ErrBitstreamReaderEOF,
		desc:     "Reached EOF of the byte stream reader",
	},
	ErrBitstreamVarInsufficientCapacity: {
		severity: ERROR,
		code:     ErrBitstreamVarInsufficientCapacity,
		desc:     "Destination variable has insufficient bit capacity",
	},
	ErrCborMajorTypeUnhandled: {
		severity: ERROR,
		code:     ErrCborMajorTypeUnhandled,
		desc:     "CBOR major type unassigned or unknown",
	},
	ErrCborAdditionalTypeUnhandled: {
		severity: ERROR,
		code:     ErrCborAdditionalTypeUnhandled,
		desc:     "CBOR additional type unhandled",
	},
	ErrCborNegativeIntsOnly: {
		severity: ERROR,
		code:     ErrCborNegativeIntsOnly,
		desc:     "Expected negative integers only",
	},
	ErrCborUnhandledReadBytesInTermsOfBits: {
		severity: ERROR,
		code:     ErrCborUnhandledReadBytesInTermsOfBits,
		desc:     "Reading multiple bytes in terms of bits is unhandled",
	},
	ErrCborBignumParsingFailed: {
		severity: ERROR,
		code:     ErrCborBignumParsingFailed,
		desc:     "Error converting string to bignum",
	},
	ErrShelleyPayloadInvalid: {
		severity: ERROR,
		code:     ErrShelleyPayloadInvalid,
		desc:     "Shelley payload is invalid",
	},
	ErrShelleyInvalidContainerMode: {
		severity: ERROR,
		code:     ErrShelleyInvalidContainerMode,
		desc:     "Invalid container mode expected",
	},
	ErrShellyUnexpectedCborItem: {
		severity: ERROR,
		code:     ErrShellyUnexpectedCborItem,
		desc:     "Unexpected CBOR item in response from node",
	},
	ErrShellyHandshakeFailed: {
		severity: ERROR,
		code:     ErrShellyHandshakeFailed,
		desc:     "Handshake negotiation failed",
	},
}

// Error string
func (e *CLIError) Error() string {
	return fmt.Sprintf("[%s.%d] %s", e.severity.String(), e.code, e.desc)
}

// Code of the error
func (e *CLIError) Code() int {
	return e.code
}

// Message of the error
func (e *CLIError) Message() string {
	if e != nil {
		return e.desc
	}
	return "Error object is undefined; unable to call Message()"
}

// NewError returns a new error instance with the code
func NewError(code int) *CLIError {
	return NewMessageErrorf(code, "")
}

// NewMessageError returns a new error instance with code and message
func NewMessageError(code int, message string) *CLIError {
	return NewMessageErrorf(code, message)
}

// NewErrorf returns a new error instance with vars
func NewErrorf(code int, v ...interface{}) *CLIError {
	err := cliErrorMap[code]

	if err.code != code {
		return nil
	}

	return &CLIError{
		severity: err.severity,
		code:     err.code,
		desc:     err.desc,
		message:  fmt.Sprintf(err.message, v...),
	}
}

// NewMessageErrorf returns a new error instance with message
func NewMessageErrorf(code int, message string, v ...interface{}) *CLIError {
	err := cliErrorMap[code]

	if err.code != code {
		return nil
	}

	if len(message) == 0 {
		// Use the message if the format string supplied is empty
		message = err.message
	}

	return &CLIError{
		severity: err.severity,
		code:     err.code,
		desc:     err.desc,
		message:  fmt.Sprintf(message, v...),
	}
}

func (t Severity) String() string {
	switch t {
	case ERROR:
		return "E"
	case WARNING:
		return "W"
	case INFO:
		return "I"
	case DEBUG:
		return "D"
	}
	return "U"
}
