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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Log(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Dbg(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Wrn(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Err(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		udpLoggerClient.Success(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, udpTimeout, time.Since(now)+udpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Log(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Dbg(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Wrn(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Err(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
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
		input := fmt.Sprintf("This is a testing message #%d", i)
		inputLines[i] = formatInput(input)
		tcpLoggerClient.Success(input)
		select {
		case line := <-consoleOut:
			outputLines[i] = line
			if failed {
				close()
				t.Errorf("Message %d did not log in timeout period of %v. Actual:  %v", i, tcpTimeout, time.Since(now)+tcpTimeout)
				t.FailNow()
			}
		case <-ticker.C:
			failed = true
			now = time.Now()
		}
	}
	close()
	if correct, expected := checkOutput(); !correct {
		t.Errorf("Lines did not match:%s", expected)
		t.FailNow()
	}
}

func TestSetcpLoggerClientFile(t *testing.T) {
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

func TestMessageType(t *testing.T) {
	if !checkMsgType(udpLoggerServer) {
		t.Errorf("expected %T actual %T", &LogMessage{}, udpLoggerServer.getMessageType())
	}

	if !checkMsgType(tcpLoggerServer) {
		t.Errorf("expected %T actual %T", &LogMessage{}, tcpLoggerServer.getMessageType())
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
	time.Sleep(50 * time.Millisecond)

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
	time.Sleep(50 * time.Millisecond)

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

func close() {
	writer.Close()
	reader.Close()
}

func checkOutput() (bool, string) {
	close()
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
	paths := strings.Split(file, "/")
	file = paths[len(paths)-1]
	return fmt.Sprintf("%s:%d -- %s", file, line+1, msg) + string(reset)
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
