package analytics

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"

  "github.com/ONSdigital/go-ns/log"
  "github.com/ONSdigital/dp-search-monitoring/config"
)

// Queue provides the ability to handle SQS messages.
type Queue struct {
	Client sqsiface.SQSAPI
	URL    string
}

// Message is a concrete representation of the SQS message
type Message struct {
	Created string `json:"created"`
	Url   string `json:"url"`
	Term  string `json:"term"`
  ListType string `json:listType`
  GaID string `json:gaID`
  PageIndex int `json:pageIndex`
  LinkIndex int `json:linkIndex`
  PageSize int `json:pageSize`
	receiptHandle string `json:receiptHandle`
}

func (m *Message) SetReceiptHandle(receiptHandle string) {
	m.receiptHandle = receiptHandle
}

func(m *Message) ReceiptHandle() string {
	return m.receiptHandle
}

// Returns a Queue struct for accessing an SQS Queue at a
// specific URL, defined by ANALYTICS_SQS_URL
func GetQueue() (Queue, error) {
  var q Queue

  cfg, err := external.LoadDefaultAWSConfig()
  if err != nil {
    return q, err
  }

  q = Queue{
    Client: sqs.New(cfg),
    URL:    config.SQSAnalyticsURL,
  }

  return q, nil
}

// GetAttributes returns attributes for the desired SQS queue
func (q *Queue) GetAttributes() (*sqs.GetQueueAttributesOutput, error) {
	// Returns the attributes for the desired Queue
	params := sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(q.URL),
		AttributeNames: []sqs.QueueAttributeName{sqs.QueueAttributeNameAll},
	}

	req := q.Client.GetQueueAttributesRequest(&params)
	resp, err := req.Send()

	log.Debug("Got attributes response", log.Data{
		"resp": resp,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetMessages returns the parsed messages from SQS if any. If an error
// occurs that error will be returned.
func (q *Queue) GetMessages(waitTimeout int64, maxNumberOfMessages int64) ([]Message, error) {
	params := sqs.ReceiveMessageInput{
		QueueUrl: aws.String(q.URL),
		MaxNumberOfMessages: &maxNumberOfMessages,
	}
	if waitTimeout > 0 {
		params.WaitTimeSeconds = aws.Int64(waitTimeout)
	}
	req := q.Client.ReceiveMessageRequest(&params)
	resp, err := req.Send()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages, %v", err)
	}

  log.Debug("Got message response", log.Data{
    "resp": resp,
		"size": len(resp.Messages),
  })

	msgs := make([]Message, len(resp.Messages))
	for i, msg := range resp.Messages {
		parsedMsg := Message{}
		if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message, %v", err)
		}

    // Add the ReceiptHandle
    parsedMsg.SetReceiptHandle(*msg.ReceiptHandle)
		msgs[i] = parsedMsg
	}

	return msgs, nil
}

func (q *Queue) DeleteMessage(receiptHandle string) (*sqs.DeleteMessageOutput, error) {
	params := sqs.DeleteMessageInput{
		QueueUrl: aws.String(q.URL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	req := q.Client.DeleteMessageRequest(&params)
	resp, err := req.Send()

	if err != nil {
		return nil, fmt.Errorf("failed to delete message with receipt %s, %v", receiptHandle, err)
	}

	log.Debug("Got delete response", log.Data{
		"resp": resp,
	})

	return resp, err
}

func (q *Queue) BatchDeleteMessages(receiptHandles []string) (*sqs.DeleteMessageBatchOutput, error) {
	entries := make([]sqs.DeleteMessageBatchRequestEntry, len(receiptHandles))

	for i, receiptHandle := range receiptHandles {
		entries[i] = sqs.DeleteMessageBatchRequestEntry{
			Id: aws.String(strconv.Itoa(i)),
			ReceiptHandle: aws.String(receiptHandle),
		}
	}

	params := sqs.DeleteMessageBatchInput{
		Entries: entries,
		QueueUrl: aws.String(q.URL),
	}

	req := q.Client.DeleteMessageBatchRequest(&params)
	resp, err := req.Send()

	if err != nil {
		return nil, fmt.Errorf("failed to batch delete messages with receipt, %v", err)
	}

	log.Debug("Got batch delete response", log.Data{
		"resp": resp,
	})

	return resp, err
}
