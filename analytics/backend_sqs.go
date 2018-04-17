package analytics

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"

	"github.com/ONSdigital/dp-search-monitoring/config"
)

//go:generate moq -pkg analytics -out sqs_mocks.go . SQSReader

// SQSReader defines an interface for reading messages from SQS
type SQSReader interface {
	GetAttributes() (*sqs.GetQueueAttributesOutput, error)
	GetMessages(waitTimeout int64, maxNumberOfMessages int64) ([]Message, error)
	DeleteMessage(receiptHandle string) (*sqs.DeleteMessageOutput, error)
	BatchDeleteMessages(receiptHandles []string) (*sqs.DeleteMessageBatchOutput, error)
}

// Queue provides the ability to handle SQS messages.
type SQSReaderImpl struct {
	SQSReader
	Client sqsiface.SQSAPI
	URL    string
}

// Message is a concrete representation of the SQS message
type Message struct {
	Created       string `json:"created"`
	Url           string `json:"url"`
	Term          string `json:"term"`
	ListType      string `json:listType`
	GaID          string `json:gaID`
	GID           string `json:gID`
	PageIndex     int    `json:pageIndex`
	LinkIndex     int    `json:linkIndex`
	PageSize      int    `json:pageSize`
	receiptHandle string `json:receiptHandle`
}

// Set the receipt handle used to delete messages from SQS
func (m *Message) SetReceiptHandle(receiptHandle string) {
	m.receiptHandle = receiptHandle
}

// Get the receipt handle used to delete messages from SQS
func (m *Message) ReceiptHandle() string {
	return m.receiptHandle
}

// Returns a Queue struct for accessing an SQS Queue at a
// specific URL, defined by ANALYTICS_SQS_URL
func GetReader() (*SQSReaderImpl, error) {
	var q SQSReaderImpl

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	q = SQSReaderImpl{
		Client: sqs.New(cfg),
		URL:    config.SQSAnalyticsURL,
	}

	return &q, nil
}

// GetAttributes returns attributes for the desired SQS queue
func (q *SQSReaderImpl) GetAttributes() (*sqs.GetQueueAttributesOutput, error) {
	// Returns the attributes for the desired Queue
	params := sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(q.URL),
		AttributeNames: []sqs.QueueAttributeName{sqs.QueueAttributeNameAll},
	}

	// Generate the request
	req := q.Client.GetQueueAttributesRequest(&params)
	// Send the request
	resp, err := req.Send()

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetMessages returns the parsed messages from SQS if any. If an error
// occurs that error will be returned.
func (q *SQSReaderImpl) GetMessages(waitTimeout int64, maxNumberOfMessages int64) ([]Message, error) {
	params := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.URL),
		MaxNumberOfMessages: &maxNumberOfMessages,
	}

	if waitTimeout > 0 {
		// Poll for messages
		params.WaitTimeSeconds = aws.Int64(waitTimeout)
	}

	// Generate the request
	req := q.Client.ReceiveMessageRequest(&params)
	// Send the request
	resp, err := req.Send()

	if err != nil {
		return nil, fmt.Errorf("failed to get messages, %v", err)
	}

	// Unmarshall the SQS messages
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

// DeleteMessages deletes a single message from SQS using the message receipt handle
func (q *SQSReaderImpl) DeleteMessage(receiptHandle string) (*sqs.DeleteMessageOutput, error) {
	params := sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.URL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	// Generate the request
	req := q.Client.DeleteMessageRequest(&params)
	// Send the request
	resp, err := req.Send()

	if err != nil {
		return nil, fmt.Errorf("failed to delete message with receipt %s, %v", receiptHandle, err)
	}

	return resp, err
}

// BatchDeleteMessages deletes one or many message(s) from SQS using the message receipt handle(s)
func (q *SQSReaderImpl) BatchDeleteMessages(receiptHandles []string) (*sqs.DeleteMessageBatchOutput, error) {
	entries := make([]sqs.DeleteMessageBatchRequestEntry, len(receiptHandles))

	// Generate list of sqs.DeleteMessageBatchRequestEntry using the message receipt handle(s)
	for i, receiptHandle := range receiptHandles {
		entries[i] = sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(strconv.Itoa(i)),
			ReceiptHandle: aws.String(receiptHandle),
		}
	}

	params := sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(q.URL),
	}

	// DeleteMessages deletes a single message from SQS using the message receipt handle
	req := q.Client.DeleteMessageBatchRequest(&params)
	// Send the request
	resp, err := req.Send()

	if err != nil {
		return nil, fmt.Errorf("failed to batch delete messages with receipt, %v", err)
	}

	return resp, err
}
