package manager

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/microsoft/go-mssqldb/azuread" // MS SQL driver with Azure AD support
	"github.com/spf13/viper"
)

// DB is the singleton struct for the database connection.
type DB struct {
	connection *sql.DB
}

var (
	instance *DB
	once     sync.Once
)

// Get returns the singleton instance with the connection from config
func Get() *DB {
	options := fmt.Sprintf("encrypt=%s", (func(opt bool) string {
		if opt {
			return "true"
		}
		return "false"
	})(viper.GetBool("db-encrypt")))
	if viper.GetBool("db-trust-cert") {
		options += "&TrustServerCertificate=true"
	}
	if viper.GetString("loglevel") == "debug" {
		// All the logs
		options += "&log=255"
	}
	return GetInstance(
		fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&%s",
			viper.GetString("db-user"),
			viper.GetString("db-password"),
			viper.GetString("db-host"),
			viper.GetString("db"),
			options,
		),
	)
}

// GetInstance returns the singleton instance of the DB struct.
func GetInstance(connString string) *DB {
	once.Do(func() {
		db, err := sql.Open("sqlserver", connString)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// Ping the database to verify the connection
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping the database: %v", err)
		}

		instance = &DB{
			connection: db,
		}
		fmt.Println("Database connection established.")
	})
	return instance
}

// GetConnection provides access to the SQL database connection.
func (d *DB) GetConnection() *sql.DB {
	return d.connection
}

// Close closes the database connection.
func (d *DB) Close() error {
	if d.connection != nil {
		return d.connection.Close()
	}
	return nil
}
