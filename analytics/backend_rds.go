package analytics

import (
	"fmt"

	"github.com/ONSdigital/dp-search-monitoring/config"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ONSdigital/go-ns/log"
)

// mysqlDB messages books to a MySQL instance.
type MySqlDB struct {
	Conn *sql.DB
}

func MySQLDriver() (*sql.DB, error) {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
				config.RdsDbUser, config.RdsDbPassword, config.RdsDbName, config.RdsDbEndpoint, config.RdsDbPort)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	log.Debug("Successfully made connection to postgres DB", log.Data{
		"dbName": config.RdsDbName,
	})

	return conn, nil
}