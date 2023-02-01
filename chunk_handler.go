package grtmp

import (
	"encoding/binary"
	"io"
)

type ChunkHandler struct {
	streamMap map[int]*Stream
	chunkSize uint32
}

func NewChunkHandler() *ChunkHandler {
	return &ChunkHandler{
		streamMap: make(map[int]*Stream),
		chunkSize: defaultChunkSize,
	}
}

type MessageHandler func(*Stream) error

func (c *ChunkHandler) StartLoop(reader io.Reader, msgHandler MessageHandler) error {
	buf := make([]byte, 11)

	for {
		// 5.3.1.1 Chunk Basic Header
		if _, err := io.ReadAtLeast(reader, buf[:1], 1); err != nil {
			return err
		}
		fmt := (buf[0] >> 6) & 0x03
		csId := int(buf[0] & 0x3f)

		switch csId {
		case 0:
			if _, err := io.ReadAtLeast(reader, buf[:1], 1); err != nil {
				return err
			}
			csId = 64 + int(buf[0])
		case 1:
			if _, err := io.ReadAtLeast(reader, buf[:2], 2); err != nil {
				return err
			}
			// Chunk Stream ID 为小端存储
			csId = 64 + int(buf[0]) + int(buf[1])<<8
		}

		// 初始化 Stream
		stream, ok := c.streamMap[csId]
		if !ok {
			stream = NewStream()
			c.streamMap[csId] = stream
		}

		// 5.3.1.2 Chunk Message Header
		switch fmt {
		case 0:
			if _, err := io.ReadAtLeast(reader, buf, 11); err != nil {
				return err
			}
			// 绝对时间戳
			stream.timestamp = GetUint24(buf)
			stream.header.Timestamp = stream.timestamp
			stream.header.MsgLen = GetUint24(buf[3:])
			stream.header.MsgTypeId = buf[6]
			stream.header.MsgStreamId = int(binary.BigEndian.Uint32(buf[7:]))

			stream.msg.Grow(stream.header.MsgLen)
		case 1:
			if _, err := io.ReadAtLeast(reader, buf[:7], 7); err != nil {
				return err
			}
			// 相对时间戳
			stream.timestamp = GetUint24(buf)
			stream.header.Timestamp += stream.timestamp
			stream.header.MsgLen = GetUint24(buf[3:])
			stream.header.MsgTypeId = buf[6]

			stream.msg.Grow(stream.header.MsgLen)
		case 2:
			if _, err := io.ReadAtLeast(reader, buf[:3], 3); err != nil {
				return err
			}
			// 相对时间戳
			stream.timestamp = GetUint24(buf)
			stream.header.Timestamp += stream.timestamp
		case 3:
			// 5.3.1.2.4 Type 3 chunks have No Message Header
		}

		// 5.3.1.3 Extended Timestamp
		if stream.timestamp >= messageHeaderMaxTimestamp {
			if _, err := io.ReadAtLeast(reader, buf[:4], 4); err != nil {
				return err
			}
			extTimestamp := binary.BigEndian.Uint32(buf)
			stream.timestamp = extTimestamp
			switch fmt {
			case 0:
				stream.header.Timestamp = stream.timestamp
			case 1:
				fallthrough
			case 2:
				stream.header.Timestamp = stream.header.Timestamp + stream.timestamp
			}
		}

		// 计算 chunk payload 长度
		var readLen uint32
		if stream.header.MsgLen <= c.chunkSize {
			readLen = stream.header.MsgLen
		} else {
			readLen = stream.header.MsgLen - stream.msg.Len()
			if readLen > c.chunkSize {
				readLen = c.chunkSize
			}
		}

		if _, err := io.ReadFull(reader, stream.msg.WriteBuffer(int(readLen))); err != nil {
			return err
		}

		if stream.msg.Len() == stream.header.MsgLen {
			stream.header.CsId = csId

			//todo: 将aggregate message拆分为多个sub message再调用msgHandler

			if err := msgHandler(stream); err != nil {
				return err
			}

			// 清空stream msg buffer
			stream.msg.Reset()
		}

	}
}

func (c *ChunkHandler) setChunkSize(size uint32) {
	c.chunkSize = size
}

func GetUint24(bs []byte) uint32 {
	return uint32(bs[0])<<16 | uint32(bs[1])<<8 | uint32(bs[0])
}
