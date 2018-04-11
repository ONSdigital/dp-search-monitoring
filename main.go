package main

import (
  "os"
  "strconv"

  "github.com/jasonlvhit/gocron"

  "github.com/ONSdigital/go-ns/log"
  "github.com/ONSdigital/dp-search-monitoring/mongo"
  "github.com/ONSdigital/dp-search-monitoring/config"
)

func main() {
  if v := os.Getenv("ANALYTICS_SQS_URL"); len(v) > 0 {
    config.SQSAnalyticsURL = v
  }

  if v := os.Getenv("MONGODB_URL"); len(v) > 0 {
    config.MongoDBUrl = v
  }

  if v := os.Getenv("SQS_WAIT_TIMEOUT"); len(v) > 0 {
    a, _ := strconv.Atoi(v)
    if a < 0 || a > 20 {
      log.Debug("SQS_WAIT_TIMEOUT must be between 0 and 20", nil)
      os.Exit(1)
    }
    config.SQSWaitTimeout = int64(a)
  }

  // Setup cron job to poll for SQS messages and insert into mongoDB
  s := gocron.NewScheduler()
  s.Every(1).Day().At("00:00").Do(mongo.ImportSQSMessages)

  _, time := gocron.NextRun()
  log.Debug("Cron job scheduled", log.Data{
    "NextRun:":   time,
  })
  // s.Start()

  <- s.Start() // To start scheduler and block
}
