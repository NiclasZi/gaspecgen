package cli

import (
	"fmt"

	viperconf "github.com/Phillezi/common/config/viper"
	zetup "github.com/Phillezi/common/logging/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:  "gaspecctl",
	Long: gaSpecCtl,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zetup.Setup()
	},
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	cobra.OnInitialize(func() { viperconf.InitConfig("nzctl") })

	rootCmd.PersistentFlags().String("loglevel", "info", "Set the logging level (info, warn, error, debug)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().String("profile", "", "Set the logging profile (production or empty)")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	rootCmd.PersistentFlags().Bool("stacktrace", false, "Show the stack trace in error logs")
	viper.BindPFlag("stacktrace", rootCmd.PersistentFlags().Lookup("stacktrace"))

	rootCmd.PersistentFlags().String("db", "mydb", "The db")
	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))

	rootCmd.PersistentFlags().String("db-host", "localhost:1433", "The DB host addr (host:port)")
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))

	rootCmd.PersistentFlags().String("db-user", "myuser@domain.com", "The DB user (user@domain.com)")
	viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))

	rootCmd.PersistentFlags().String("db-password", "mypassword", "The DB password")
	viper.BindPFlag("db-password", rootCmd.PersistentFlags().Lookup("db-password"))

	rootCmd.PersistentFlags().Bool("db-encrypt", true, "If encryption should be used")
	viper.BindPFlag("db-encrypt", rootCmd.PersistentFlags().Lookup("db-encrypt"))

	rootCmd.PersistentFlags().Bool("db-trust-cert", false, "If client should trust server cert")
	viper.BindPFlag("db-trust-cert", rootCmd.PersistentFlags().Lookup("db-trust-cert"))

	rootCmd.AddCommand(versionCmd)
}

func ExecuteE() error {
	return rootCmd.Execute()
}

func GetRootCMD() *cobra.Command {
	return rootCmd
}
