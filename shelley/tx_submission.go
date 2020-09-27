package shelley

////////////////////////////////////////////////////////////////////////////////
//
// txSubmissionMessage
//     = msgRequestTxIds
//     / msgReplyTxIds
//     / msgRequestTxs
//     / msgReplyTxs
//     / tsMsgDone
//     / msgReplyKTnxBye
//
// msgRequestTxIds = [0, tsBlocking, txCount, txCount]
// msgReplyTxIds   = [1, [ *txIdAndSize] ]
// msgRequestTxs   = [2, tsIdList ]
// msgReplyTxs     = [3, tsIdList ]
// tsMsgDone       = [4]
// msgReplyKTnxBye = [5]
//
// tsBlocking      = false / true
// txCount         = word16
// tsIdList        = [ *txId ] ; The codec only accepts infinite-length list encoding for tsIdList !
// txIdAndSize     = [txId, txSizeInBytes]
// txId            = int
// txSizeInBytes   = word32
//
////////////////////////////////////////////////////////////////////////////////

type word16 uint
type word32 uint
type word64 uint

type TxSubmissionMessageType uint
type txSizeInBytes word32
type tsBlocking bool
type txCount word16
type tsIdList []*TxID
type TxID int

const (
	TxSubmissionMessageRequestTxIdsType TxSubmissionMessageType = 0
	TxSubmissionMessageReplyTxIdsType                           = 1
	TxSubmissionMessageRequestTxsType                           = 2
	TxSubmissionMessageReplyTxsType                             = 3
	TxSubmissionTsMsgDoneType                                   = 4
	TxSubmissionMessageReplyKTnxByeType                         = 5
)

type TxSubmissionMessageRequestTxIds struct {
	MessageType TxSubmissionMessageType
	TsBlocking  tsBlocking
	TxCount1    word16
	TxCount2    word16
}

type TxSubmissionMessageReplyTxIds struct {
	MessageType TxSubmissionMessageType
	TxIdAndSize []*TxIdAndSize
}

type TxSubmissionMessageRequestTxs struct {
	MessageType TxSubmissionMessageType
	TsIDList    []int
}

type TxSubmissionMessageReplyTxs struct {
	MessageType TxSubmissionMessageType
	TsIdList    []int
}

type TxSubmissionMessageTsMsgDone struct {
	MessageType TxSubmissionMessageType
}

type TxSubmissionMessageReplyKTnxBye struct {
	MessageType TxSubmissionMessageType
}

type TxIdAndSize struct {
	TxId          int
	TxSizeInBytes word32
}
