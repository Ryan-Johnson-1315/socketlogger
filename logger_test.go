package socketlogger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	udpLoggerServer LoggerServer
	udpLoggerClient LoggerClient
	tcpLoggerServer LoggerServer
	tcpLoggerClient LoggerClient
	consoleOut      chan string
	inputLines      []string
	outputLines     []string
	reader, writer  *os.File
	longWait        int
)

const (
	NUM_LINES  int           = 5000
	udpTimeout time.Duration = 10 * time.Millisecond
	// Account that TCP is a much slower socket
	tcpTimeout time.Duration = 30 * time.Millisecond
)

func TestUDPLog(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestUDPLog", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Log(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestUDPDbg(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestUDPDbg", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Dbg(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestUDPWrn(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestUDPWrn", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Wrn(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestUDPErr(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestUDPErr", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Err(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestUDPSuccess(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestUDPSuccess", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Success(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestTCPLog(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(tcpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestTCPLog", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Log(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestTCPDbg(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(tcpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestTCPDbg", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Dbg(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestTCPWrn(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(tcpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestTCPWrn", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Wrn(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestTCPErr(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(tcpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestTCPErr", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Err(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestTCPSuccess(t *testing.T) {
	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(tcpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TestTCPSuccess", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Success(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestSetTcpLoggerClientFile(t *testing.T) {
	now := time.Now().Format("2006-01-02T15:04:05") + ".log"
	dir := t.TempDir()
	if err := logFile(NewUdpLoggerServer(), dir, now); err != "" {
		t.Error(err)
	}

	if err := logFile(NewTcpLoggerServer(), dir, now); err != "" {
		t.Error(err)
	}

	if err := logFile(NewTcpLoggerServer(), dir+"/not_there_yet", now); err != "" {
		t.Error(err)
	}

	if err := logFile(NewUdpLoggerServer(), "/dev/should_fail", now); err == "" {
		t.Error(err)
	}
}

func TestUDPNative(t *testing.T) {
	server := NewUdpLoggerServer()
	server.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 40000,
	})
	defer server.Shutdown()

	client := NewUdpLoggerClient()
	client.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 40000,
	})

	defer client.Disconnect()
	log.SetFlags(NativeFlags)
	log.SetOutput(udpLoggerClient)

	time.Sleep(time.Millisecond * 1500)

	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> UDP Native", i)
		inputLines[i] = formatNativeInput(input)
		log.Println(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestTCPNative(t *testing.T) {
	log.SetFlags(NativeFlags)
	log.SetOutput(tcpLoggerClient)

	capture()
	failed := false
	var now time.Time
	for i := 0; i < NUM_LINES; i++ {
		ticker := time.NewTicker(udpTimeout * time.Duration(longWait))
		input := fmt.Sprintf("This is a testing message #%d -> TCP Native", i)
		inputLines[i] = formatNativeInput(input)
		log.Println(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				closeWriters()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	closeWriters()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestMessageType(t *testing.T) {
	if !checkMsgType(udpLoggerServer) {
		t.Errorf("expected %T actual %T", &LogMessage{}, udpLoggerServer.getMessageType())
	}

	if !checkMsgType(tcpLoggerServer) {
		t.Errorf("expected %T actual %T", &LogMessage{}, tcpLoggerServer.getMessageType())
	}
}

func TestTCPDisconnect(t *testing.T) {
	server := NewTcpLoggerServer()
	dir := t.TempDir()
	server.SetLogFile(dir, "tcp_disconnect.log")
	server.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 43000,
	})

	logger := NewTcpLoggerClient()
	logger.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 43000,
	})
	expectedLines := 250
	for i := 0; i < expectedLines; i++ {
		logger.Log("TCP disconnect %d", i)
	}

	logger.Disconnect()                // Blocking
	time.Sleep(500 * time.Millisecond) // Just enough to let the messages get to the server
	server.Shutdown()                  // Blocking

	dat, _ := os.ReadFile(filepath.Join(dir, "tcp_disconnect.log"))
	lines := strings.Split(string(dat), "\n")

	if len(lines) != expectedLines+4 { // +4 is the server and client and tcp disconnect output that gets written to the file, and the last newline after the split
		t.Errorf("Output lines did not match expected. Expected: %d, Actual %d, %s", expectedLines+3, len(lines), filepath.Join(dir, "tcp_disconnect.log"))
	}
}

func TestUDPDisconnect(t *testing.T) {
	server := NewUdpLoggerServer()
	dir := t.TempDir()
	server.SetLogFile(dir, "udp_disconnect.log")
	server.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 43000,
	})

	logger := NewUdpLoggerClient()
	logger.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 43000,
	})
	expectedLines := 250
	for i := 0; i < expectedLines; i++ {
		logger.Log("UDP disconnect %d", i)
	}

	logger.Disconnect()                // Blocking
	time.Sleep(500 * time.Millisecond) // Just enough to let the messages get to the server
	server.Shutdown()                  // Blocking

	dat, _ := os.ReadFile(filepath.Join(dir, "udp_disconnect.log"))
	lines := strings.Split(string(dat), "\n")

	if len(lines) != expectedLines+3 { // +3 is the server and client output that gets written to the file, and the last newline after the split
		t.Errorf("Output lines did not match expected. Expected: %d, Actual %d, %s", expectedLines+3, len(lines), filepath.Join(dir, "udp_disconnect.log"))
	}
}

func BenchmarkUDPLog(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		udpLoggerClient.Log("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkUDPDbg(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		udpLoggerClient.Dbg("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkUDPWrn(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		udpLoggerClient.Wrn("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkUDPErr(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		udpLoggerClient.Err("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkUDPSuccess(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		udpLoggerClient.Success("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkTCPLog(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		tcpLoggerClient.Log("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkTCPDbg(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		tcpLoggerClient.Dbg("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkTCPWrn(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		tcpLoggerClient.Wrn("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkTCPErr(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		tcpLoggerClient.Err("Benching with a long message %s", "***********************************************************************************")
	}
}

func BenchmarkTCPSuccess(b *testing.B) {
	for i := 0; i < NUM_LINES; i++ {
		tcpLoggerClient.Success("Benching with a long message %s", "***********************************************************************************")
	}
}

func init() {
	tcp := 60001
	tcpLoggerServer = NewTcpLoggerServer()
	tcpLoggerServer.SetTimeFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	tcpLoggerServer.Bind(Connection{
		Addr: "127.0.0.1",
		Port: tcp,
	})

	tcpLoggerClient = NewTcpLoggerClient()
	tcpLoggerClient.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: tcp,
	})

	udp := 60000
	udpLoggerServer = NewUdpLoggerServer()
	udpLoggerServer.SetTimeFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	udpLoggerServer.Bind(Connection{
		Addr: "127.0.0.1",
		Port: udp,
	})

	udpLoggerClient = NewUdpLoggerClient()
	udpLoggerClient.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: udp,
	})

	consoleOut = make(chan string, 25)
	inputLines = make([]string, NUM_LINES)
	outputLines = make([]string, NUM_LINES)
	time.Sleep(100 * time.Millisecond)
	log.Println("Multiplying wait time by", os.Getenv("LONG_WAIT"))
	if env := os.Getenv("LONG_WAIT"); env != "" {
		var err error
		longWait, err = strconv.Atoi(env)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		longWait = 1
	}
}

func capture() {
	var err error

	reader, writer, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr

	os.Stderr = writer
	os.Stdout = writer
	log.SetOutput(writer)

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			consoleOut <- scanner.Text()
		}

		os.Stdout = stdout
		os.Stderr = stderr
	}()
	time.Sleep(50 * time.Millisecond) // Reset everything
}

func closeWriters() {
	writer.Close()
	reader.Close()
}

func checkOutput() (bool, string) {
	closeWriters()
	success := true
	expected := ""
	for i, fullLine := range outputLines {
		output := strings.Split(fullLine, "| ")[1]
		if same := strings.Compare(output, inputLines[i]); same != 0 {
			success = false
			expected = fmt.Sprintf(`
Output:   "%s"
Expected: "%s"`, strings.TrimSuffix(output, "\n"), inputLines[i])
			break
		}
	}
	return success, expected
}

func formatInput(msg string) string {
	_, file, line, _ := runtime.Caller(1)
	fname := filepath.Base(file)
	return fmt.Sprintf("%s:%d -- %s", fname, line+1, msg) + string(reset)
}

func formatNativeInput(msg string) string {
	_, file, line, _ := runtime.Caller(1)
	fname := filepath.Base(file)
	return fmt.Sprintf("%s:%d %s", fname, line+1, msg)
}

func logFile(server LoggerServer, dir, fname string) string {
	msg := ""
	server.SetLogFile(dir, fname)

	if !fileDirExists(dir, fname) {
		msg = fmt.Sprintf("%s log file was not created", filepath.Join(dir, fname))
	}
	return msg
}

func checkMsgType(server LoggerServer) bool {
	switch server.getMessageType().(type) {
	case *LogMessage:
		return true
	default:
		return false
	}
}
