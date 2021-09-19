package main

import (
	"os"
	"time"

	"github.com/Ryan-Johnson-1315/socketlogger"
)

func main() {
	pid := os.Getpid()
	udp := socketlogger.NewUdpLoggerClient()
	udp.Connect(socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 40000,
	})
	// Makes sure all of the messages get written over the socket
	defer udp.Disconnect()

	tcp := socketlogger.NewTcpLoggerClient()
	tcp.Connect(socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 40001,
	})
	// Makes sure all of the messages get written over the socket
	defer tcp.Disconnect()

	f := func(logger socketlogger.LoggerClient, protocol string) {
		for i := 0; i < 100; i++ {
			switch i % 5 {
			case 0:
				logger.Log("%s testing LOG message from process %d", protocol, pid)
				time.Sleep(time.Millisecond * time.Duration(pid%800))
			case 1:
				logger.Dbg("%s testing DBG message from process %d", protocol, pid)
				time.Sleep(time.Nanosecond * time.Duration(pid%200))
			case 2:
				logger.Wrn("%s testing WARN message from process %d", protocol, pid)
				time.Sleep(time.Millisecond * time.Duration(pid%1200))
			case 3:
				logger.Success("%s testing SUCCESS message from process %d", protocol, pid)
				time.Sleep(time.Nanosecond * time.Duration(pid%300))
			case 4:
				logger.Err("%s testing ERROR message from process %d", protocol, pid)
				time.Sleep(time.Millisecond * time.Duration(pid%100))
			}
		}
	}

	// start first on on thread
	go f(udp, "UDP")
	// wait for this one to finish
	f(tcp, "TCP")
}
