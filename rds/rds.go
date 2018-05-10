package rds

import (
	"database/sql"

	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"fmt"
	"github.com/ONSdigital/dp-search-monitoring/config"
)

type RdsSQLClient struct {
	db *sql.DB
}

func New() (*RdsSQLClient, error) {
	db, err := MySQLDriver()
	if err != nil {
		return nil, err
	}
	return &RdsSQLClient{db}, err
}

//CREATE TABLE messages (
//created date NOT NULL,
//url text NOT NULL,
//term text NOT NULL,
//listType text NOT NULL,
//gaID text NOT NULL,
//gID text NOT NULL,
//pageIndex integer NOT NULL,
//linkIndex integer NOT NULL,
//pageSize integer NOT NULL
//)

func (client *RdsSQLClient) Insert(message *analytics.Message) error {

	query := fmt.Sprintf("INSERT INTO %s (created, url, term, listType, gaID, gID, pageIndex, linkIndex, pageSize) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", config.RdsDbTable)

	stmt, err := client.db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		message.Created, message.Url, message.Term, message.ListType, message.GaID, message.GID, message.PageIndex, message.LinkIndex, message.PageSize)

	if err != nil {
		return err
	}

	return nil
}

func (client *RdsSQLClient) Close() {
	client.db.Close()
}