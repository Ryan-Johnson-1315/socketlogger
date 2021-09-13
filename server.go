package socketlogger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Server interface {
	Bind(c Connection) error
	Shutdown()

	start()
	buildSocket(c Connection) (net.Conn, error)
	listenForMsgsOnSocket(sock net.Conn, msgs chan SocketMessage)
	getMessageType() SocketMessage
	init(interface{})
	write(chan SocketMessage) // Either log file, csv file or console, implemented in server type
	setFlushChannel(chan bool)
}

type server struct {
	comms
	msgs         chan SocketMessage
	this         interface{}
	closeSockets chan bool // This channel will notify to close the sockets
	flushed      chan bool // This makes Shutdown() blocking, allowing everything to be written to console/log file
}

func (s *server) Bind(c Connection) error {
	var err error = nil
	s.sock, err = s.this.(Server).buildSocket(c)
	if err != nil {
		return err
	} else {
		s.start()
		return nil
	}
}

func (s *server) Shutdown() {
	close(s.closeSockets)
	close(s.msgs) // Notifies writer to finish writing
	<-s.flushed   // Waits for writer to flush all data
}

func (s *server) init(i interface{}) {
	if inst, ok := i.(Server); !ok {
		panic(fmt.Errorf("instance is not of type Server! Type: %T", inst))
	} else {
		s.msgs = make(chan SocketMessage, 100)
		s.this = i
		s.closeSockets = make(chan bool)
		s.flushed = make(chan bool)
	}
}

func (s *server) start() {
	s.this.(Server).setFlushChannel(s.flushed)
	go s.this.(Server).write(s.msgs)
	go s.listenForMsgsOnSocket(s.comms.sock, s.msgs)
}

func (s *server) listenForMsgsOnSocket(sock net.Conn, msgs chan SocketMessage) {
	if sock != nil {
		reader := bufio.NewReaderSize(sock, bufSize)
		dec := json.NewDecoder(reader)
		inst, ok := s.this.(Server)
		if !ok {
			panic(fmt.Errorf("server is not a Sever type. Type %T", inst))
		}

		decoded := make(chan SocketMessage, 100)
		running := true
		socketDisconnected := make(chan bool)
		go func() {
			for {
				msg := inst.getMessageType()
				if err := dec.Decode(&msg); err != nil {
					if running && err == io.EOF {
						time.Sleep(50 * time.Nanosecond) // Make sure this message gets printed last
						s.msgs <- newLogMessage(MessageLevelDbg, "Socket disconnected %s", sock.RemoteAddr())
					} else if running {
						s.msgs <- newLogMessage(MessageLevelErr, "ERROR!! %v, unexpected error: %v", err, running)
					}
					sock.Close()
					break
				} else {
					decoded <- msg
				}
			}
			socketDisconnected <- true
		}()

	L:
		for {
			select {
			case msg := <-decoded:
				msgs <- msg
			case <-s.closeSockets:
				running = false
				break L
			}
		}
		sock.Close()
		<-socketDisconnected
	}
}

type udpserver struct {
	server
}

func (u *udpserver) buildSocket(c Connection) (net.Conn, error) {
	sock, err := net.ListenUDP(udpProtocol, &net.UDPAddr{
		IP:   net.ParseIP(c.Addr),
		Port: c.Port,
	})
	if err != nil {
		log.Println(newLogMessage(MessageLevelErr, "Could not create %s: %v", "UDP Server", err))
	} else {
		log.Println(newLogMessage(MessageLevelSuccess, "%s listening at %s", "UDP Server", sock.LocalAddr()))
	}

	return sock, err
}

type tcpserver struct {
	server
}

// Satisfies Server interface
func (t *tcpserver) buildSocket(c Connection) (net.Conn, error) {
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", fmt.Sprintf("%v:%d", c.Addr, c.Port))
	if err == nil {
		log.Println(newLogMessage(MessageLevelSuccess, "%s listening at %s:%d", "TCP Server", c.Addr, c.Port))
		go func() {
			defer listener.Close()
			for {
				// Listen for an incoming connection.
				conn, err := listener.Accept()
				if err != nil {
					log.Println(newLogMessage(MessageLevelErr, "Error accepting: %v", err.Error()).String())
					continue
				}

				go t.this.(Server).listenForMsgsOnSocket(conn, t.msgs)
			}
		}()
	}

	return nil, err
}
