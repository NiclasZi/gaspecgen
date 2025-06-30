DECLARE @BOM_List TABLE (
    [QTY] INT,
    [Art_nr] NVARCHAR(255),
    [Ref_Designator] NVARCHAR(MAX)
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

SELECT distinct [Art_nr], qty
      --,[Description]
	  ,(case 
	when Std.designation IS NOT NULL THEN Std.designation
	--when left(SUBSTRING(t2.EAN11,1,5),1) in ('B','W','C','U') then Std.Ben
	else t3.MAKTX
	end) AS Description
	,std.Dimension
      --,[QTY]   
      --,[Value]
      --,[Cell]
      --,[Part Number]
      --,[SAP_status]
	  ,SM.Name
	  ,PT.PurchaseText
	  --,PT.ChangedAt
	  ,[Ref_Designator]

  FROM 
    @BOM_List BOM
  Join [DETPLAN].[dbo].[SplitStandardMaterialPurchaseText] PT
  on BOM.Art_nr=PT.Artnr
And PT.ISactive=1

Join [MSupply].[dbo].[StandardManufacturer] SM
on SM.[ManufacturerCode]=PT.[ManufacturerCode]

Left Join Msupply.[dbo].[StandardMaterial] std
	on std.Artnr=BOM.art_nr
	left JOIN SAP.dbo.MaterialText_MAKT t3
	ON BOM.art_nr = t3.MATNR
	AND t3.SPRAS = 'E'