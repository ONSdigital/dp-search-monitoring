package rds

import (
	"fmt"

	"github.com/ONSdigital/dp-search-monitoring/config"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ONSdigital/go-ns/log"
)

var createTableStatements = []string{
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		created date NOT NULL,
		url text NOT NULL,
		term text NOT NULL,
		listType text NOT NULL,
		gaID text NOT NULL,
		gID text NOT NULL,
		pageIndex integer NOT NULL,
		linkIndex integer NOT NULL,
		pageSize integer NOT NULL
	)`, config.RdsDbTable),
}

// createTable creates the table
func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
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

	// Create table, if necessary
	err = createTable(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}