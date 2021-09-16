package main

import (
	"github.com/Ryan-Johnson-1315/socketlogger"
)

func main() {
	csv := socketlogger.NewUdpCsvClient()
	csv.Connect(socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 0,
	}, socketlogger.Connection{
		Addr: "127.0.0.1",
		Port: 50000,
	})

	// Create the multiple csv files for on client
	class2019 := "class_2019.csv"
	class2020 := "class_2020.csv"

	db2019 := [][]interface{}{
		{"Jackelyn Vanhoose", "2000-09-29", 3.4, "Pomona College, California"},
		{"Christoper Westover", "2000-10-12", 3.9, "Stanford University, California"},
		{"Marcus Courtois", "2001-10-19", 3.30, "Union College, New York"},
		{"Leena Bodner", "2000-10-27", 4.0, "United States Air Force Academy, Colorado"},
		{"Hoyt Stermer", "2001-12-09", 2.0, "University of California, Los Angeles, California"},
	}

	db2020 := [][]interface{}{
		{"Carina Haws", "2001-09-16", 3.1, "Washington University in Saint Louis, Missouri"},
		{"Alice Gill", "2002-09-17", 3.9, "Rice University, Texas"},
		{"Sharika Teeple", "2001-09-27", 2.9, "University of Maryland, College Park, Maryland"},
		{"Earl Friel", "2002-12-03", 3.6, "University of Richmond, Virginia"},
		{"Earlean Numbers", "2002-12-23", 3.3, "Mount Holyoke College, Massachusetts"},
	}

	csv.NewCsvFile(class2019, []interface{}{"Name", "Birthday", "GPA", "College"})
	csv.NewCsvFile(class2020, []interface{}{"Name", "Birthday", "GPA", "College"})

	for _, row := range db2019 {
		csv.AppendRow(class2019, row)
	}

	for _, row := range db2020 {
		csv.AppendRow(class2020, row)
	}

	csv.Disconnect()
}
