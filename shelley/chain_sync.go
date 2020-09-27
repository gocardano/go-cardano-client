package shelley

////////////////////////////////////////////////////////////////////////////////
//
// chainSyncMessage
//     = msgRequestNext
//     / msgAwaitReply
//     / msgRollForward
//     / msgRollBackward
//     / msgFindIntersect
//     / msgIntersectFound
//     / msgIntersectNotFound
//     / chainSyncMsgDone
//
// msgRequestNext         = [0]
// msgAwaitReply          = [1]
// msgRollForward         = [2, wrappedHeader, tip]
// msgRollBackward        = [3, point, tip]
// msgFindIntersect       = [4, points]
// msgIntersectFound      = [5, point, tip]
// msgIntersectNotFound   = [6, tip]
// chainSyncMsgDone       = [7]
//
// wrappedHeader = #6.24(bytes .cbor blockHeader)
// tip = [point, uint]
//
// points = [ *point ]
//
////////////////////////////////////////////////////////////////////////////////

type ChainSyncMessageType uint

const (
	ChainSyncMessageRequestNextType       ChainSyncMessageType = 0
	ChainSyncMessageAwaitReplyType                             = 1
	ChainSyncMessageRollForwardType                            = 2
	ChainSyncMessageRollBackwardType                           = 3
	ChainSyncMessageFindIntersectType                          = 4
	ChainSyncMessageIntersectFoundType                         = 5
	ChainSyncMessageIntersectNotFoundType                      = 6
	ChainSyncMessageDoneType                                   = 7
)

type ChainSyncMessageRequestNext struct {
	MessageType ChainSyncMessageType
}

type ChainSyncMessageAwaitReply struct {
	MessageType ChainSyncMessageType
}

type ChainSyncMessageRollForward struct {
	MessageType   ChainSyncMessageType
	WrappedHeader *WrappedHeader
	Tip           *Tip
}

type ChainSyncMessageRollBackward struct {
	MessageType ChainSyncMessageType
	Point       *Point
	Tip         *Tip
}

type ChainSyncMessageFindIntersect struct {
	MessageType ChainSyncMessageType
	Points      []*Point
}

type ChainSyncMessageIntersectFound struct {
	MessageType ChainSyncMessageType
	Point       *Point
	Tip         *Tip
}

type ChainSyncMessageIntersectNotFound struct {
	MessageType ChainSyncMessageType
	Tip         *Tip
}

type ChainSyncMessageDone struct {
	MessageType ChainSyncMessageType
}

type Tip struct {
	Point *Point
	Value uint
}

type WrappedHeader struct {
	Value []byte
}

type Point struct {
	Origin          []interface{}
	BlockHeaderHash *BlockHeaderHash
}

type BlockHeaderHash struct {
	SlotNo word64
	Value  int
}
