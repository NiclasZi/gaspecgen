package cmd

import (
	"fmt"

	"github.com/Phillezi/nz-mssql/internal/config"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nz-mssql",
	Short: "CLI app for querying a ms SQL db",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := viper.GetString("loglevel")
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			logrus.Warnf("Invalid log level %s, falling back to INFO", level)
			lvl = logrus.InfoLevel
		}
		logrus.SetLevel(lvl)

		logrus.Debugf("Logging level set to %s", lvl)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "See the version of the binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: " + viper.GetString("release"))
	},
}

func init() {

	cobra.OnInitialize(config.InitConfig)

	// Persistent flags
	rootCmd.PersistentFlags().String("loglevel", "info", "Set the logging level (info, warn, error, debug)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().String("db-host", "localhost:3306", "The DB host addr")
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))

	rootCmd.AddCommand(versionCmd)

}
