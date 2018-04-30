package rds

import (
	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"github.com/ONSdigital/go-ns/log"

	_ "database/sql/driver"
	"database/sql"
)

type RdsSQLClient struct {
	db *sql.DB
}

func New() (*RdsSQLClient, error) {
	db, err := analytics.MySQLDriver()
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

	stmt, err := client.db.Prepare(
		"INSERT INTO messages (created, url, term, listType, gaID, gID, pageIndex, linkIndex, pageSize) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		message.Created, message.Url, message.Term, message.ListType, message.GaID, message.GID, message.PageIndex, message.LinkIndex, message.PageSize)

	if err != nil {
		return err
	}

	log.Debug("Got message", log.Data{
		"Created": message.Created,
		"Url": message.Url,
		"Term": message.Term,
		"ListType": message.ListType,
		"GaID": message.GaID,
		"GID": message.GID,
		"PageIndex": message.PageIndex,
		"LinkIndex": message.LinkIndex,
		"PageSize": message.PageSize,
		"ReceiptHandle": message.ReceiptHandle(),
	})

	return nil
}

func (client *RdsSQLClient) Close() {
	client.db.Close()
}