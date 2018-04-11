package mongo

import (
  "gopkg.in/mgo.v2"

  "github.com/ONSdigital/go-ns/log"
  "github.com/ONSdigital/dp-search-monitoring/config"
  "github.com/ONSdigital/dp-search-monitoring/analytics"
)

func Import() {
  q, err := analytics.GetReader()

  if err != nil {
    log.Error(err, nil)
    return
  }

  // Wraps ImportSQSMessages and logs any errors raised
  log.Debug("Starting import.", nil)

  err = importSQSMessages(q)
  if err != nil {
    log.Error(err, nil)
  }
  log.Debug("Import complete.", nil)
}

func importSQSMessages(q analytics.SQSReader) error {
  // Reads messages from SQS queue and impors them into mongoDB

  // Create mongo client
  session, err := mgo.Dial(config.MongoDBUrl)
  if err != nil {
    return err
  }

  // Load the collection
  c := session.DB(config.MongoDBDatabase).C(config.MongoDBCollection)

  // Defer closing of the mongoDB session
  defer session.Close()

  // Check for more messages
  msgs, err := q.GetMessages(config.SQSWaitTimeout, config.MaxSQSMessages)
  if err != nil {
    return err
  }

  var count int64

  for ok := (len(msgs) > 0); ok; {
    // Loop through the messages and insert into mongo
    receiptHandles := make([]string, len(msgs))
    for i := range msgs {
      msg := msgs[i]

      err := c.Insert(&msg)

      if err != nil {
        return err
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

    // Batch delete processed messages
    q.BatchDeleteMessages(receiptHandles)

    // Check for more messages
    msgs, err = q.GetMessages(config.SQSWaitTimeout, config.MaxSQSMessages)
    if err != nil {
      return err
    }
    ok = (len(msgs) > 0)
  }

  log.Debug("Insert complete", log.Data{
    "total": count,
  })

  return nil
}
