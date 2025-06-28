package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/microsoft/go-mssqldb/azuread"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DB struct {
	connection *sql.DB
}

var (
	instance *DB
	lastErr  error
	once     sync.Once
)

// Get returns the singleton instance with the connection from config
func Get() (*DB, error) {
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
// returns nil on failure
func GetInstance(connString string) (*DB, error) {
	once.Do(func() {
		db, err := sql.Open("sqlserver", connString)
		if err != nil {
			//zap.L().Error("Failed to connect to the database", zap.Error(err))
			lastErr = fmt.Errorf("failed to connect to the database, err: %s", err.Error())
			return
		}

		// Ping the database to verify the connection
		if err := db.Ping(); err != nil {
			//zap.L().Error("Failed to ping the database", zap.Error(err))
			lastErr = fmt.Errorf("failed to ping the database, err: %s", err.Error())
		}

		instance = &DB{
			connection: db,
		}
		zap.L().Debug("Connected to db")
	})
	return instance, lastErr
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
