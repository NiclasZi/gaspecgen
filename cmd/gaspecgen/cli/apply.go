package cli

import (
	"fmt"
	"os"

	"github.com/Phillezi/gaspecgen/db"
	"github.com/Phillezi/gaspecgen/pkg/generator"
	"github.com/Phillezi/gaspecgen/pkg/loader"
	"github.com/Phillezi/gaspecgen/pkg/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var applyCmd = &cobra.Command{
	Use:   "apply [template.sql]",
	Short: "Apply template SQL file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templatePath := args[0]
		dataPath := viper.GetString("input")

		sqlBytes, err := os.ReadFile(templatePath)
		if err != nil {
			zap.L().Fatal("Failed to read SQL template", zap.Error(err))
		}

		loaderOpts := loader.LoadOpts{
			Sheet:      viper.GetString("sheet-name"),
			SheetIndex: viper.GetInt("sheet-index"),
		}

		r := renderer.NewGoTemplateRenderer()
		var query string

		if dataPath != "" {
			ld, err := loader.GetLoader(dataPath, loaderOpts)
			if err != nil {
				zap.L().Fatal("Failed to get loader", zap.Error(err))
			}

			dataRows, err := ld.Load(dataPath)
			if err != nil {
				zap.L().Fatal("Failed to load input data", zap.Error(err))
			}

			q, err := r.Render(string(sqlBytes), *renderer.FromMapArr(dataRows))
			if err != nil {
				zap.L().Fatal("Failed to render sql query with input data", zap.Error(err))
			}
			query = q
		} else {
			q, err := r.Render(string(sqlBytes), renderer.QueryData{})
			if err != nil {
				zap.L().Fatal("Failed to render sql query with input data", zap.Error(err))
			}
			query = q
		}

		fmt.Println("===QUERY===")
		fmt.Println(query)

		db, err := db.Get()
		if err != nil {
			zap.L().Fatal("Failed to connect to the database", zap.Error(err))
		}
		defer func() {
			if err := db.Close(); err != nil {
				zap.L().Fatal("Failed to close the database connection", zap.Error(err))
			}
		}()

		stmt, err := db.GetConnection().Prepare(query)
		if err != nil {
			zap.L().Fatal("Failed to prepare query", zap.Error(err))
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			zap.L().Fatal("Query execution failed", zap.Error(err))
		}
		defer rows.Close()

		g, err := generator.GetGenerator(viper.GetString("output"), generator.GenerationOptions{
			SheetName: viper.GetString("sheet"),
		})
		if err != nil {
			zap.L().Fatal("Failed to get generator", zap.Error(err))
		}

		columns, err := rows.Columns()
		if err != nil {
			zap.L().Fatal("Failed to get columns", zap.Error(err))
		}

		var results []map[string]string
		for rows.Next() {
			values := make([]any, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				zap.L().Fatal("Failed to scan row", zap.Error(err))
			}

			rowMap := make(map[string]string)
			for i, col := range columns {
				var val string
				if b, ok := values[i].([]byte); ok {
					val = string(b)
				} else if values[i] != nil {
					val = fmt.Sprintf("%v", values[i])
				} else {
					val = ""
				}
				rowMap[col] = val
			}
			results = append(results, rowMap)
		}

		if err := g.Generate(results); err != nil {
			zap.L().Fatal("Failed to generate output", zap.Error(err))
		} else {
			zap.L().Info("Done!")
		}

	},
}

func init() {
	applyCmd.Flags().StringP("input", "i", "", "CSV or XLSX file to inject values from")
	applyCmd.Flags().StringP("output", "o", "", "Output file path for results (not implemented yet)")
	applyCmd.Flags().IntP("sheet-index-in", "s", 0, "Sheet index to get values from (only applies when using xlsx input), zero indexed so first is 0")
	applyCmd.Flags().StringP("sheet-name-in", "S", "", "Sheet name to get values from (only applies when using xlsx input), takes priority over sheet-index-in")
	applyCmd.Flags().String("sheet", "", "Sheet name to output result to (only applies when using xlsx output)")

	viper.BindPFlag("input", applyCmd.Flags().Lookup("input"))
	viper.BindPFlag("output", applyCmd.Flags().Lookup("output"))
	viper.BindPFlag("sheet-index-in", applyCmd.Flags().Lookup("sheet-index-in"))
	viper.BindPFlag("sheet-name-in", applyCmd.Flags().Lookup("sheet-name-in"))
	viper.BindPFlag("sheet", applyCmd.Flags().Lookup("sheet"))

	rootCmd.AddCommand(applyCmd)
}
