package grtmp

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestStreamMsgBuf_WriteBuffer(t *testing.T) {
	ast := assert.New(t)
	buf := &StreamMsgBuf{
		buf: bytes.NewBuffer(make([]byte, 1)),
		wp:  0,
	}
	buf.Grow(100)
	rbuf := bytes.NewBuffer([]byte("test stream msg buf write buffer"))
	l := len("test stream msg buf write buffer")
	_, err := io.ReadAtLeast(rbuf, buf.WriteBuffer(l), l)

	ast.Equal(nil, err)
	ast.Equal("test stream msg buf write buffer", string(buf.Bytes()))
}

func TestStreamMsgBuf_Next(t *testing.T) {
	ast := assert.New(t)
	buf := &StreamMsgBuf{
		buf: bytes.NewBuffer(make([]byte, 10)),
		wp:  0,
	}
	copy(buf.WriteBuffer(15), []byte("testteststestss"))

	ast.Equal("test", string(buf.Next(4)))
	ast.Equal("tests", string(buf.Next(5)))
	ast.Equal("testss", string(buf.Next(6)))
}

func TestSliceCap(t *testing.T) {
	buf := make([]byte, 0, 100)
	fmt.Println(len(buf), cap(buf))
	buf = buf[:10]
	fmt.Println(len(buf), cap(buf))
}
