package generator

import (
	"os"
	"sort"

	"github.com/xuri/excelize/v2"
)

type XLSXGenerator struct {
	Filename  string
	OutSheet  string
	Overwrite bool
}

func (g *XLSXGenerator) Generate(data []map[string]string) error {
	var f *excelize.File
	var err error

	// Check if file exists
	if _, err = os.Stat(g.Filename); err == nil {
		// Open existing file
		f, err = excelize.OpenFile(g.Filename)
		if err != nil {
			return err
		}
	} else {
		// Create new file
		f = excelize.NewFile()
	}

	sheet := g.OutSheet
	if sheet == "" {
		sheet = "Sheet1"
	}

	// Check if sheet exists, create if not
	if index, _ := f.GetSheetIndex(sheet); index == -1 {
		f.NewSheet(sheet)
	}

	// If overwrite is true, delete all rows in the sheet
	if g.Overwrite {
		// Remove sheet and recreate it (excelize does not provide a direct clear sheet method)
		err = f.DeleteSheet(sheet)
		if err != nil {
			return err
		}
		f.NewSheet(sheet)
	}

	if len(data) == 0 {
		return f.SaveAs(g.Filename)
	}

	// Determine consistent column order from first row
	columns := make([]string, 0, len(data[0]))
	for col := range data[0] {
		columns = append(columns, col)
	}
	sort.Strings(columns)

	// Find start row:
	// If overwrite, start at row 1,
	// else append below existing rows
	startRow := 1
	if !g.Overwrite {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return err
		}
		startRow = len(rows) + 1
	}

	// Write headers if starting at row 1
	if startRow == 1 {
		for colIdx, col := range columns {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, startRow)
			f.SetCellValue(sheet, cell, col)
		}
		startRow++ // next row for data
	}

	// Write data rows
	for rowIdx, row := range data {
		for colIdx, col := range columns {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, startRow+rowIdx)
			f.SetCellValue(sheet, cell, row[col])
		}
	}

	return f.SaveAs(g.Filename)
}
