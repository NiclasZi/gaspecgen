DECLARE @BOM_List TABLE (
    [QTY] INT,
    [Art_nr] NVARCHAR(255),
    [Ref_Designator] TEXT
);

INSERT INTO @BOM_List (
    [QTY],
    [Art_nr],
    [Ref_Designator]
)
VALUES
{{- $length := len .Rows }}
{{- range $index, $row := .Rows }}
    ({{ $row.qty }}, '{{ $row.artNr }}', '{{ $row.refDesignator }}'){{ if lt $index (sub $length 1) }},{{ end }}
{{- end }};

SELECT DISTINCT
    BOM.[Art_nr],
    BOM.[QTY],
    CASE
        WHEN Std.designation IS NOT NULL THEN Std.designation
        ELSE t3.MAKTX
    END AS [Description],
    Std.[Dimension],
    SM.[Name] AS ManufacturerName,
    PT.[PurchaseText]
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
