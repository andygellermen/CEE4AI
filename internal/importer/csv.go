package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type CSVData struct {
	Headers []string
	Rows    []map[string]string
}

func ReadCSV(path string) (*CSVData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return &CSVData{}, nil
		}
		return nil, fmt.Errorf("read csv headers: %w", err)
	}

	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	rows := make([]map[string]string, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read csv row: %w", err)
		}

		if len(record) != len(headers) {
			return nil, fmt.Errorf("row has %d values, expected %d", len(record), len(headers))
		}

		row := make(map[string]string, len(headers))
		for i, header := range headers {
			row[header] = strings.TrimSpace(record[i])
		}
		rows = append(rows, row)
	}

	return &CSVData{
		Headers: headers,
		Rows:    rows,
	}, nil
}
