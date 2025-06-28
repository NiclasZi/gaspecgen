package loader

import (
	"github.com/Phillezi/common/utils/or"
	"github.com/iancoleman/strcase"
	"github.com/xuri/excelize/v2"
)

type XLSXLoader struct {
	Sheet       string
	SheetIndex  int
	ToCamelCase bool
}

func (l *XLSXLoader) Load(path string) ([]map[string]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	sheetName := or.Call(
		func() string { return l.Sheet },
		func() string { return f.GetSheetName(l.SheetIndex) },
	)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	headers := rows[0]
	if l.ToCamelCase {
		for i, h := range headers {
			headers[i] = strcase.ToLowerCamel(h)
		}
	}

	var results []map[string]string
	for _, row := range rows[1:] {
		record := map[string]string{}
		for i, h := range headers {
			if i < len(row) {
				record[h] = row[i]
			} else {
				record[h] = ""
			}
		}
		results = append(results, record)
	}

	return results, nil
}
