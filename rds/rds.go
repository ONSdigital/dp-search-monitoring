package rds

import (
	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type RdsSQLClient struct {
	db *mysql.MySQLDriver
}

func New() (*RdsSQLClient, error) {
	db, err := analytics.MySQLDriver()
	if err != nil {
		return nil, err
	}
	return &RdsSQLClient{db}, err
}

func (client *RdsSQLClient) Insert(message *analytics.Message) error {
	fmt.Println(message)
	return nil
}

func (client *RdsSQLClient) Close() {

}