package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// ReadCSV takes a file path and returns its contents as a slice of slices of strings.
func ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = true

	var lines [][]string
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error in line: %v\n", record)
			return nil, err
		}
		lines = append(lines, record)
	}

	return lines, nil
}
