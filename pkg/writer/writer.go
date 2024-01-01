package writer

import (
	"encoding/csv"
	"os"
)

// WriteCSV takes a 2D slice of strings (data) and writes it to a file specified by filePath.
func WriteCSV(data [][]string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err // return any error encountered while writing
		}
	}

	return nil // return nil if no error occurred
}
