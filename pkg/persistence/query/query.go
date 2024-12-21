package query

import (
	"database/sql"
	"fmt"
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

	ef, err := excelize.OpenFile(inputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := ef.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Hardcoded sheet for now
	erows, err := ef.GetRows("PC2123X1 ArtNr")
	if err != nil {
		logrus.Fatalf("Failed to get rows: %v", err)
	}

	var artNrIndex int = -1
	if len(erows) > 0 {
		for i, header := range erows[0] {
			if strings.TrimSpace(header) == "Art_nr" {
				artNrIndex = i
				break
			}
		}
	}

	// Error if "Art_nr" column is not found
	if artNrIndex == -1 {
		logrus.Fatalf("Failed to find 'Art_nr' column in the header")
	}

	// Extract all values from the "Art_nr" column
	var artNrList []string
	for i, row := range erows {
		if i == 0 {
			continue // Skip header row
		}
		if len(row) > artNrIndex { // Ensure row has enough columns
			artNr := strings.TrimSpace(row[artNrIndex])
			if artNr != "" {
				artNrList = append(artNrList, artNr)
			}
		}
	}

	logrus.Debugln("found these ArtNrs:", artNrList)

	// Build query parameter placeholders for `Art_nr`
	placeholders := make([]string, len(artNrList))
	for i := range artNrList {
		placeholders[i] = fmt.Sprintf("@p%d", i+1)
	}

	query, args := buildQuery(artNrList)

	logrus.Debugln("query:", query)

	// Prepare the statement
	stmt, err := db.GetConnection().Prepare(query)
	if err != nil {
		logrus.Fatalf("Failed to prepare query: %v", err)
	}
	defer stmt.Close()

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

// buildQuery generates the SQL query and arguments for the given Art_nr list.
func buildQuery(artNrList []string) (string, []interface{}) {
	// Prepare placeholders and arguments for the INSERT statement
	placeholders := make([]string, len(artNrList))
	args := make([]interface{}, len(artNrList))
	for i, artNr := range artNrList {
		placeholders[i] = fmt.Sprintf("(@p%d)", i+1)
		args[i] = artNr
	}

	// Base query
	query := `
DECLARE @ArtNrList TABLE ([Art_nr] NVARCHAR(255));

-- Insert values into the table variable
INSERT INTO @ArtNrList ([Art_nr])
VALUES %s;

-- Main query
SELECT ArtNrList.[Art_nr],
       (CASE
            WHEN Std.designation IS NOT NULL THEN Std.designation
            ELSE t3.MAKTX
       END) AS Description,
       PT.QTY,
       PT.Value,
       SM.Name,
       PT.PurchaseText,
       PT.ChangedAt,
       PT.[Ref Designator]
FROM @ArtNrList AS ArtNrList
LEFT JOIN [DETPLAN].[dbo].[SplitStandardMaterialPurchaseText] PT
    ON ArtNrList.Art_nr = PT.Artnr
   AND PT.ISactive = 1
LEFT JOIN [MSupply].[dbo].[StandardManufacturer] SM
    ON SM.[ManufacturerCode] = PT.[ManufacturerCode]
LEFT JOIN Msupply.[dbo].[StandardMaterial] Std
    ON Std.Artnr = ArtNrList.Art_nr
LEFT JOIN SAP.dbo.MaterialText_MAKT t3
    ON ArtNrList.Art_nr = t3.MATNR
   AND t3.SPRAS = 'E';
`

	// Format the query with placeholders
	query = fmt.Sprintf(query, strings.Join(placeholders, ", "))

	return query, args
}
