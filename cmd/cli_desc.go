package cmd

import (
	"log"

	"github.com/Phillezi/nz-mssql/pkg/commands/desc"
	"github.com/spf13/cobra"
)

var descCmd = &cobra.Command{
	Use:   "desc",
	Short: "Get descriptions by article numbers",
	Run: func(cmd *cobra.Command, args []string) {
		inputFilePath, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}
		outputFilePath, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}
		desc.Desc(inputFilePath, outputFilePath)
	},
}

func init() {
	descCmd.Flags().StringP("input", "i", "", "Input file (that has the article numbers)")
	descCmd.Flags().StringP("output", "o", "./out.xlsx", "Output file (that has the descriptions)")
	rootCmd.AddCommand(descCmd)
}
