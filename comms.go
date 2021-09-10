package socketlogger

import (
	"net"
)

type comms struct {
	sock               net.Conn
	connectionProtocol string // TCP/UDP
}
