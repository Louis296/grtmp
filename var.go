package grtmp

const (
	grtmpVersion = "v0.01"
)

const (
	// rtmp协议中三字节时间戳最大值
	messageHeaderMaxTimestamp uint32 = 0xFFFFFF
	defaultChunkSize          uint32 = 128
)

// rtmp message type id
const (
	// SetChunkSizeMessage
	//
	// 5.4. Protocol Control Messages
	SetChunkSizeMessage              uint8 = 1
	AbortMessage                     uint8 = 2
	AcknowledgementMessage           uint8 = 3
	WindowAcknowledgementSizeMessage uint8 = 5
	SetPeerBandwidthMessage          uint8 = 6

	// UserControlMessage
	//
	// 6.2. User Control Messages
	UserControlMessage uint8 = 4

	// CommandMessageAMF0
	//
	// 7. RTMP Command Messages
	CommandMessageAMF0      uint8 = 20
	CommandMessageAMF3      uint8 = 17
	DataMessageAMF0         uint8 = 18
	DataMessageAMF3         uint8 = 15
	SharedObjectMessageAMF0 uint8 = 19
	SharedObjectMessageAMF3 uint8 = 16
	AudioMessage            uint8 = 8
	VideoMessage            uint8 = 9
	AggregateMessage        uint8 = 22
)

// 7.1.7. User Control Message Events
const (
	UserControlMessageStreamBegin      uint8 = 0
	UserControlMessageStreamEOF        uint8 = 1
	UserControlMessageStreamDry        uint8 = 2
	UserControlMessageSetBufferLength  uint8 = 3
	UserControlMessageStreamIsRecorded uint8 = 4
	UserControlMessagePingRequest      uint8 = 6
	UserControlMessagePingResponse     uint8 = 7
)
