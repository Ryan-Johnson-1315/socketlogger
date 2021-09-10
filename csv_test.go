package socketlogger

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	udpCsvServer CsvServer
	udpCsvClient CsvClient
	tcpCsvServer CsvServer
	tcpCsvClient CsvClient
	outputDir    string
)

func init() {
	t := &testing.T{}
	outputDir = t.TempDir()
	udp := 50000
	udpCsvServer = NewUdpCsvServer()
	udpCsvServer.SetOutputCsvDirectory(outputDir)
	udpCsvServer.Bind(Connection{
		Addr: "127.0.0.1",
		Port: udp,
	})
	time.Sleep(1 * time.Second)

	udpCsvClient = NewUdpCsvClient()
	udpCsvClient.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: udp,
	})
	time.Sleep(1 * time.Second)

	tcp := 50001
	tcpCsvServer = NewTcpCsvServer()
	tcpCsvServer.SetOutputCsvDirectory(outputDir)
	tcpCsvServer.Bind(Connection{
		Addr: "127.0.0.1",
		Port: tcp,
	})
	time.Sleep(1 * time.Second)

	tcpCsvClient = NewTcpCsvClient()
	tcpCsvClient.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: tcp,
	})
	time.Sleep(1 * time.Second)
}

func TestCsvCreation(t *testing.T) {
	udpFname := "udp_testing.csv"
	udpCsvClient.NewCsvFile(udpFname, []interface{}{"hello", "from", "udp", "land", udpFname})
	time.Sleep(100 * time.Millisecond) // Creating the file takes longer than sending the message
	if !fileDirExists(outputDir, udpFname) {
		t.Errorf("%s log file was not created", filepath.Join(outputDir, udpFname))
	}

	tcpFname := "tcp_testing.csv"
	tcpCsvClient.NewCsvFile(tcpFname, []interface{}{"hello", "from", "tcp", "land", tcpFname})
	time.Sleep(100 * time.Millisecond) // Creating the file takes longer than sending the message
	if !fileDirExists(outputDir, tcpFname) {
		t.Errorf("%s log file was not created", filepath.Join(outputDir, tcpFname))
	}
}

func TestCsvDuplicate(t *testing.T) {
	udp := []struct {
		input  string
		output string
	}{
		{"udp_1.csv", "udp_1_1.csv"},
		{"udp_2.csv", "udp_2_1.csv"},
		{"udp_3.csv", "udp_3_1.csv"},
		{"udp_4.csv", "udp_4_1.csv"},
		{"udp_5.csv", "udp_5_1.csv"},
		{"udp_6.csv", "udp_6_1.csv"},
		{"udp_7.csv", "udp_7_1.csv"},
	}

	for _, test := range udp {
		createFile(filepath.Join(outputDir, test.input))
		udpCsvClient.AppendRow(test.input, []interface{}{"this", "is", "the", "duplicate", "headers"})
		time.Sleep(100 * time.Millisecond)
		if !fileDirExists(outputDir, test.output) {
			t.Errorf("%s log file was not created", filepath.Join(outputDir, test.output))
		}
	}

	tcp := []struct {
		input  string
		output string
	}{
		{"tcp_1.csv", "tcp_1_1.csv"},
		{"tcp_2.csv", "tcp_2_1.csv"},
		{"tcp_3.csv", "tcp_3_1.csv"},
		{"tcp_4.csv", "tcp_4_1.csv"},
		{"tcp_5.csv", "tcp_5_1.csv"},
		{"tcp_6.csv", "tcp_6_1.csv"},
		{"tcp_7.csv", "tcp_7_1.csv"},
	}

	for _, test := range tcp {
		createFile(filepath.Join(outputDir, test.input))
		tcpCsvClient.AppendRow(test.input, []interface{}{"this", "is", "the", "duplicate", "headers"})
		time.Sleep(100 * time.Millisecond)
		if !fileDirExists(outputDir, test.output) {
			t.Errorf("%s log file was not created", filepath.Join(outputDir, test.output))
		}
	}
}

func TestOutput(t *testing.T) {
	rows := [][]interface{}{
		{"one", "two", 3, 4.4, "yolo"},
		{"one", "two", 3, 4.4, "yolo"},
		{"one", "two", 3, 4.4, "yolo"},
		{"one", "two", 3, 4.4, "yolo"},
		{"one", "two", 3, 4.4, "yolo"},
		{"one", "two", 3, 4.4, "yolo"},
		{"one", "two", 3, 4.4, "yolo"},
	}
	udpFname := "udp_output_test.csv"
	udpCsvClient.NewCsvFile(udpFname, nil)

	for _, row := range rows {
		udpCsvClient.AppendRow(udpFname, row)
	}

	time.Sleep(500 * time.Millisecond)

	ufile, uerr := os.Open(filepath.Join(outputDir, udpFname))
	if uerr != nil {
		t.Errorf("Error opening file %v", uerr)
	}
	ureader := csv.NewReader(ufile)
	outputRows, uerr := ureader.ReadAll()
	if uerr != nil {
		t.Errorf("Error reading csv file %v", uerr)
	}

	if len(outputRows) == 0 {
		t.Errorf("No rows written to file!")
	}

	for i, row := range outputRows {
		for j, cell := range row {
			if fmt.Sprint(rows[i][j]) != cell {
				t.Errorf("Row %d, Col %d did not match. Actual: %s, Expected %s", i, j, fmt.Sprint(rows[i][j]), cell)
			}
		}
	}

	tcpFname := "tcp_output_test.csv"
	for _, row := range rows {
		tcpCsvClient.AppendRow(tcpFname, row)
	}

	time.Sleep(500 * time.Millisecond)

	tfile, terr := os.Open(filepath.Join(outputDir, tcpFname))
	if terr != nil {
		t.Errorf("Error opening file %v", terr)
	}
	treader := csv.NewReader(tfile)
	outputRows, err := treader.ReadAll()
	if err != nil {
		t.Errorf("Error reading csv file %v", err)
	}

	if len(outputRows) == 0 {
		t.Errorf("No rows written to file!")
	}

	for i, row := range outputRows {
		for j, cell := range row {
			if fmt.Sprint(rows[i][j]) != cell {
				t.Errorf("Row %d, Col %d did not match. Actual: %s, Expected %s", i, j, fmt.Sprint(rows[i][j]), cell)
			}
		}
	}
}

func TestBadCsvFile(t *testing.T) {
	userv := NewUdpCsvServer()
	userv.SetOutputCsvDirectory("/dev")
	userv.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 43201,
	})

	uclient := NewUdpCsvClient()
	uclient.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 43201,
	})
	uclient.AppendRow("should_fail.csv", nil)
	time.Sleep(100 * time.Millisecond)
	if fileDirExists("/dev", "should_fail.csv") {
		t.Error("Created file when should not have created it")
	}

	tserv := NewTcpCsvServer()
	tserv.SetOutputCsvDirectory("/dev")
	tserv.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 43202,
	})

	tclient := NewTcpCsvClient()
	tclient.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 43202,
	})
	tclient.AppendRow("should_fail.csv", nil)
	time.Sleep(100 * time.Millisecond)
	if fileDirExists("/dev", "should_fail.csv") {
		t.Error("Created file when should not have created it")
	}
}

func createFile(path string) {
	os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
}
