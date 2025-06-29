package cli

import (
	"fmt"
	"time"

	"github.com/NiclasZi/gaspecgen/db"
	"github.com/NiclasZi/gaspecgen/internal/server"
	"github.com/NiclasZi/gaspecgen/util"
	viperconf "github.com/Phillezi/common/config/viper"
	"github.com/Phillezi/common/interrupt"
	zetup "github.com/Phillezi/common/logging/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:  "gaspecgen",
	Long: gaSpecGen,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zetup.Setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.Get()
		if err != nil {
			zap.L().Fatal("Failed to connect to db", zap.Error(err))
			interrupt.GetInstance().Shutdown()
		}
		defer func() {
			if err := db.Close(); err != nil {
				zap.L().Error("Failed to close the database connection", zap.Error(err))
			}
		}()
		s := server.New(interrupt.GetInstance().Context(), 8080)
		var errCh chan error = make(chan error, 1)
		go func() {
			if err := s.Start(); err != nil {
				errCh <- err
			}
		}()

		if viper.GetBool("open-browser") {
			time.AfterFunc(500*time.Millisecond, func() { util.Open(s.Addr()) })
		}

		select {
		case err := <-errCh:
			zap.L().Fatal("Server sent error", zap.Error(err))
			return // not necessary
		case <-interrupt.GetInstance().Context().Done():
			return
		}

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
	cobra.OnInitialize(func() { viperconf.InitConfig("gaspecgen", "config") })

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

	rootCmd.Flags().Bool("open-browser", false, "Open the url in the browser on server startup")
	viper.BindPFlag("open-browser", rootCmd.Flags().Lookup("open-browser"))

	rootCmd.AddCommand(versionCmd)
}

func ExecuteE() error {
	return rootCmd.Execute()
}

func GetRootCMD() *cobra.Command {
	return rootCmd
}
