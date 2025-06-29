-- Create DETPLAN database if not exists
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'DETPLAN')
BEGIN
    CREATE DATABASE DETPLAN;
END
GO

USE DETPLAN;
GO

-- Create login if it doesn't exist
IF NOT EXISTS (SELECT * FROM sys.server_principals WHERE name = 'myuser@domain.com')
BEGIN
    CREATE LOGIN [myuser@domain.com] WITH PASSWORD = 'mypassword';
END
GO

-- Create user for existing login if login exists and user does not exist
IF EXISTS (SELECT * FROM sys.server_principals WHERE name = 'myuser@domain.com')
BEGIN
    IF NOT EXISTS (SELECT * FROM sys.database_principals WHERE name = 'myuser')
    BEGIN
        CREATE USER [myuser] FOR LOGIN [myuser@domain.com];
        EXEC sp_addrolemember 'db_datareader', [myuser];
        EXEC sp_addrolemember 'db_datawriter', [myuser];
    END
END
ELSE
BEGIN
    PRINT 'Login myuser@domain.com does not exist. Please create the login first.';
END
GO

-- Create table SplitStandardMaterialPurchaseText if not exists
IF NOT EXISTS (
    SELECT * FROM sys.tables WHERE name = 'SplitStandardMaterialPurchaseText' AND schema_id = SCHEMA_ID('dbo')
)
BEGIN
CREATE TABLE [dbo].[SplitStandardMaterialPurchaseText] (
    Artnr NVARCHAR(50) NOT NULL,
    ISactive BIT NOT NULL,
    QTY INT NULL,
    Value DECIMAL(18, 2) NULL,
    PurchaseText NVARCHAR(255) NULL,
    ChangedAt DATETIME NULL,
    [Ref Designator] NVARCHAR(255) NULL,
    ManufacturerCode NVARCHAR(50) NULL
);
END
GO

-- Insert sample data if not already there
IF NOT EXISTS (SELECT 1 FROM [dbo].[SplitStandardMaterialPurchaseText] WHERE Artnr = 'ArtNr1')
BEGIN
    INSERT INTO [dbo].[SplitStandardMaterialPurchaseText] 
    (Artnr, ISactive, QTY, Value, PurchaseText, ChangedAt, [Ref Designator], ManufacturerCode)
    VALUES ('ArtNr1', 1, 10, 100.50, 'Purchase text for ArtNr1', GETDATE(), 'RD1', 'Manuf1');
END
GO

IF NOT EXISTS (SELECT 1 FROM [dbo].[SplitStandardMaterialPurchaseText] WHERE Artnr = 'ArtNr2')
BEGIN
    INSERT INTO [dbo].[SplitStandardMaterialPurchaseText] 
    (Artnr, ISactive, QTY, Value, PurchaseText, ChangedAt, [Ref Designator], ManufacturerCode)
    VALUES ('ArtNr2', 1, 20, 200.75, 'Purchase text for ArtNr2', GETDATE(), 'RD2', 'Manuf2');
END
GO

-- Create MSupply database if not exists
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'MSupply')
BEGIN
    CREATE DATABASE MSupply;
END
GO

USE MSupply;
GO

-- Create table StandardManufacturer if not exists
IF NOT EXISTS (
    SELECT * FROM sys.tables WHERE name = 'StandardManufacturer' AND schema_id = SCHEMA_ID('dbo')
)
BEGIN
CREATE TABLE [dbo].[StandardManufacturer] (
    ManufacturerCode NVARCHAR(50) NOT NULL PRIMARY KEY,
    Name NVARCHAR(255) NULL
);
END
GO

-- Insert sample data if not exists
IF NOT EXISTS (SELECT 1 FROM [dbo].[StandardManufacturer] WHERE ManufacturerCode = 'Manuf1')
BEGIN
    INSERT INTO [dbo].[StandardManufacturer] (ManufacturerCode, Name)
    VALUES ('Manuf1', 'Manufacturer 1');
END
GO

IF NOT EXISTS (SELECT 1 FROM [dbo].[StandardManufacturer] WHERE ManufacturerCode = 'Manuf2')
BEGIN
    INSERT INTO [dbo].[StandardManufacturer] (ManufacturerCode, Name)
    VALUES ('Manuf2', 'Manufacturer 2');
END
GO

-- Create table StandardMaterial if not exists
IF NOT EXISTS (
    SELECT * FROM sys.tables WHERE name = 'StandardMaterial' AND schema_id = SCHEMA_ID('dbo')
)
BEGIN
CREATE TABLE [dbo].[StandardMaterial] (
    Artnr NVARCHAR(50) NOT NULL PRIMARY KEY,
    Designation NVARCHAR(255) NULL,
    Dimension NVARCHAR(255) NULL
);
END
GO

-- Insert sample data if not exists
IF NOT EXISTS (SELECT 1 FROM [dbo].[StandardMaterial] WHERE Artnr = 'ArtNr1')
BEGIN
    INSERT INTO [dbo].[StandardMaterial] (Artnr, Designation, Dimension)
    VALUES ('ArtNr1', 'Designation 1', 'Dim1');
END
GO

IF NOT EXISTS (SELECT 1 FROM [dbo].[StandardMaterial] WHERE Artnr = 'ArtNr2')
BEGIN
    INSERT INTO [dbo].[StandardMaterial] (Artnr, Designation, Dimension)
    VALUES ('ArtNr2', 'Designation 2', 'Dim2');
END
GO

-- Create SAP database if not exists
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'SAP')
BEGIN
    CREATE DATABASE SAP;
END
GO

USE SAP;
GO

-- Create table MaterialText_MAKT if not exists
IF NOT EXISTS (
    SELECT * FROM sys.tables WHERE name = 'MaterialText_MAKT' AND schema_id = SCHEMA_ID('dbo')
)
BEGIN
CREATE TABLE [dbo].[MaterialText_MAKT] (
    MATNR NVARCHAR(50) NOT NULL,
    SPRAS NVARCHAR(2) NOT NULL,
    MAKTX NVARCHAR(255) NULL,
    CONSTRAINT PK_MaterialText_MAKT PRIMARY KEY (MATNR, SPRAS)
);
END
GO

-- Insert sample data if not exists
IF NOT EXISTS (SELECT 1 FROM [dbo].[MaterialText_MAKT] WHERE MATNR = 'ArtNr1' AND SPRAS = 'E')
BEGIN
    INSERT INTO [dbo].[MaterialText_MAKT] (MATNR, SPRAS, MAKTX)
    VALUES ('ArtNr1', 'E', 'Material description for ArtNr1');
END
GO

IF NOT EXISTS (SELECT 1 FROM [dbo].[MaterialText_MAKT] WHERE MATNR = 'ArtNr2' AND SPRAS = 'E')
BEGIN
    INSERT INTO [dbo].[MaterialText_MAKT] (MATNR, SPRAS, MAKTX)
    VALUES ('ArtNr2', 'E', 'Material description for ArtNr2');
END
GO
