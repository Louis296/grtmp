package grtmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	c0Len     = 1
	c1Len     = 1536
	c2Len     = 1536
	s0Len     = 1
	s1Len     = 1536
	s2Len     = 1536
	c0c1Len   = 1537
	s0s1Len   = 1537
	s0s1s2Len = 3073
)

const (
	version uint8 = 3
)

var random1528 []byte

type HandshakeServer struct {
	c0c1 []byte
}

func (s *HandshakeServer) ReadC0C1(reader io.Reader) error {
	c0c1 := make([]byte, c0c1Len)
	if _, err := io.ReadAtLeast(reader, c0c1, c0c1Len); err != nil {
		return err
	}
	s.c0c1 = c0c1
	return nil
}

func (s *HandshakeServer) WriteS0S1S2(writer io.Writer) error {
	s0s1s2 := make([]byte, s0s1s2Len)
	// s0
	s0s1s2[0] = version

	// s1
	s1 := s0s1s2[s0Len:s0s1Len]
	binary.BigEndian.PutUint32(s1, uint32(time.Now().UnixNano()))
	binary.BigEndian.PutUint32(s1[4:], 0)
	copy(s1[8:], random1528)

	// s2
	s2 := s0s1s2[s0s1Len:]
	binary.BigEndian.PutUint32(s2, uint32(time.Now().UnixNano()))
	c1 := s.c0c1[c0Len:]
	copy(s2[4:], c1[:4])
	copy(s2[8:], c1[8:])

	_, err := writer.Write(s0s1s2)
	return err
}

func (s *HandshakeServer) ReadC2(reader io.Reader) error {
	c2 := make([]byte, c2Len)
	_, err := io.ReadAtLeast(reader, c2, c2Len)
	return err
}

func init() {
	random1528 = make([]byte, 1528)
	key := []byte(fmt.Sprintf("grtmp %v by louis296", grtmpVersion))
	for i := 0; i < 1528; i += len(key) {
		copy(random1528[i:], key)
	}
}
