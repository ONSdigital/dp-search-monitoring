package analytics

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/rds/rdsutils"
	"github.com/ONSdigital/dp-search-monitoring/config"

	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func MySQLDriver() (*mysql.MySQLDriver, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	// Init RDS client
	r := rds.New(cfg)

	// Get credentials
	c := r.Credentials

	// Build auth token
	authToken, err := rdsutils.BuildAuthToken(config.RdsDbEndpoint, cfg.Region, config.RdsDbUser, c)
	if err != nil {
		return nil, err
	}

	// Create the MySQL DNS string for the DB connection
	// user:password@protocol(endpoint)/dbname?<params>
	dnsStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true",
		config.RdsDbUser, authToken, config.RdsDbEndpoint, config.RdsDbName,
	)

	driver := mysql.MySQLDriver{}
	// Use db to perform SQL operations on database
	if _, err = sql.Open("mysql", dnsStr); err != nil {
		panic(err)
	}

	return &driver, err
}