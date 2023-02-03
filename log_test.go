package grtmp

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	ast := assert.New(t)

	errMsg := "test error log content"
	warnMsg := "test warn log content"
	infoMsg := "test info log content"
	l := len("[grtmp] 2023/02/03 22:39:31.629421 ")
	buf := &bytes.Buffer{}
	SetLogOutput(buf)

	SetLogLevel(Info)
	logger.Error(errMsg)
	ast.Equal(errPrefix+errMsg+"\n", buf.String()[l:])
	buf.Reset()
	logger.Warn(warnMsg)
	ast.Equal(warnPrefix+warnMsg+"\n", buf.String()[l:])
	buf.Reset()
	logger.Info(infoMsg)
	ast.Equal(infoPrefix+infoMsg+"\n", buf.String()[l:])
	buf.Reset()

	SetLogLevel(Warn)
	logger.Error(errMsg)
	ast.Equal(errPrefix+errMsg+"\n", buf.String()[l:])
	buf.Reset()
	logger.Warn(warnMsg)
	ast.Equal(warnPrefix+warnMsg+"\n", buf.String()[l:])
	buf.Reset()
	logger.Info(infoMsg)
	ast.Equal("", buf.String())
	buf.Reset()

	SetLogLevel(Error)
	logger.Error(errMsg)
	ast.Equal(errPrefix+errMsg+"\n", buf.String()[l:])
	buf.Reset()
	logger.Warn(warnMsg)
	logger.Info(infoMsg)
	ast.Equal("", buf.String())
	buf.Reset()

	SetLogLevel(Silent)
	logger.Error(errMsg)
	logger.Warn(warnMsg)
	logger.Info(infoMsg)
	ast.Equal("", buf.String())
	buf.Reset()
}
