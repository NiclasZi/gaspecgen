package generator

import (
	"fmt"
	"io"

	"github.com/Phillezi/common/utils/or"
)

type Generator interface {
	Generate(data []map[string]string) error
	GenerateIO(w io.Writer, data []map[string]string) error
}

type GenerationOptions struct {
	SheetName string
}

func GetGenerator(path string, generatorOptions ...GenerationOptions) (Generator, error) {
	if path == "" {
		return &CLIGenerator{}, nil
	}
	opt := or.Or(generatorOptions...)

	switch {
	case hasSuffixCI(path, ".csv"):
		return &CSVGenerator{Filename: path}, nil
	case hasSuffixCI(path, ".xlsx"):
		return &XLSXGenerator{Filename: path, OutSheet: opt.SheetName, Overwrite: true}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", path)
	}
}

func hasSuffixCI(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
