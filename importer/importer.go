package importer

import (
	"github.com/ONSdigital/go-ns/log"

	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"github.com/ONSdigital/dp-search-monitoring/config"
	"github.com/ONSdigital/dp-search-monitoring/mongo"
	"github.com/ONSdigital/dp-search-monitoring/rds"
)

//go:generate moq -pkg importer -out importer_mocks.go . ImportClient

type ImportClient interface {
	Insert(message *analytics.Message) error
	Close()
}

func getConfiguredImportClient() (ImportClient, error) {
	// Returns the configured import client. Note that if we get here,
	// there must be a valid configuration to pass switch statement in main
	var c ImportClient
	var err error

	switch config.Backend {
	case "MONGO":
		c, err = mongo.New()
		break
	case "RDS_POSTGRES":
		c, err = rds.New()
		break
	}

	return c, err
}

func Import() error {
	q, err := analytics.GetReader()

	if err != nil {
		log.Error(err, nil)
		return err
	}

	c, err := getConfiguredImportClient()

	if err != nil {
		log.Error(err, nil)
		return err
	}

	defer c.Close()

	// Wraps ImportSQSMessages and logs any errors raised
	log.Debug("Starting import.", nil)

	count, err := ImportSQSMessages(q, c)
	if err != nil {
		log.Error(err, nil)
		return err
	}

	log.Debug("Insert complete", log.Data{
		"total": count,
	})
	return nil
}

func ImportSQSMessages(q analytics.SQSReader, c ImportClient) (int64, error) {
	// Pulls a list of SQS messages, stores them, and removes from the the queue
	var count int64

	// Get initial batch of messages
	msgs, err := q.GetMessages(config.SQSWaitTimeout, config.MaxSQSMessages)
	if err != nil {
		return count, err
	}

	// Loop over messages
	for ok := len(msgs) > 0; ok; {
		// Loop through the messages and insert into database
		receiptHandles := make([]string, len(msgs))
		for i := range msgs {
			msg := msgs[i]

			// Do the insert
			err := c.Insert(&msg)

			if err != nil {
				return count, err
			} else {
				// Inserted document, can remove from SQS
				receiptHandles[i] = msg.ReceiptHandle()
			}
			count++
		}

		if config.Verbose {
			log.Debug("Insert progress:", log.Data{
				"total":       count,
				"messageSize": len(msgs),
			})
		}

		if config.SQSDeleteEnabled {
			// Batch delete processed messages
			resp, err := q.BatchDeleteMessages(receiptHandles)

			if err != nil {
				return count, err
			}

			if config.Verbose && config.SuperVerbose {
				log.Debug("Got BatchDeleteResponse", log.Data{
					"response": resp,
				})
			}
		} else if config.Verbose {
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
