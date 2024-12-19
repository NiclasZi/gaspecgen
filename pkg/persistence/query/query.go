package query

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/Phillezi/nz-mssql/pkg/persistence/manager"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func PrelESpec(inputFilePath, outputFilePath string) {

	db := manager.Get()
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Errorf("Failed to close the database connection: %v", err)
		}
	}()

	inputData, err := os.ReadFile(inputFilePath)
	if err != nil {
		logrus.Fatalf("Failed to read input file: %v", err)
	}

	// Process input file contents
	// Assume the file contains a list of `Art_nr` values, one per line
	artNrList := strings.Split(string(inputData), "\n")

	// Build query parameter placeholders for `Art_nr`
	placeholders := make([]string, len(artNrList))
	for i := range artNrList {
		placeholders[i] = fmt.Sprintf("@p%d", i+1)
	}

	// SQL query with IN clause
	query := fmt.Sprintf(`
SELECT [Art_nr]
      ,(CASE
         WHEN Std.designation IS NOT NULL THEN Std.designation
         ELSE t3.MAKTX
       END) AS Description
      ,[QTY]
      ,[Value]
      ,SM.Name
      ,PT.PurchaseText
      ,PT.ChangedAt
      ,[Ref Designator]
  FROM [DETPLAN].[dbo].['PC2123X1 ArtNr$'] BOM
  JOIN [DETPLAN].[dbo].[SplitStandardMaterialPurchaseText] PT
    ON BOM.Art_nr = PT.Artnr
   AND PT.ISactive = 1
  JOIN [MSupply].[dbo].[StandardManufacturer] SM
    ON SM.[ManufacturerCode] = PT.[ManufacturerCode]
  LEFT JOIN [Msupply].[dbo].[StandardMaterial] Std
    ON Std.Artnr = BOM.Art_nr
  LEFT JOIN [SAP].[dbo].[MaterialText_MAKT] t3
    ON BOM.Art_nr = t3.MATNR
   AND t3.SPRAS = 'E'
 WHERE BOM.Art_nr IN (%s);
`, strings.Join(placeholders, ","))

	// Prepare the statement
	stmt, err := db.GetConnection().Prepare(query)
	if err != nil {
		logrus.Fatalf("Failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Execute the query with the input values
	args := make([]interface{}, len(artNrList))
	for i, artNr := range artNrList {
		args[i] = strings.TrimSpace(artNr) // Clean up whitespace
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		logrus.Fatalf("Query execution failed: %v", err)
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Set headers in the first row
	headers := []string{
		"Art_nr", "Description", "Qty", "Value", "Name", "PurchaseText", "ChangedAt", "RefDesignator",
	}
	for i, header := range headers {
		col := string(rune('A' + i)) // Convert index to Excel column letter
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write query results to the Excel file
	rowIndex := 2
	for rows.Next() {
		var artNr, description, name, purchaseText, changedAt, refDesignator string
		var qty, value sql.NullInt64

		if err := rows.Scan(&artNr, &description, &qty, &value, &name, &purchaseText, &changedAt, &refDesignator); err != nil {
			logrus.Fatalf("Failed to scan row: %v", err)
		}

		// Write data to each column
		values := []interface{}{
			artNr, description, qty.Int64, value.Int64, name, purchaseText, changedAt, refDesignator,
		}
		for i, v := range values {
			col := string(rune('A' + i))
			cell := fmt.Sprintf("%s%d", col, rowIndex)
			f.SetCellValue(sheetName, cell, v)
		}

		rowIndex++
	}

	// Check for errors in the result set
	if err := rows.Err(); err != nil {
		logrus.Fatalf("Error in result set: %v", err)
	}

	// Save the Excel file
	if err := f.SaveAs(outputFilePath); err != nil {
		logrus.Fatalf("Failed to save Excel file: %v", err)
	}

	fmt.Printf("Query results written to %s\n", outputFilePath)
}
