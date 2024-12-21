//go:debug x509negativeserial=1
package main

import (
	"github.com/Phillezi/nz-mssql/cmd"
	"github.com/spf13/viper"
)

var buildTimestamp = "19700101000000"

func main() {
	viper.Set("release", "release-"+buildTimestamp)
	cmd.Execute()
}
