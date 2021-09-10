package socketlogger

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Client interface {
	Connect(client, server Connection) error

	start()
	buildSocket(local, remote Connection) (string, net.Conn, error)
	writeOverSocket(chan SocketMessage)
	setMsgChannel(chan SocketMessage)
	init(i interface{})
}

type client struct {
	comms
	msgsToSend chan SocketMessage
	remoteAddr net.Addr
	this       interface{}
}

func (c *client) Connect(client, server Connection) error {
	var err error = nil
	inst, ok := c.this.(Client)
	if !ok {
		err = fmt.Errorf(`type is not interface type "Server". Type %t`, c.this)
	} else {
		c.connectionProtocol, c.sock, err = inst.buildSocket(client, server)
		c.msgsToSend <- newLogMessage(MessageLevelSuccess, "Built %s at %s", c.connectionProtocol, c.sock.LocalAddr())
	}
	c.start()
	return err
}

func (c *client) start() {
	c.this.(Client).setMsgChannel(c.msgsToSend)
	go c.this.(Client).writeOverSocket(c.msgsToSend)
}

func (c *client) init(i interface{}) {
	if inst, ok := i.(Client); !ok {
		panic(fmt.Errorf("instance is not of type Client! Type: %T", inst))
	} else {
		c.this = i
		c.msgsToSend = make(chan SocketMessage, 100)
	}
}

type udpClient struct {
	client
}

func (u *udpClient) buildSocket(local Connection, remote Connection) (string, net.Conn, error) {
	sock, err := net.ListenUDP(udpProtocol, &net.UDPAddr{
		IP:   net.ParseIP(local.Addr),
		Port: local.Port,
	})

	u.remoteAddr = &net.UDPAddr{
		IP:   net.ParseIP(remote.Addr),
		Port: remote.Port,
	}
	return "UDP Client", sock, err
}

func (u *udpClient) writeOverSocket(msgsToSend chan SocketMessage) {
	sock, goodSock := u.sock.(*net.UDPConn)
	addr, goodAddr := u.remoteAddr.(*net.UDPAddr)

	if !goodSock || !goodAddr {
		panic(fmt.Errorf("udp client socket/server addr is not *net.UDPConn/*net.UDPAddr. Type: %T/%T", sock, addr))
	} else {
		for msg := range msgsToSend {
			bytes, _ := json.Marshal(msg)
			sock.WriteToUDP(bytes, addr)
		}
	}
}

type tcpClient struct {
	client
	connected bool
}

func (t *tcpClient) buildSocket(local Connection, remote Connection) (string, net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr(tcpProtocol, fmt.Sprintf("%s:%d", remote.Addr, remote.Port))
	if err != nil {
		log.Fatal(err)
	}

	sock, err := net.DialTCP(tcpProtocol, nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	} else {
		time.Sleep(50 * time.Millisecond)
	}

	t.connected = (err == nil)
	return "TCP Client", sock, err
}

func (t *tcpClient) writeOverSocket(msgsToSend chan SocketMessage) {
	if t.connected {
		sock, ok := t.sock.(*net.TCPConn)
		if !ok {
			log.Fatal(fmt.Errorf("socket is not *net.TCPConn type. Type: %T", t.sock))
		} else {
			for msg := range msgsToSend {
				bytes, _ := json.Marshal(msg)
				sock.Write(bytes)
			}
		}
	}
}
