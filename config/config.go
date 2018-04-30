package config

var SQSAnalyticsURL = ""

var SQSWaitTimeout int64 = 20

var MaxSQSMessages int64 = 10

var SQSDeleteEnabled = false

var MongoDBUrl = "localhost:27017"

var MongoDBDatabase = "local"

var MongoDBCollection = "searchstats"

var RunAllOnStartup = true

var TimeUnit = "DAYS"

var TimeWindow uint64 = 1

var AtTime = "00:00"

var Backend = "RDS_POSTGRES"

// RDS

var RdsDbUser = ""

var RdsDbPassword = ""

var RdsDbName = ""

var RdsDbEndpoint = ""

var RdsDbPort int64 = 5432  // Default postgres port