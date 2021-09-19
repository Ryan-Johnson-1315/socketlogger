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
	log.SetFlags(socketlogger.EmbeddedFlags)
	
	// Also works with custom loggers
	ErrorLogger := log.New(logger, "ERROR: ", socketlogger.EmbeddedFlags)
	defer logger.Disconnect()

	for i := 0; i < 250; i++ {
		log.Println("Testing from the log.Println() method")
		logger.Dbg("testing")
		ErrorLogger.Println("hello world")
	}
}
