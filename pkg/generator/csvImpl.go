package generator

import (
	"encoding/csv"
	"os"
	"sort"
)

type CSVGenerator struct {
	Filename string
}

func (g *CSVGenerator) Generate(data []map[string]string) error {
	if len(data) == 0 {
		return nil
	}

	f, err := os.Create(g.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Extract headers from first row and sort for consistency
	headers := make([]string, 0, len(data[0]))
	for k := range data[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	// Write header row
	if err := w.Write(headers); err != nil {
		return err
	}

	// Write each row in the order of headers
	for _, row := range data {
		record := make([]string, len(headers))
		for i, h := range headers {
			record[i] = row[h]
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}
