package grtmp

import (
	"encoding/binary"
	"net"
	"sync"
)

type StreamMsgHandler interface {
	handleSetChunkSizeMessage(stream *Stream) error
	handleCommandMessageAMF0(stream *Stream) error
	handleCommandMessageAMF3(stream *Stream) error
	handleDataMessageAMF0(stream *Stream) error
	handleDataMessageAMF3(stream *Stream) error
	handleAcknowledgeMessage(stream *Stream) error
	handleUserControlMessage(stream *Stream) error
}

type ServerSession struct {
	hs   *HandshakeServer
	ch   *ChunkHandler
	conn net.Conn

	closeOnce sync.Once
}

func NewServerSession(conn net.Conn) *ServerSession {
	return &ServerSession{
		hs:   &HandshakeServer{},
		ch:   NewChunkHandler(),
		conn: conn,
	}
}

func (s *ServerSession) StartLoop() {
	if err := s.handshake(); err != nil {
		s.close(err)
		return
	}
	err := s.ch.StartLoop(s.conn, s.msgHandler)
	s.close(err)
	return
}

func (s *ServerSession) handshake() error {
	if err := s.hs.ReadC0C1(s.conn); err != nil {
		return err
	}
	if err := s.hs.WriteS0S1S2(s.conn); err != nil {
		return err
	}
	if err := s.hs.ReadC2(s.conn); err != nil {
		return err
	}
	return nil
}

func (s *ServerSession) msgHandler(stream *Stream) error {
	switch stream.header.MsgTypeId {
	case SetChunkSizeMessage:
		return s.handleSetChunkSizeMessage(stream)
	case CommandMessageAMF0:
		return s.handleCommandMessageAMF0(stream)
	case CommandMessageAMF3:
		return s.handleCommandMessageAMF3(stream)
	case DataMessageAMF0:
		return s.handleDataMessageAMF0(stream)
	case DataMessageAMF3:
		return s.handleDataMessageAMF3(stream)
	case AcknowledgementMessage:
		return s.handleAcknowledgeMessage(stream)
	case UserControlMessage:
		return s.handleUserControlMessage(stream)
	case AudioMessage:
		fallthrough
	case VideoMessage:
		//todo: 处理rtmp音视频数据
	default:
		// 未知消息类型
		logger.Warn("unknown rtmp message type, type id = %v", stream.header.MsgTypeId)
	}
	return nil
}

func (s *ServerSession) close(err error) {
	s.closeOnce.Do(func() {
		if err != nil {
			logger.Warn("rtmp session close, err = %v", err)
		}
		if s.conn != nil {
			_ = s.conn.Close()
		}
	})
}

//-------- implement of StreamMsgHandler interface ---------

func (s *ServerSession) handleSetChunkSizeMessage(stream *Stream) error {
	size := binary.BigEndian.Uint32(stream.msg.Bytes())
	// 5.4.1. Valid sizes are 1 to 2147483647 (0x7FFFFFFF) inclusive
	if size >= 1 && size <= 0x7FFFFFFF {
		s.ch.setChunkSize(size)
	} else {
		// 违法 chunk size，正常情况下不会执行到此处，处理方式待定
		logger.Error("illegal chunk size %v", size)
	}
	return nil
}

func (s *ServerSession) handleCommandMessageAMF0(stream *Stream) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerSession) handleCommandMessageAMF3(stream *Stream) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerSession) handleDataMessageAMF0(stream *Stream) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerSession) handleDataMessageAMF3(stream *Stream) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerSession) handleAcknowledgeMessage(stream *Stream) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerSession) handleUserControlMessage(stream *Stream) error {
	//TODO implement me
	panic("implement me")
}
