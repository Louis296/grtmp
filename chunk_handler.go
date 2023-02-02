package grtmp

import (
	"encoding/binary"
	"errors"
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
	var subStream *Stream

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

		// 读取 chunk payload
		if _, err := io.ReadFull(reader, stream.msg.WriteBuffer(int(readLen))); err != nil {
			return err
		}

		if stream.msg.Len() == stream.header.MsgLen {
			stream.header.CsId = csId

			// 将aggregate message拆分为多个sub message再调用msgHandler
			if stream.header.MsgTypeId == AggregateMessage {
				isFirstSub := false
				baseTimeStamp := uint32(0)

				if subStream == nil {
					subStream = NewStream()
				}
				subStream.header.CsId = stream.header.CsId

				for stream.msg.Len() > 0 {
					// 读取sub message header
					if stream.msg.Len() < 11 {
						// 头信息不完整
						// todo: 考虑封装chunk解析中的错误类型
						return errors.New("wrong sub message header len, chunk parse stop")
					}
					subStream.header.MsgTypeId = stream.msg.Next(1)[0]
					subStream.header.MsgLen = GetUint24(stream.msg.Next(3))
					subStream.timestamp = binary.BigEndian.Uint32(stream.msg.Next(4))
					subStream.header.MsgStreamId = int(GetUint24(stream.msg.Next(3)))

					// 计算 timestamp
					if isFirstSub {
						baseTimeStamp = subStream.header.Timestamp
						isFirstSub = false
					}
					subStream.header.Timestamp = stream.header.Timestamp + subStream.timestamp - baseTimeStamp

					// 读取sub message data
					if stream.msg.Len() < subStream.header.MsgLen {
						// 数据信息不完整
						return errors.New("wrong sub message data len, chunk parse stop")
					}
					copy(subStream.msg.WriteBuffer(int(subStream.header.MsgLen)), stream.msg.Next(int(subStream.header.MsgLen)))

					// 调用上层的msgHandler
					if err := msgHandler(subStream); err != nil {
						return err
					}

					// 跳过 Back Pointer
					// 标准中未指定 size，其他源码实现中使用了4Byte，故暂用4Byte
					if stream.msg.Len() < 4 {
						// Back Pointer长度不足
					}
					stream.msg.Next(4)
				}

				// 清空sub msg buffer
				subStream.msg.Reset()
			} else {
				if err := msgHandler(stream); err != nil {
					return err
				}

				// 清空stream msg buffer
				stream.msg.Reset()
			}
		}

	}
}

func (c *ChunkHandler) setChunkSize(size uint32) {
	c.chunkSize = size
}

func GetUint24(bs []byte) uint32 {
	return uint32(bs[0])<<16 | uint32(bs[1])<<8 | uint32(bs[0])
}
