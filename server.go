package grtmp

import (
	"net"
)

type Server struct {
	addr string
	ln   net.Listener
}

func Listen(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		addr: addr,
		ln:   ln,
	}, nil
}

func (s *Server) Serve() error {
	logger.Info("grtmp server listening addr %v", s.addr)
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	session := NewServerSession(conn)
	session.StartLoop()
}
