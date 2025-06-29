package loader

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/Phillezi/common/utils/or"
	"github.com/Phillezi/gaspecgen/util"
)

type Loader interface {
	Load(path string) ([]map[string]string, error)
	LoadIO(r io.Reader) ([]map[string]string, error)
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

func GetLoaderIO(filename string, r io.Reader, loadingOpts ...LoadOpts) (Loader, error) {
	switch ext := filepath.Ext(filename); ext {
	case ".csv", ".CSV":
		return &CSVLoader{}, nil
	case ".xlsx", ".XLSX":
		opt := or.Or(loadingOpts...)
		return &XLSXLoader{
			Sheet:       opt.Sheet,
			SheetIndex:  opt.SheetIndex,
			ToCamelCase: *or.Or(opt.CamelCase, util.PtrOf(true)),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", filename)
	}
}

func hasSuffixCI(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
