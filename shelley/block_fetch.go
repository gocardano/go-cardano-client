package shelley

////////////////////////////////////////////////////////////////////////////////
//
// blockFetchMessage
//      = msgRequestRange
//      / msgClientDone
//      / msgStartBatch
//      / msgNoBlocks
//      / msgBlock
//      / msgBatchDone
//
// msgRequestRange = [0, point, point]
// msgClientDone   = [1]
// msgStartBatch   = [2]
// msgNoBlocks     = [3]
// msgBlock        = [4, #6.24(bytes .cbor block)]
// msgBatchDone    = [5]
//
////////////////////////////////////////////////////////////////////////////////

// BlockFetchMessageType identify the message type for the block fetch protocol
type BlockFetchMessageType uint

const (
	BlockFetchMessageRequestRangeType BlockFetchMessageType = 0
	BlockFetchMessageClientDoneType   BlockFetchMessageType = 1
	BlockFetchMessageStartBatchType   BlockFetchMessageType = 2
	BlockFetchMessageNoBlocksType     BlockFetchMessageType = 3
	BlockFetchMessageBlockType        BlockFetchMessageType = 4
	BlockFetchMessageBatchDoneType    BlockFetchMessageType = 5
)

type BlockFetchMessageRequestRange struct {
	MessageType BlockFetchMessageType
	Point1      *Point
	Point2      *Point
}

type BlockFetchMessageClientDone struct {
	MessageType BlockFetchMessageType
}

type BlockFetchMessageStartBatch struct {
	MessageType BlockFetchMessageType
}

type BlockFetchMessageNoBlocks struct {
	MessageType BlockFetchMessageType
}

type BlockFetchMessageBlock struct {
	MessageType BlockFetchMessageType
	Value       []byte
}

type BlockFetchMessageBatchDone struct {
	MessageType BlockFetchMessageType
}
