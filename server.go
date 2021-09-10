package socketlogger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Server interface {
	Bind(c Connection) error

	start()
	buildSocket(c Connection) (net.Conn, error)
	listenForMsgsOnSocket(sock net.Conn, msgs chan SocketMessage)
	getMessageType() SocketMessage
	init(i interface{})
	write(chan SocketMessage) // Either log file, csv file or console, implemented in server type
}

type server struct {
	comms
	msgs chan SocketMessage
	this interface{}
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

func (s *server) init(i interface{}) {
	if inst, ok := i.(Server); !ok {
		panic(fmt.Errorf("instance is not of type Server! Type: %T", inst))
	} else {
		s.msgs = make(chan SocketMessage, 100)
		s.this = i
	}
}

func (s *server) start() {
	go s.this.(Server).write(s.msgs)
	go s.listenForMsgsOnSocket(s.comms.sock, s.msgs)
}

func (s *server) listenForMsgsOnSocket(sock net.Conn, msgs chan SocketMessage) {
	if sock != nil {
		go func() {
			reader := bufio.NewReaderSize(sock, bufSize)
			dec := json.NewDecoder(reader)
			inst, ok := s.this.(Server)
			if !ok {
				panic(fmt.Errorf("server is not a Sever type. Type %T", inst))
			}
			for {
				msg := inst.getMessageType()
				err := dec.Decode(&msg)
				if err != nil {
					er := newLogMessage(MessageLevelErr, "ERROR!! %v", err)
					log.Println(er.String())
					sock.Close()
					break
				} else {
					msgs <- msg
				}
			}
		}()
	}
}

type udpServer struct {
	server
}

func (u *udpServer) buildSocket(c Connection) (net.Conn, error) {
	sock, err := net.ListenUDP(udpProtocol, &net.UDPAddr{
		IP:   net.ParseIP(c.Addr),
		Port: c.Port,
	})
	log.Println(newLogMessage(MessageLevelSuccess, "%s listening at %s", "UDP Server", sock.LocalAddr()))

	return sock, err
}

type tcpServer struct {
	server
}

// Satisfies Server interface
func (t *tcpServer) buildSocket(c Connection) (net.Conn, error) {
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
