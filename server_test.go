package grtmp

import (
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	server, err := Listen(":9000")
	if err != nil {
		t.Failed()
	}
	go func() {
		err := server.Serve()
		if err != nil {
			t.Failed()
		}
	}()
	_, err = net.Dial("tcp", ":9000")
	if err != nil {
		t.Failed()
	}
}
