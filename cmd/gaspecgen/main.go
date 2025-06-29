package main

import (
	"os"

	"github.com/NiclasZi/gaspecgen/cmd/gaspecgen/cli"
)

func main() {
	if err := cli.ExecuteE(); err != nil {
		os.Exit(1)
	}
}
