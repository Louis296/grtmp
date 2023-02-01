package grtmp

import "bytes"

const (
	initMsgBufLen = 4096
)

type Stream struct {
	header    MessageHeader
	msg       StreamMsg
	timestamp uint32
}

func NewStream() *Stream {
	return &Stream{msg: StreamMsg{
		buf: bytes.NewBuffer(make([]byte, initMsgBufLen)), wp: 0},
	}
}

type StreamMsg struct {
	// 禁止直接操作该buf，对于stream msg的所有操作都要使用封装好的方法
	buf *bytes.Buffer
	// 写指针，指向当前buf中未被写入的第一项
	// todo: 读取时，wp向前移动
	// 写入时，wp向后移动
	wp int
}

func (s *StreamMsg) Grow(n uint32) {
	s.buf.Grow(int(n))
}

func (s *StreamMsg) Len() uint32 {
	return uint32(s.wp)
}

// WriteBuffer
// 返回当前Buffer的写指针后长度为n的待写入区域，
// 并将写指针向后移动n位。
func (s *StreamMsg) WriteBuffer(n int) []byte {
	if s.wp+n >= s.buf.Len() {
		s.buf.Grow(n)
	}
	s.wp += n
	return s.buf.Bytes()[s.wp-n : s.wp]
}

func (s *StreamMsg) Reset() {
	s.buf.Reset()
	s.wp = 0
}
