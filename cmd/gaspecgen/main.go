package main

import (
	"os"

	"github.com/Phillezi/nz-mssql/cmd/gaspecgen/cli"
)

func main() {
	if err := cli.ExecuteE(); err != nil {
		os.Exit(1)
	}
}
