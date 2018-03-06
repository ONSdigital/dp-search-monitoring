package main

import (
  "net/http"
  "os"
  "strconv"
  "encoding/json"

  "github.com/jasonlvhit/gocron"

  "github.com/ONSdigital/go-ns/log"
  "github.com/ONSdigital/dp-search-monitoring/mongo"
  "github.com/ONSdigital/dp-search-monitoring/config"
  "github.com/ONSdigital/dp-search-monitoring/analytics"
)

const CONTENT_TYPE = "Content-Type"
const APPLICATION_JSON = "application/json"

func attributesHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
    return
  }

  // Get Queue attributes from SQS
  q, err := analytics.GetQueue()
  if err != nil {
    log.Error(err, nil)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  atts, err := q.GetAttributes()
  if err != nil {
    log.Error(err, nil)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Set content header
  w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)

  // Encode the JSON response
  json.NewEncoder(w).Encode(atts)
  return
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
    return
  }
  // Get messages from SQS
  q, err := analytics.GetQueue()
  if err != nil {
    log.Error(err, nil)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  msgs, err := q.GetMessages(config.SQSWaitTimeout, config.MaxSQSMessages)
  if err != nil {
    log.Error(err, nil)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Set content header
  w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)

  // Encode the JSON response
  json.NewEncoder(w).Encode(msgs)
  return
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
    return
  }

  // Trigger an sync between SQS and mongo
  err := mongo.ImportSQSMessages()
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  // Set 200 status OK
  w.WriteHeader(http.StatusOK)
  return
}

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
  s.Start()

  // <- s.Start() // To start scheduler and block

  // Init the server
  mux := http.NewServeMux()

  // Add handlers to the ServeMux
  mux.HandleFunc("/messages", messageHandler)
  mux.HandleFunc("/attributes", attributesHandler)
  mux.HandleFunc("/sync", syncHandler)

  if err := http.ListenAndServe(config.BindAddr, mux); err != nil {
    log.Error(err, nil)
    os.Exit(2)
  }
}
