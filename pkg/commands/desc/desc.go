package desc

import (
	"github.com/Phillezi/nz-mssql/pkg/persistence/query"
	"github.com/sirupsen/logrus"
)

func Desc(inputFilePath, outputFilePath string) {
	if inputFilePath == "" {
		logrus.Fatalln("input file not specified, use -i to specify it")
	}
	query.PrelESpec(inputFilePath, outputFilePath)
}
