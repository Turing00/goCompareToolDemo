//This is goCompareTool program to check up big files with low footprint of memory.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	numberOfRecord                                                                   int
	fileName1, fileName2, tokenizedCsvFile1RecordField, tokenizedCsvFile2RecordField string
	encodingSymbols                                                                  = map[string]string{"filler": "", "hyphen": "-", "semicolon": ";", "space": " "}
	// standardContractFileFormat = map[int]string{
	//  1: "...",
	// }
	// standardCashflowFileFormat = map[int]string{
	//  1: "...",
	// }
)

func evalCsvRecordOddFieldValue(input string) (output string) {
	switch input {
	case encodingSymbols["filler"]:
		output = "<FILLER>"
	case encodingSymbols["space"]:
		output = "<SPACE>"
	default:
		output = input
	}
	return
}

func checkError(reason string, err error) {
	if err != nil {
		log.Fatal(reason, err)
	}
}

func main() {
	if (len(os.Args) - 1) != 2 {
		fmt.Println("Usage: dumb_diff.exe [filename1] [filename2]")
		os.Exit(1)
	}

	fileName1Match, _ := regexp.MatchString("^.*C[T|F].*$", os.Args[1])
	fileName2Match, _ := regexp.MatchString("^.*C[T|F].*$", os.Args[2])

	if fileName1Match && fileName2Match {
		fileName1, fileName2 = os.Args[1], os.Args[2]
	}

	//handle opening of first and second csv's style file with semicolon separator
	csvFile1Handler, err := os.Open(fileName1)
	checkError("Cannot open file : ", err)
	defer csvFile1Handler.Close()

	csvFile2Handler, err := os.Open(fileName2)
	checkError("Cannot open file : ", err)
	defer csvFile1Handler.Close()

	//read and perform comparison of content from file handler
	csvFile1Scanner, csvFile2Scanner := bufio.NewScanner(csvFile1Handler), bufio.NewScanner(csvFile2Handler)
	fmt.Printf("%s\n", strings.Repeat(encodingSymbols["hyphen"], 6))
	for csvFile1Scanner.Scan() && csvFile2Scanner.Scan() {
		tokenizedCsvFile1Record, tokenizedCsvFile2Record := strings.Split(string(csvFile1Scanner.Bytes()), encodingSymbols["semicolon"]), strings.Split(string(csvFile2Scanner.Bytes()), encodingSymbols["semicolon"])
		numberOfRecord++

		for index1, index2 := 0, 0; index1 < len(tokenizedCsvFile1Record) && index2 < len(tokenizedCsvFile2Record); index1, index2 = index1+1, index2+1 {
			if tokenizedCsvFile1Record[index1] != tokenizedCsvFile2Record[index2] {
				tokenizedCsvFile1RecordField, tokenizedCsvFile2RecordField = evalCsvRecordOddFieldValue(tokenizedCsvFile1Record[index1]), evalCsvRecordOddFieldValue(tokenizedCsvFile2Record[index2])
				fmt.Printf("LINE(%d) : field[%d] - /file %s /token %s ≠ /file %s /token %s\n", numberOfRecord, (index1 + 1), fileName1, tokenizedCsvFile1RecordField, fileName2, tokenizedCsvFile2RecordField)
			}
		}
		fmt.Printf("%s\n", strings.Repeat(encodingSymbols["hyphen"], 6))
	}
}
