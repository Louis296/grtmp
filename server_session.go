package grtmp

import (
	"fmt"
	"net"
	"sync"
)

type ServerSession struct {
	hs   HandshakeServer
	conn net.Conn

	closeOnce sync.Once
}

func NewServerSession(conn net.Conn) *ServerSession {
	return &ServerSession{
		hs:   HandshakeServer{},
		conn: conn,
	}
}

func (s *ServerSession) StartLoop() error {
	if err := s.handshake(); err != nil {
		s.close(err)
		return err
	}
	//todo: 处理后续rtmp消息
	return nil
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

func (s *ServerSession) close(err error) {
	s.closeOnce.Do(func() {
		if err != nil {
			//todo: log
			fmt.Printf("rtmp session close, err = %v \n", err)
		}
		if s.conn != nil {
			_ = s.conn.Close()
		}
	})
}
