package config

var Verbose = true  // Enable logging of events

var SuperVerbose = false  // Enable logging of low-priority messages (i.e batch delete responses)

var RunAllOnStartup = true

var TimeUnit = "DAYS"

var TimeWindow uint64 = 1

var AtTime = "00:00"

var Backend = "RDS_POSTGRES"