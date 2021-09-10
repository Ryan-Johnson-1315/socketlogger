package socketlogger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type LoggerServer interface {
	SetLogFile(string, string) error
	SetTimeFlags(flags int) error
	Server
}

type loggerserver struct{}

func (l *loggerserver) SetLogFile(dir, name string) error {
	var err error
	var logFile *os.File
	if !fileDirExists(dir, "") {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	logFile, err = os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	return err
}

func (l *loggerserver) SetTimeFlags(flags int) error {
	log.SetFlags(flags)
	return nil
}

func (l *loggerserver) getMessageType() SocketMessage {
	return &LogMessage{}
}

func (l *loggerserver) write(msgs chan SocketMessage) {
	for msg := range msgs {
		log.Printf("%s\n", msg)
	}
}

type UdpLoggerServer struct {
	loggerserver
	udpServer
}

func NewUdpLoggerServer() LoggerServer {
	u := &UdpLoggerServer{}
	u.init(u)
	return u
}

type TcpLoggerServer struct {
	tcpServer
	loggerserver
}

func NewTcpLoggerServer() LoggerServer {
	t := &TcpLoggerServer{}
	t.init(t)
	return t
}

type LoggerClient interface {
	Log(format string, args ...interface{})
	Wrn(format string, args ...interface{})
	Dbg(format string, args ...interface{})
	Err(format string, args ...interface{})
	Success(format string, args ...interface{})
	Client
}

type loggerclient struct {
	msgsToSend chan SocketMessage
}

func (l *loggerclient) setMsgChannel(msgsToSend chan SocketMessage) {
	l.msgsToSend = msgsToSend
}

func (l *loggerclient) Log(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	l.msgsToSend <- newLogMessageCaller(MessageLevelLog, file, line, ok, format, args...)
}

func (l *loggerclient) Wrn(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	l.msgsToSend <- newLogMessageCaller(MessageLevelWrn, file, line, ok, format, args...)
}

func (l *loggerclient) Dbg(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	l.msgsToSend <- newLogMessageCaller(MessageLevelDbg, file, line, ok, format, args...)
}

func (l *loggerclient) Err(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	l.msgsToSend <- newLogMessageCaller(MessageLevelErr, file, line, ok, format, args...)
}

func (l *loggerclient) Success(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	l.msgsToSend <- newLogMessageCaller(MessageLevelSuccess, file, line, ok, format, args...)
}

type UdpLoggerClient struct {
	loggerclient
	udpClient
}

func NewUdpLoggerClient() LoggerClient {
	u := &UdpLoggerClient{}
	u.init(u)
	return u
}

type TcpLoggerClient struct {
	loggerclient
	tcpClient
}

func NewTcpLoggerClient() LoggerClient {
	t := &TcpLoggerClient{}
	t.init(t)
	return t
}
