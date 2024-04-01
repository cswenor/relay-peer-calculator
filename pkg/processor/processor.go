package processor

import (
	"sort"
	"strings"
)

// PeerData stores the average number of peers for each host on specific dates.
type PeerData map[string]map[string]string // Map of hosts to a map of dates to average peers

// FirstKey returns the first key from the PeerData map.
func FirstKey(data PeerData) string {
	for key := range data {
		return key // Return the first key encountered
	}
	return "" // Return empty if no data
}

// ProcessPeers processes the CSV records to remove duplicate date columns.
func ProcessPeers(records [][]string) PeerData {
	processedPeers := make(PeerData) // Initialize your data structure

	if len(records) == 0 {
		return processedPeers // Return if empty
	}

	// Iterate through the records
	for i, row := range records {
		if i == 0 {
			continue // Skip the header row (dates)
		}

		host := row[0] // The host identifier

		// Initialize the map for this host if it's the first time seeing it
		if _, exists := processedPeers[host]; !exists {
			processedPeers[host] = make(map[string]string)
		}

		// Iterate through the columns for this row (each column is a date)
		for j, avgPeers := range row {
			if j == 0 {
				continue // Skip the host identifier
			}

			date := records[0][j] // The date corresponding to this column

			// If this date hasn't been seen for this host, add the average peers
			if _, exists := processedPeers[host][date]; !exists {
				processedPeers[host][date] = avgPeers
			}
		}
	}

	return processedPeers
}

// Convert the processed data back to a 2D slice (similar to CSV format) if needed
func PeersToSlice(peerData PeerData, dates []string) [][]string {
	var result [][]string
	result = append(result, dates) // Add dates as header

	for host, datesMap := range peerData {
		var row []string
		row = append(row, host)          // Add host
		for _, date := range dates[1:] { // Skip the first header as it's 'host\Time'
			row = append(row, datesMap[date])
		}
		result = append(result, row)
	}

	return result
}

// PeerDataToCSV converts PeerData into a 2D slice of strings in the desired CSV format.
func PeerDataToCSV(peerData PeerData) [][]string {
	// Initialize a structure to hold all unique dates
	allDates := make(map[string]bool)
	for _, dates := range peerData {
		for date := range dates {
			allDates[date] = true
		}
	}

	// Convert the map of all dates to a slice and sort them
	sortedDates := make([]string, 0, len(allDates))
	for date := range allDates {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates) // Sort the dates in ascending order

	// Prepare the 2D slice for the CSV output
	output := make([][]string, 1+len(peerData)) // Allocate space for header plus each host
	header := make([]string, 1, len(sortedDates)+1)
	header[0] = "Relay/Date"                // First cell is the label for the relay column
	header = append(header, sortedDates...) // Append all sorted dates to the header
	output[0] = header                      // Set the header as the first row

	i := 1 // Start filling from the second row
	for host, dates := range peerData {
		row := make([]string, len(sortedDates)+1) // +1 for the host label
		row[0] = host                             // Set host name as the first column
		for j, date := range sortedDates {
			if val, ok := dates[date]; ok {
				row[j+1] = val // +1 because the first column is the host name
			} else {
				row[j+1] = "0" // Fill with "0" if no data for this date
			}
		}
		output[i] = row
		i++
	}

	return output
}

// FillMissingWithZero ensures that every host has a value for every date in PeerData.
// If a value is missing or empty, it fills in "0".
func FillMissingWithZero(peerData PeerData) {
	// First, collect all unique dates across all hosts
	allDates := make(map[string]bool)
	for _, dates := range peerData {
		for date := range dates {
			allDates[date] = true
		}
	}

	// Convert the map of all dates to a sorted slice
	sortedDates := make([]string, 0, len(allDates))
	for date := range allDates {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates) // Sort the dates in ascending order

	// Ensure every host has a value for every date
	for host := range peerData {
		for _, date := range sortedDates {
			value, exists := peerData[host][date]
			if !exists || value == "" {
				peerData[host][date] = "0" // Fill missing or empty values with "0"
			}
		}
	}
}

// FilterDataByMonth filters PeerData to only include data from the specified month.
// The month should be in the format "YYYY-MM".
func FilterDataByMonth(peerData PeerData, month string) PeerData {
	filteredData := make(PeerData) // Initialize a new PeerData structure for filtered data

	// Iterate through the data and keep only the data for the specified month
	for host, dates := range peerData {
		for date, value := range dates {
			if strings.HasPrefix(date, month) { // Check if the date starts with the specified month
				// Initialize the inner map for this host if it hasn't been already
				if _, exists := filteredData[host]; !exists {
					filteredData[host] = make(map[string]string)
				}
				// Add the data for this date to the filtered data
				filteredData[host][date] = value
			}
		}
	}

	return filteredData
}
