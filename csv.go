package socketlogger

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CsvServer interface {
	SetOutputCsvDirectory(string)
	Server
}

type csvserver struct {
	writers   map[string]*csv.Writer
	outputDir string
	flush     chan bool
}

func (c *csvserver) SetOutputCsvDirectory(dir string) {
	c.outputDir = dir
	if !fileDirExists(dir, "") {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Print(newLogMessage(MessageLevelWrn, "Could not make directory \"%s\", %v", dir, err))
		}
	}
}

func (c *csvserver) buildCsvFile(msg *CsvMessage) *csv.Writer {
	if !fileDirExists(c.outputDir, msg.Filename) ||
		(fileDirExists(c.outputDir, msg.Filename) && c.writers[msg.Filename] == nil) {
		duplicate := false
		fname := filepath.Join(c.outputDir, msg.Filename)
		for i := 1; ; i++ {
			if fileExists(fname) {
				duplicate = true
				fname = filepath.Join(c.outputDir, fmt.Sprintf("%s_%d.csv", strings.Split(msg.Filename, ".csv")[0], i))
			} else {
				break
			}
		}

		if duplicate {
			log.Print(newLogMessage(MessageLevelWrn, "Found previous %s, creating %s", msg.Filename, fname))
		}

		fptr, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)

		if err != nil {
			log.Print(newLogMessage(MessageLevelErr, "Could not open file: %s -> %v", fname, err))
			return nil
		} else {
			log.Print(newLogMessage(MessageLevelSuccess, "File created %s", fname))
		}

		csvWriter := csv.NewWriter(fptr)
		if err != nil {
			log.Print(newLogMessage(MessageLevelErr, "Could not open build csv.Writer -> %v", err))
			return nil
		}
		c.writers[msg.Filename] = csvWriter
	}

	return c.writers[msg.Filename]
}

func (c *csvserver) getMessageType() SocketMessage {
	return &CsvMessage{}
}

func (c *csvserver) initCsvServer() {
	c.writers = make(map[string]*csv.Writer)
}

func (c *csvserver) write(msgs chan SocketMessage) {
	for msg := range msgs {
		if msg.Type() == Csv {
			inst := msg.(*CsvMessage)
			if inst.Filename != "" {
				writer := c.buildCsvFile(inst)
				if writer == nil {
					log.Print(newLogMessage(MessageLevelErr, "csv writer returned as nil!"))
					continue
				}
				// Only need to write the row if it is there
				if len(inst.Row) > 0 {
					writer.Write(transform(inst.Row))
					writer.Flush() // flushes headers & data
				}
			}
		}
	}
	c.flush <- true
}

func (c *csvserver) setFlushChannel(flush chan bool) {
	c.flush = flush
}

type TcpCsvServer struct {
	tcpserver
	csvserver
}

func NewTcpCsvServer() CsvServer {
	t := &TcpCsvServer{}
	t.init(t)
	t.initCsvServer()
	return t
}

type UdpCsvServer struct {
	udpserver
	csvserver
}

func NewUdpCsvServer() CsvServer {
	u := &UdpCsvServer{}
	u.init(u)
	u.initCsvServer()
	return u
}

type CsvClient interface {
	NewCsvFile(fname string, headers []interface{})
	AppendRow(fname string, row []interface{})
	Client
}

type csvclient struct {
	msgsToSend chan SocketMessage
}

func (c *csvclient) setMsgChannel(msgsToSend chan SocketMessage) {
	c.msgsToSend = msgsToSend
}

func (c *csvclient) NewCsvFile(fname string, headers []interface{}) {
	c.msgsToSend <- newCsvMessage(fname, headers)
}

func (c *csvclient) AppendRow(fname string, row []interface{}) {
	c.msgsToSend <- newCsvMessage(fname, row)
}

type UdpCsvClient struct {
	csvclient
	udpClient
}

func NewUdpCsvClient() CsvClient {
	u := &UdpCsvClient{}
	u.init(u)
	return u
}

type TcpCsvClient struct {
	csvclient
	tcpClient
}

func NewTcpCsvClient() CsvClient {
	t := &TcpCsvClient{}
	t.init(t)
	return t
}
