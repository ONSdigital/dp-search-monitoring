package main

import (
	"os"
	"strconv"

	"github.com/jasonlvhit/gocron"

	"github.com/ONSdigital/dp-search-monitoring/config"
	"github.com/ONSdigital/dp-search-monitoring/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func main() {
	// Setup config
	if v := os.Getenv("ANALYTICS_SQS_URL"); len(v) > 0 {
		config.SQSAnalyticsURL = v
	}

	if v := os.Getenv("SQS_WAIT_TIMEOUT"); len(v) > 0 {
		a, _ := strconv.Atoi(v)
		if a < 0 || a > 20 {
			log.Debug("SQS_WAIT_TIMEOUT must be between 0 and 20", nil)
			os.Exit(1)
		}
		config.SQSWaitTimeout = int64(a)
	}

	if v := os.Getenv("MONGODB_URL"); len(v) > 0 {
		config.MongoDBUrl = v
	}

	if v := os.Getenv("MONGO_DB"); len(v) > 0 {
		config.MongoDBDatabase = v
	}

	if v := os.Getenv("MONGO_COLLECTION"); len(v) > 0 {
		config.MongoDBCollection = v
	}

	if v := os.Getenv("RUN_ON_STARTUP"); len(v) > 0 {
		val, err := strconv.ParseBool(v)

		if err != nil {
			log.Debug("Unable to convert 'RUN_ON_STARTUP' val to bool", log.Data{
				"Value": val,
			})
			os.Exit(1)
		}
		config.RunAllOnStartup = val
	}

	if v := os.Getenv("TIME_UNIT"); len(v) > 0 {
		config.TimeUnit = v
	}

	if v := os.Getenv("AT_TIME"); len(v) > 0 {
		config.AtTime = v
	}

	// Setup cron job to poll for SQS messages and insert into mongoDB
	s := gocron.NewScheduler()

	// Schedule import by specified TimeUnit
	switch config.TimeUnit {
	case "DAYS":
		s.Every(1).Day().At(config.AtTime).Do(mongo.Import)
		break
	case "HOURS":
		s.Every(1).Hour().Do(mongo.Import)
		break
	default:
		log.Debug("Unknown 'TIME_UNIT'.", log.Data{
			"TimeUnit": config.TimeUnit,
		})
		os.Exit(1)
	}

	// Run import on initial startup?
	if config.RunAllOnStartup {
		s.RunAll()
	}

	// Log the time of the next run
	_, time := s.NextRun()
	log.Debug("Cron job scheduled", log.Data{
		"NextRun": time,
	})

	<-s.Start() // Start scheduler and block
}
