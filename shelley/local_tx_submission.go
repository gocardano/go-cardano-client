package shelley

////////////////////////////////////////////////////////////////////////////////
//
// localTxSubmissionMessage
//     = msgSubmitTx
//     / msgAcceptTx
//     / msgRejectTx
//     / ltMsgDone
//
// msgSubmitTx = [0, transaction ]
// msgAcceptTx = [1]
// msgRejectTx = [2, rejectReason ]
// ltMsgDone   = [3]
//
////////////////////////////////////////////////////////////////////////////////

type LocalTxSubmissionMessageType uint
type transaction int
type rejectReason int

const (
	LocalMessageSubmissionMsgSubmitTxType LocalTxSubmissionMessageType = 0
	LocalMessageSubmissionMsgAcceptTx                                  = 1
	LocalMessageSubmissionMsgRejectTx                                  = 2
	LocalMessageSubmissionLtMsgDone                                    = 3
)

type LocalTxSubmissionMessageMsgSubmitTx struct {
	Type        LocalTxSubmissionMessageType
	Transaction transaction
}

type LocalTxSubmissionMessageMsgAcceptTx struct {
	Type LocalTxSubmissionMessageType
}
type LocalTxSubmissionMessageMsgRejectTx struct {
	Type         LocalTxSubmissionMessageType
	RejectReason rejectReason
}

type LocalTxSubmissionMessageLtMsgDone struct {
	Type LocalTxSubmissionMessageType
}
