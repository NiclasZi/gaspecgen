package generator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type CLIGenerator struct{}

func (g *CLIGenerator) Generate(data []map[string]string) error {
	if len(data) == 0 {
		return errors.New("no data to generate")
	}

	// Use tabwriter to print a neat table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print headers from keys of the first map
	firstRow := data[0]
	headers := []string{}
	for k := range firstRow {
		headers = append(headers, k)
	}
	// Print header row
	for _, h := range headers {
		fmt.Fprintf(w, "%s\t", h)
	}
	fmt.Fprintln(w)

	// Print separator row
	for range headers {
		fmt.Fprintf(w, "--------\t")
	}
	fmt.Fprintln(w)

	// Print data rows
	for _, row := range data {
		for _, h := range headers {
			fmt.Fprintf(w, "%s\t", row[h])
		}
		fmt.Fprintln(w)
	}

	return w.Flush()
}

func (g *CLIGenerator) GenerateIO(w io.Writer, data []map[string]string) error {
	if len(data) == 0 {
		return errors.New("no data to generate")
	}

	// Use tabwriter to print a neat table
	ww := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Print headers from keys of the first map
	firstRow := data[0]
	headers := []string{}
	for k := range firstRow {
		headers = append(headers, k)
	}
	// Print header row
	for _, h := range headers {
		fmt.Fprintf(ww, "%s\t", h)
	}
	fmt.Fprintln(ww)

	// Print separator row
	for range headers {
		fmt.Fprintf(ww, "--------\t")
	}
	fmt.Fprintln(ww)

	// Print data rows
	for _, row := range data {
		for _, h := range headers {
			fmt.Fprintf(ww, "%s\t", row[h])
		}
		fmt.Fprintln(ww)
	}

	return ww.Flush()
}
