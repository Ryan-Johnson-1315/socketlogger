package socketlogger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	MessageType  string
	messageLevel int
	color        string
)

const (
	Log                 MessageType  = "log"
	Csv                 MessageType  = "csv"
	MessageLevelLog     messageLevel = 0
	MessageLevelWrn     messageLevel = 1
	MessageLevelSuccess messageLevel = 2
	MessageLevelErr     messageLevel = 3
	MessageLevelDbg     messageLevel = 4
	reset               color        = "\033[0m"
	red                 color        = "\033[31m"
	green               color        = "\033[32m"
	yellow              color        = "\033[33m"
	cyan                color        = "\033[36m"
	udpProtocol         string       = "udp"
	tcpProtocol         string       = "tcp"
	bufSize             int          = 16384
)

type SocketMessage interface {
	String() string
	Type() MessageType
}

type LogMessage struct {
	Caller   string       `json:"caller"`
	LogLevel messageLevel `json:"level"`
	Message  string       `json:"message"`
}

type CsvMessage struct {
	Caller   string        `json:"caller"`
	Row      []interface{} `json:"row"`
	Filename string        `json:"csv_filename"`
}

type Connection struct {
	Addr string
	Port int
}

func (l LogMessage) String() string {
	str := string(reset)
	switch l.LogLevel {
	case MessageLevelWrn:
		str += string(yellow)
	case MessageLevelSuccess:
		str += string(green)
	case MessageLevelErr:
		str += string(red)
	case MessageLevelDbg:
		str += string(cyan)
	}
	return str + fmt.Sprintf(" | %s -- %s%s", l.Caller, strings.TrimSuffix(l.Message, "\n"), string(reset))
}

func (LogMessage) Type() MessageType {
	return Log
}

func newLogMessageCaller(lvl messageLevel, file string, line int, ok bool, format string, args ...interface{}) SocketMessage {
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return &LogMessage{
		LogLevel: lvl,
		Caller:   caller,
		Message:  fmt.Sprintf(format, args...),
	}
}

func newLogMessage(lvl messageLevel, format string, args ...interface{}) SocketMessage {
	_, file, line, ok := runtime.Caller(1)
	return newLogMessageCaller(lvl, file, line, ok, format, args...)
}

func (CsvMessage) Type() MessageType {
	return Csv
}

func (c CsvMessage) String() string {
	return fmt.Sprintf("Filename: %v, Row: %v", c.Filename, c.Row)
}

func transform(row []interface{}) []string {
	data := make([]string, len(row))
	for i := 0; i < len(row); i++ {
		if msg, ok := row[i].(*LogMessage); ok {
			data[i] = msg.String()
		} else {
			data[i] = fmt.Sprint(row[i])
		}
	}
	return data
}

func newCsvMessage(fname string, row []interface{}) SocketMessage {
	caller := "unknown"
	_, file, line, ok := runtime.Caller(1)
	if ok {
		paths := strings.Split(file, "/")
		caller = fmt.Sprintf("%s:%d", paths[len(paths)-1], line)
	}

	return &CsvMessage{
		Caller:   caller,
		Filename: fname,
		Row:      row,
	}
}

func fileDirExists(dir, file string) bool {
	return fileExists(filepath.Join(dir, file))
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
