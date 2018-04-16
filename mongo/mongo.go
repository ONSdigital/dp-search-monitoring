package mongo

import (
  "gopkg.in/mgo.v2"

  "github.com/ONSdigital/go-ns/log"
  "github.com/ONSdigital/dp-search-monitoring/config"
  "github.com/ONSdigital/dp-search-monitoring/analytics"
)

//go:generate moq -pkg mongo -out mongo_mocks.go . MongoClient

type MongoClient interface {
  Insert(message *analytics.Message) error
}

type MongoClientImpl struct {
  session *mgo.Session
}

func NewMongoClient() (*MongoClientImpl, error) {
  // Create mongo client
  session, err := mgo.Dial(config.MongoDBUrl)
  if err != nil {
    return nil, err
  }
  return &MongoClientImpl{session}, nil
}

func (client *MongoClientImpl) Insert(message *analytics.Message) error {
  c := client.session.DB(config.MongoDBDatabase).C(config.MongoDBCollection)
  return c.Insert(message)
}

func PullMessages(q analytics.SQSReader, client MongoClient) (int64, error) {
  // Pulls a list of SQS messages into mongo and removes from the the queue
  var count int64

  // Get initial batch of messages
  msgs, err := q.GetMessages(config.SQSWaitTimeout, config.MaxSQSMessages)
  if err != nil {
    return count, err
  }

  for ok := len(msgs) > 0; ok; {
    // Loop through the messages and insert into mongo
    receiptHandles := make([]string, len(msgs))
    for i := range msgs {
      msg := msgs[i]

      err := client.Insert(&msg)

      if err != nil {
        return count, err
      } else {
        // Inserted document, can remove from SQS
        receiptHandles[i] = msg.ReceiptHandle()
      }
      count++
    }
    log.Debug("Insert progress:", log.Data{
      "total": count,
      "messageSize": len(msgs),
    })

    if config.SQSDeleteEnabled {
      // Batch delete processed messages
      resp, err := q.BatchDeleteMessages(receiptHandles)

      if err != nil {
        return count, err
      }

      log.Debug("Got BatchDeleteResponse", log.Data{
        "response": resp,
      })
    } else {
      log.Debug("Currently configured to prevent deletion of SQS messages", nil)
    }

    if err != nil {
      return count, err
    }

    // Check for more messages
    msgs, err = q.GetMessages(config.SQSWaitTimeout, config.MaxSQSMessages)
    if err != nil {
      return count, err
    }
    ok = len(msgs) > 0
  }

  return count, nil
}

func Import() {
  q, err := analytics.GetReader()

  if err != nil {
    log.Error(err, nil)
    return
  }

  // Wraps ImportSQSMessages and logs any errors raised
  log.Debug("Starting import.", nil)

  client, err := NewMongoClient()
  if err != nil {
    log.Error(err, nil)
  }

  count, err := PullMessages(q, client)
  if err != nil {
    log.Error(err, nil)
  }

  log.Debug("Insert complete", log.Data{
    "total": count,
  })
}
