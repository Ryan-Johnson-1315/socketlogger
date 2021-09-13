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
	outputDir  string
	incrememnt int
)

func init() {
	t := &testing.T{}
	outputDir = t.TempDir()
	incrememnt = 0
}

func TestCsvCreation(t *testing.T) {
	userver, uclient := newUdp()
	tserver, tclient := newTcp()
	defer func() {
		uclient.Disconnect()
		tclient.Disconnect()
		time.Sleep(50 * time.Millisecond)
		tserver.Shutdown()
		userver.Shutdown()
		time.Sleep(50 * time.Millisecond)
	}()

	uclient.NewCsvFile("udp_creation.csv", []interface{}{"hello", "from", "udp", "land", "udp_creation.csv"})

	time.Sleep(100 * time.Millisecond) // Creating the file takes longer than sending the message
	if !fileDirExists(outputDir, "udp_creation.csv") {
		t.Errorf("%s log file was not created", filepath.Join(outputDir, "udp_creation.csv"))
	}

	tclient.NewCsvFile("tcp_creation.csv", []interface{}{"hello", "from", "udp", "land", "tcp_creation.csv"})
	time.Sleep(100 * time.Millisecond) // Creating the file takes longer than sending the message
	if !fileDirExists(outputDir, "tcp_creation.csv") {
		t.Errorf("%s log file was not created", filepath.Join(outputDir, "tcp_creation.csv"))
	}
}

func TestDisconnectCsv(t *testing.T) {
	userver, uclient := newUdp()
	tserver, tclient := newTcp()

	rows := [][]interface{}{
		{"Carina Haws", "2001-09-16", 3.1, "Washington University in Saint Louis, Missouri"},
		{"Alice Gill", "2002-09-17", 3.9, "Rice University, Texas"},
		{"Sharika Teeple", "2001-09-27", 2.9, "University of Maryland, College Park, Maryland"},
		{"Earl Friel", "2002-12-03", 3.6, "University of Richmond, Virginia"},
		{"Earlean Numbers", "2002-12-23", 3.3, "Mount Holyoke College, Massachusetts"},
	}
	ufname := "udp_disconnect.csv"
	uclient.NewCsvFile(ufname, nil)
	for _, row := range rows {
		uclient.AppendRow(ufname, row)
	}

	uclient.Disconnect()
	time.Sleep(500 * time.Millisecond) // Just enough to let the messages get to the server over the socket
	userver.Shutdown()

	file, terr := os.Open(filepath.Join(outputDir, ufname))
	if terr != nil {
		t.Errorf("Error opening file %v", terr)
	}
	ureader := csv.NewReader(file)
	uoutputRows, uerr := ureader.ReadAll()
	if uerr != nil {
		t.Errorf("Error reading csv file %v", uerr)
	}

	if len(uoutputRows) == 0 {
		t.Errorf("No rows written to file!")
	}

	for i, row := range uoutputRows {
		for j, cell := range row {
			if fmt.Sprint(rows[i][j]) != cell {
				t.Errorf("Row %d, Col %d did not match. Actual: %s, Expected %s", i, j, fmt.Sprint(rows[i][j]), cell)
			}
		}
	}

	tfname := "udp_disconnect.csv"
	tclient.NewCsvFile(tfname, nil)
	for _, row := range rows {
		tclient.AppendRow(tfname, row)
	}

	tclient.Disconnect()
	time.Sleep(500 * time.Millisecond) // Just enough to let the messages get to the server over the socket
	tserver.Shutdown()

	tfile, terr := os.Open(filepath.Join(outputDir, tfname))
	if terr != nil {
		t.Errorf("Error opening file %v", terr)
	}
	treader := csv.NewReader(tfile)
	toutputRows, terr := treader.ReadAll()
	if terr != nil {
		t.Errorf("Error reading csv file %v", terr)
	}

	if len(toutputRows) == 0 {
		t.Errorf("No rows written to file!")
	}

	for i, row := range toutputRows {
		for j, cell := range row {
			if fmt.Sprint(rows[i][j]) != cell {
				t.Errorf("Row %d, Col %d did not match. Actual: %s, Expected %s", i, j, fmt.Sprint(rows[i][j]), cell)
			}
		}
	}
}

func TestCsvDuplicate(t *testing.T) {
	userver, uclient := newUdp()
	tserver, tclient := newTcp()
	defer func() {
		uclient.Disconnect()
		tclient.Disconnect()
		time.Sleep(50 * time.Millisecond)
		tserver.Shutdown()
		userver.Shutdown()
	}()
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
		uclient.AppendRow(test.input, []interface{}{"this", "is", "the", "duplicate", "headers"})
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
		tclient.AppendRow(test.input, []interface{}{"this", "is", "the", "duplicate", "headers"})
		time.Sleep(100 * time.Millisecond)
		if !fileDirExists(outputDir, test.output) {
			t.Errorf("%s log file was not created", filepath.Join(outputDir, test.output))
		}
	}
}

func TestOutput(t *testing.T) {
	userver, uclient := newUdp()
	tserver, tclient := newTcp()
	defer func() {
		uclient.Disconnect()
		tclient.Disconnect()
		time.Sleep(50 * time.Millisecond)
		tserver.Shutdown()
		userver.Shutdown()
	}()
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
	uclient.NewCsvFile(udpFname, nil)

	for _, row := range rows {
		uclient.AppendRow(udpFname, row)
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
		tclient.AppendRow(tcpFname, row)
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
	userver, uclient := newUdp()
	tserver, tclient := newTcp()
	defer func() {
		uclient.Disconnect()
		tclient.Disconnect()
		time.Sleep(50 * time.Millisecond)
		tserver.Shutdown()
		userver.Shutdown()
	}()
	userv := NewUdpCsvServer()
	userv.SetOutputCsvDirectory("/dev")
	userv.Bind(Connection{
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

	tclient.AppendRow("should_fail.csv", nil)
	time.Sleep(100 * time.Millisecond)
	if fileDirExists("/dev", "should_fail.csv") {
		t.Error("Created file when should not have created it")
	}
}

func createFile(path string) {
	os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
}

func newUdp() (CsvServer, CsvClient) {
	server := NewUdpCsvServer()
	server.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 42000 + incrememnt,
	})
	server.SetOutputCsvDirectory(outputDir)

	client := NewUdpCsvClient()
	client.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 42000 + incrememnt,
	})
	incrememnt++
	return server, client
}

func newTcp() (CsvServer, CsvClient) {
	server := NewTcpCsvServer()
	server.Bind(Connection{
		Addr: "127.0.0.1",
		Port: 43001 + incrememnt,
	})
	server.SetOutputCsvDirectory(outputDir)
	client := NewTcpCsvClient()
	client.Connect(Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, Connection{
		Addr: "127.0.0.1",
		Port: 43001 + incrememnt,
	})
	incrememnt++
	return server, client
}
