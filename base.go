package grtmp

type MessageHeader struct {
	CsId int
	// 绝对时间戳
	Timestamp   uint32
	MsgLen      uint32
	MsgTypeId   uint8
	MsgStreamId int
}
