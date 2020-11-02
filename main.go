package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type report struct {
	XMLName xml.Name `xml:"Report"`
	Ps      []p      `xml:"rp>pss>ps"`
}

type p struct {
	Number string `xml:"m,attr"`
	Price  string `xml:"awod,attr"`
}

func xls(a report, b map[string]float64) {
	xlsx := excelize.NewFile()
	// Create a new sheet.
	index := xlsx.NewSheet("Sheet1")
	// Set value of a cell.
	i := 1
	for key, value := range b {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), key)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), value)
		i++
	}
	xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), "Итого:")
	xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), itog(a.Ps))

	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	error := xlsx.SaveAs("./Book1.xlsx")
	if error != nil {
		fmt.Println(error)
	}
}

func itog(f []p) float64 {
	var q float64
	for _, item := range f {
		str := strings.Replace(item.Price, ",", ".", -1)
		i, err := strconv.ParseFloat(str, 32)
		if err == nil {
			q = q + i
		}
	}
	p := math.Round(q*100) / 100
	return p
}

func main() {
	csvFile, err := os.Open("mvz.csv")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("\nSuccessfully opened file mvz.csv\n")

	defer func() {
		if err = csvFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 2
	reader.Comment = '#'

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	xmlFile, err := os.Open("report.xml")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Successfully opened file report.xml\n\n")

	defer func() {
		if err = xmlFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	xmlReadFile, err := ioutil.ReadAll(xmlFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	report := report{}
	xml.Unmarshal(xmlReadFile, &report)

	var sprav = make(map[string]string)
	for _, value := range rawCSVdata {
		sprav[value[0]] = value[1]
	}

	var result = make(map[string]float64)
	for _, item := range report.Ps {
		st := strings.Replace(item.Price, ",", ".", -1)
		strToInt, err := strconv.ParseFloat(st, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		result[sprav[item.Number]] += strToInt
	}

	for _, orderForNumber := range report.Ps {
		fmt.Printf("%+v\n", orderForNumber)
	}

	xls(report, result)

	in := bufio.NewScanner(os.Stdin)
	in.Scan()
}
