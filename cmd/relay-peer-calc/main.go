package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"relay-peer-calc/pkg/processor"
	"relay-peer-calc/pkg/reader"
	"relay-peer-calc/pkg/writer"
	"strings"
)

func main() {
	inputDir := "../../data/input/"
	outputFilePath := "../../data/output/compiled.csv"

	// Get all files in the input directory
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", inputDir, err)
	}

	// Initialize the data structure to hold processed peer data
	processedPeers := processor.PeerData{}

	// Filter for only CSV files and process them one by one
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".csv") {
			continue // Skip directories and non-CSV files
		}

		filePath := filepath.Join(inputDir, file.Name())
		data, err := reader.ReadCSV(filePath)
		if err != nil {
			log.Printf("Failed to read from %s: %v", filePath, err)
			continue // Skip this file and move to the next
		}

		// Process the data to deduplicate and compile it
		for host, dates := range processor.ProcessPeers(data) {
			if _, exists := processedPeers[host]; !exists {
				processedPeers[host] = make(map[string]string)
			}
			for date, avgPeers := range dates {
				if _, exists := processedPeers[host][date]; !exists {
					processedPeers[host][date] = avgPeers
				}
			}
		}
	}

	monthToFilter := "2023-12" // Specify the month you want to filter by

	// Filter the processed peers data to include only data for the specified month
	filteredPeers := processor.FilterDataByMonth(processedPeers, monthToFilter)

	// Fill in missing data points with "0"
	processor.FillMissingWithZero(filteredPeers)

	// Convert the processed data to the desired 2D slice format using the new function
	processedSlice := processor.PeerDataToCSV(filteredPeers)

	// Write the processed data to a new CSV file
	if err := writer.WriteCSV(processedSlice, outputFilePath); err != nil {
		log.Fatalf("Failed to write to %s: %v", outputFilePath, err)
	}

	fmt.Println("CSV files have been successfully merged, deduplicated, and written to", outputFilePath)
}
