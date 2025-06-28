package loader

import (
	"fmt"

	"github.com/Phillezi/common/utils/or"
	"github.com/Phillezi/nz-mssql/util"
)

type Loader interface {
	Load(path string) ([]map[string]string, error)
}

type LoadOpts struct {
	Sheet      string
	SheetIndex int
	CamelCase  *bool
}

func GetLoader(path string, loadingOpts ...LoadOpts) (Loader, error) {
	switch {
	case hasSuffixCI(path, ".csv"):
		return &CSVLoader{}, nil
	case hasSuffixCI(path, ".xlsx"):
		opt := or.Or(loadingOpts...)
		return &XLSXLoader{Sheet: opt.Sheet, SheetIndex: opt.SheetIndex, ToCamelCase: *or.Or(opt.CamelCase, util.PtrOf(true))}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", path)
	}
}

func hasSuffixCI(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
