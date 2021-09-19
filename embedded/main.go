package main

import (
	"log"

	"github.com/Ryan-Johnson-1315/socketlogger"
)

func main() {
	logger := socketlogger.NewUdpLoggerClient()
	logger.Connect(socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 40000,
	})
	log.SetOutput(logger)
	log.SetFlags((log.Flags() | log.Lshortfile) &^ (log.Ldate | log.Ltime))
	defer logger.Disconnect()

	for i := 0; i < 250; i++ {
		log.Println("Testing from the log.Prinln() method")
	}
}
