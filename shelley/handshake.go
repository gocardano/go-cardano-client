package shelley

// handshakeMessage
//     = msgProposeVersions
//     / msgAcceptVersion
//     / msgRefuse
//
// msgProposeVersions = [0, versionTable]
// msgAcceptVersion   = [1, versionNumber, extraParams]
// msgRefuse          = [2, refuseReason ]
//
// ; CDDL is not expressive enough to describe the all possible values of proposeVersions.
// ; proposeVersions is a tables that maps version numbers to version parameters.
// ; The codec requires that the keys are unique and in ascending order.
// ; This specification only enumerates version numbers from 0..2.
//
// versionNumber = 0 / 1 / 2  ; The test instance of handshake only supports version numbers 1,2 and 3.
// ; versionNumber = uint     ; A real instance may support for example any unsigned integer as version number.
//
// params       = any
// extraParams  = any
// versionTable =
//     { ? 0 => params
//     , ? 1 => params
//     , ? 2 => params
//     }
//
// refuseReason
//     = refuseReasonVersionMismatch
//     / refuseReasonHandshakeDecodeError
//     / refuseReasonRefused
//
// refuseReasonVersionMismatch      = [0, [ *versionNumber ] ]
// refuseReasonHandshakeDecodeError = [1, versionNumber, tstr]
// refuseReasonRefused              = [2, versionNumber, tstr]

import (
	"fmt"

	"github.com/masterjk/go-cardano-client/cbor"
	"github.com/masterjk/go-cardano-client/errors"
	"github.com/masterjk/go-cardano-client/multiplex"

	log "github.com/sirupsen/logrus"
)

const (
	handshakeMessagePropose uint8 = 0
	handshakeMessageAccept  uint8 = 1
	handshakeMessageRefuse  uint8 = 2

	handshakeRefuseReasonVersionMismatch      uint8 = 0
	handshakeRefuseReasonHandshakeDecodeError uint8 = 1
	handshakeRefuseReasonRefused              uint8 = 2
)

type handshakeResponse struct {
	accepted      bool
	versionNumber uint16
	extraParams   uint32
	refuseReason  string
}

func handshakeRequest() *cbor.Array {

	// msgProposeVersions = [0, versionTable]
	// versionTable =
	//     { ? 0 => params
	//     , ? 1 => params
	//     , ? 2 => params
	//     }

	arr := cbor.NewArray()
	arr.Add(cbor.NewPositiveInteger8(handshakeMessagePropose))
	versionTable := cbor.NewMap()
	arr.Add(versionTable)

	versionTable.Add(cbor.NewPositiveInteger32(1), cbor.NewPositiveInteger(764824073))
	versionTable.Add(cbor.NewPositiveInteger32(2), cbor.NewPositiveInteger(764824073))
	versionTable.Add(cbor.NewPositiveInteger32(3), cbor.NewPositiveInteger(764824073))
	return arr
}

func parseHandshakeResponse(c *multiplex.Container) (*handshakeResponse, error) {

	if !c.Header().IsFromResponder() {
		log.WithField("mode", c.Header().ContainerMode()).Error("Expected container mode from responder")
		return nil, errors.NewError(errors.ErrShelleyInvalidContainerMode)
	}

	if len(c.DataItems()) != 1 && c.DataItems()[0].MajorType() != cbor.MajorTypeArray {
		log.Error("Handshake response is expecting an array response with 3 items")
		return nil, errors.NewError(errors.ErrShellyUnexpectedCborItem)
	}

	var response *handshakeResponse

	// msgAcceptVersion   = [1, versionNumber, extraParams]
	// msgRefuse          = [2, refuseReason ]
	// refuseReason
	//     = refuseReasonVersionMismatch
	//     / refuseReasonHandshakeDecodeError
	//     / refuseReasonRefused
	//
	// refuseReasonVersionMismatch      = [0, [ *versionNumber ] ]
	// refuseReasonHandshakeDecodeError = [1, versionNumber, tstr]
	// refuseReasonRefused              = [2, versionNumber, tstr]

	arr := c.DataItems()[0].(*cbor.Array)
	status := arr.Get(0).(*cbor.PositiveInteger8)

	switch status.ValueAsUint8() {
	case handshakeMessageAccept:
		response = &handshakeResponse{
			accepted:      true,
			versionNumber: arr.Get(1).(*cbor.PositiveInteger8).ValueAsUint16(),
			extraParams:   arr.Get(2).(*cbor.PositiveInteger32).ValueAsUint32(),
		}
		break

	case handshakeMessageRefuse:

		// refuseReasonVersionMismatch      = [0, [ *versionNumber ] ]
		// refuseReasonHandshakeDecodeError = [1, versionNumber, tstr]
		// refuseReasonRefused              = [2, versionNumber, tstr]

		refuseReasonArray := arr.Get(1).(*cbor.Array)

		switch refuseReasonArray.Get(0).(*cbor.PositiveInteger8).ValueAsUint8() {
		case handshakeRefuseReasonVersionMismatch:
			parameters := refuseReasonArray.Get(1).(*cbor.Array).ValuesAsString()
			response = &handshakeResponse{
				accepted:     false,
				refuseReason: fmt.Sprintf("Version mistmatch [Parameters: %s]", parameters),
			}
			break
		case handshakeRefuseReasonHandshakeDecodeError:
			response = &handshakeResponse{
				accepted:     false,
				refuseReason: fmt.Sprintf("Handshake Decode Error: %s", refuseReasonArray.ValuesAsString()),
			}
			break
		case handshakeRefuseReasonRefused:
			response = &handshakeResponse{
				accepted:     false,
				refuseReason: fmt.Sprintf("Handshake Refused Error: %s", refuseReasonArray.ValuesAsString()),
			}
			break
		default:
			response = &handshakeResponse{
				accepted:     false,
				refuseReason: fmt.Sprintf("Unhandled handshake error reason: %s", refuseReasonArray.ValuesAsString()),
			}
			break
		}
		break
	}

	return response, nil
}
