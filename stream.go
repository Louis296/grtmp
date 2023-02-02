package grtmp

import "bytes"

const (
	initMsgBufLen = 4096
)

type Stream struct {
	header    MessageHeader
	msg       StreamMsgBuf
	timestamp uint32
}

func NewStream() *Stream {
	return &Stream{msg: StreamMsgBuf{
		buf: bytes.NewBuffer(make([]byte, initMsgBufLen)), wp: 0},
	}
}

type StreamMsgBuf struct {
	// 禁止直接操作该buf，对于stream msg的所有操作都必须使用封装好的方法
	buf *bytes.Buffer
	// 写指针，指向当前buf中未被写入的第一项
	// 读取时，wp向前移动
	// 写入时，wp向后移动
	// 理论上wp==buf.Len()，目前仅用作冗余
	wp int
}

func (s *StreamMsgBuf) Grow(n uint32) {
	s.buf.Grow(int(n))
}

func (s *StreamMsgBuf) Len() uint32 {
	return uint32(s.wp)
}

// WriteBuffer
// 返回当前Buffer的写指针后长度为n的待写入区域，
// 并将写指针向后移动n位。
func (s *StreamMsgBuf) WriteBuffer(n int) []byte {
	if s.wp+n >= s.buf.Cap() {
		s.buf.Grow(n)
	}

	// 使buf的底层数组长度延长n位
	s.buf.Write(make([]byte, n))

	s.wp += n
	return s.buf.Bytes()[s.wp-n : s.wp]
}

// Next
// 返回Buffer中前n个byte，将写指针向前移动n位
// 当n大于Buffer中已有数据时，返回nil
func (s *StreamMsgBuf) Next(n int) []byte {
	if n > s.wp {
		return nil
	}
	res := s.buf.Next(n)
	s.wp -= n
	return res
}

// Bytes
// 返回Buffer中所有数值，但不改变写指针
func (s *StreamMsgBuf) Bytes() []byte {
	res := s.buf.Bytes()[:s.wp]
	return res
}

func (s *StreamMsgBuf) Reset() {
	s.buf.Reset()
	s.wp = 0
}
