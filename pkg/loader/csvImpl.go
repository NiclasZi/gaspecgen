package loader

import (
	"encoding/csv"
	"os"
)

type CSVLoader struct{}

func (l *CSVLoader) Load(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var rows []map[string]string
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		row := map[string]string{}
		for i, h := range headers {
			if i < len(record) {
				row[h] = record[i]
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
