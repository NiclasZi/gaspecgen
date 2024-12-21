package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Phillezi/nz-mssql/pkg/persistence/manager"
	"github.com/Phillezi/nz-mssql/pkg/xlsx"
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

	qtyList, err := xlsx.GetByHeader[int](erows, "QTY")
	if err != nil {
		logrus.Fatal(err)
	}
	artNrList, err := xlsx.GetByHeader[string](erows, "Art_nr")
	if err != nil {
		logrus.Fatal(err)
	}
	refDesignatorList, err := xlsx.GetByHeader[string](erows, "Ref Designator")
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugln("found these qtys:", qtyList)
	logrus.Debugln("found these ArtNrs:", artNrList)
	logrus.Debugln("found these refdesignators:", refDesignatorList)

	query, args := buildQuery(qtyList, artNrList, refDesignatorList)

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

		// Scan the row values into variables
		if err := rows.Scan(&artNr, &description, &qty, &value, &name, &purchaseText, &changedAt, &refDesignator); err != nil {
			logrus.Fatalf("Failed to scan row: %v", err)
		}

		// Prepare the row data, handling NullInt64 values
		values := []interface{}{
			artNr,
			description,
			// Handle NullInt64 for Qty and Value
			getNullInt64Value(qty),
			getNullInt64Value(value),
			name,
			purchaseText,
			changedAt,
			refDesignator,
		}

		// Write the values to the respective cells in the row
		for i, v := range values {
			col := string(rune('A' + i)) // Convert index to Excel column letter
			cell := fmt.Sprintf("%s%d", col, rowIndex)
			f.SetCellValue(sheetName, cell, v)
		}

		// Move to the next row (in the excel sheet)
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
func buildQuery(qtyList []int, artNrList []string, refDesignatorList []string) (string, []interface{}) {
	// Prepare placeholders and arguments for the INSERT statement

	if len(qtyList) != len(artNrList) || len(artNrList) != len(refDesignatorList) {
		logrus.Fatal("Different length of qty, artnr and refdesignator (lists) is not allowed!")
	}

	placeholderLen := len(artNrList)
	argsLen := len(qtyList) + len(artNrList) + len(refDesignatorList)
	placeholders := make([]string, placeholderLen)
	args := make([]interface{}, argsLen)

	j := 0
	for i := 0; i < placeholderLen; i++ {
		placeholders[i] = fmt.Sprintf("(@p%d, @p%d, @p%d)", j+1, j+2, j+3)
		args[j] = qtyList[i]
		args[j+1] = artNrList[i]
		args[j+2] = refDesignatorList[i]
		j += 3
	}

	// Base query
	query := `
DECLARE @BOM_List TABLE (
	[QTY] INT,
	[Art_nr] NVARCHAR(255),
	[Ref_Designator] TEXT
);

-- Insert values into the table variable
INSERT INTO @BOM_List (
	[QTY],
	[Art_nr],
	[Ref_Designator]
)
VALUES %s;

-- Query
SELECT DISTINCT
    BOM.[Art_nr],                                  -- Column 1: Art_nr
    BOM.[QTY],                                     -- Column 2: QTY
    CASE
        WHEN Std.designation IS NOT NULL THEN Std.designation
        ELSE t3.MAKTX
    END AS [Description],                          -- Column 3: Description (conditional)
    Std.[Dimension],                               -- Column 4: Dimension
    SM.[Name] AS ManufacturerName,                 -- Column 5: ManufacturerName
    PT.[PurchaseText]                              -- Column 6: PurchaseText
FROM 
    @BOM_List BOM
JOIN 
    [DETPLAN].[dbo].[SplitStandardMaterialPurchaseText] PT
    ON BOM.[Art_nr] = PT.[Artnr]
    AND PT.[ISactive] = 1
JOIN 
    [MSupply].[dbo].[StandardManufacturer] SM
    ON SM.[ManufacturerCode] = PT.[ManufacturerCode]
LEFT JOIN 
    [MSupply].[dbo].[StandardMaterial] Std
    ON BOM.[Art_nr] = Std.[Artnr]
LEFT JOIN SAP.dbo.MaterialText_MAKT t3
	ON BOM.art_nr = t3.MATNR
	AND t3.SPRAS = 'E';
`

	// Format the query with placeholders
	query = fmt.Sprintf(query, strings.Join(placeholders, ", "))

	return query, args
}

func getNullInt64Value(n sql.NullInt64) interface{} {
	if n.Valid {
		return n.Int64
	}
	return nil // Or use 0 or another default value if you prefer
}
