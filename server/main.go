package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Ryan-Johnson-1315/socketlogger"
)

func startLogger(server socketlogger.LoggerServer, ip, dir, file string, port int, micro bool) {
	server.SetLogFile(dir, file)

	if micro {
		server.SetTimeFlags(log.Ldate | log.Ltime)
	} else {
		server.SetTimeFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	}
	err := server.Bind(socketlogger.Connection{
		Addr: ip,
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	servers = append(servers, server)
}

func startCsv(server socketlogger.CsvServer, ip, dir string, port int) {
	server.SetOutputCsvDirectory(dir)

	err := server.Bind(socketlogger.Connection{
		Addr: ip,
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	servers = append(servers, server)
}

var servers []socketlogger.Server

func main() {
	servers = make([]socketlogger.Server, 0)
	ip := flag.String("ip", "127.0.0.1", "IP addr to bind to for logger")
	// Logger configs
	ludp := flag.Int("log_udp", 0, "Enable UDP log messages on this port")
	ltcp := flag.Int("log_tcp", 0, "Enable TCP server log messages on this port")
	ldir := flag.String("log_dir", "logs", "Default directory to save log files to")
	lmicro := flag.Bool("lsecs", false, "Turn off microseconds to log output")
	lext := flag.String("log_ext", "log", "Log file extension")

	// CSV configs
	cudp := flag.Int("csv_udp", 0, "Port to start UDP csv server")
	ctcp := flag.Int("csv_tcp", 0, "Port to start TCP csv server")
	cdir := flag.String("csv_dir", "csv", "Default directory to save csv files to")
	flag.Parse()

	now := time.Now().Format("2006-01-02T15:04:05") + "." + *lext
	logfile := filepath.Join(*ldir, now)

	if *ludp != 0 {
		startLogger(socketlogger.NewUdpLoggerServer(), *ip, *ldir, now, *ludp, *lmicro)
	}

	if *ltcp != 0 {
		startLogger(socketlogger.NewTcpLoggerServer(), *ip, *ldir, now, *ltcp, *lmicro)
	}

	if *ltcp != 0 || *ludp != 0 {
		defer func() {
			log.Println("Log file written to:", logfile)
		}()
	}

	if *cudp != 0 {
		startCsv(socketlogger.NewUdpCsvServer(), *ip, *cdir, *cudp)
	}
	if *ctcp != 0 {
		startCsv(socketlogger.NewTcpCsvServer(), *ip, *cdir, *ctcp)
	}

	if *ctcp != 0 || *cudp != 0 {
		defer func() {
			log.Println("CSV files written to:", *cdir)
		}()
	}

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println()

	for _, server := range servers {
		server.Shutdown()
	}
}
