package analytics

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"
)

// MockedReceiveMsgs mocks the reading of SQS messages for testing
type MockedReceiveMsgs struct {
	sqsiface.SQSAPI
}

var messages = []Message{
	{
		Created:   "Now",
		Url:       "/test/url",
		Term:      "test_term",
		ListType:  "test_list_type",
		GaID:      "testgaID",
		GID:       "testgID",
		PageIndex: 0,
		LinkIndex: 1,
		PageSize:  2,
	},
}

// Mock the ReceiveMessageRequest method to return a single message
func (m MockedReceiveMsgs) ReceiveMessageRequest(in *sqs.ReceiveMessageInput) sqs.ReceiveMessageRequest {
	// Only need to return mocked response output

	var output sqs.ReceiveMessageOutput

	if len(messages) > 0 {
		// Message hasn't been 'deleted', so add it to the response
		message := messages[0]

		// Marshall the message to JSON
		body, err := json.Marshal(message)

		if err != nil {
			panic(err)
		}

		// Build the output message
		output = sqs.ReceiveMessageOutput{
			Messages: []sqs.Message{
				{
					Body:          aws.String(string(body)),
					ReceiptHandle: aws.String("testHandle"),
				},
			},
		}
	} else {
		// Message has been 'deleted', so return an empty slice
		output = sqs.ReceiveMessageOutput{
			Messages: []sqs.Message{},
		}
	}
	// Build and return the ReceiveMessageRequest
	return sqs.ReceiveMessageRequest{
		Request: &aws.Request{
			Data: &output,
		},
	}
}

// Mocks the DeleteMessageRequest method to set messages to nil (simulating the deletion of our only message)
func (m MockedReceiveMsgs) DeleteMessageRequest(in *sqs.DeleteMessageInput) sqs.DeleteMessageRequest {
	messages = nil
	return sqs.DeleteMessageRequest{
		Request: &aws.Request{
			Data: &sqs.DeleteMessageOutput{},
		},
	}
}

// Mocks the DeleteMessageBatchRequest method to set messages to nil (simulating the deletion of all messages)
func (m MockedReceiveMsgs) DeleteMessageBatchRequest(in *sqs.DeleteMessageBatchInput) sqs.DeleteMessageBatchRequest {
	messages = nil
	return sqs.DeleteMessageBatchRequest{
		Request: &aws.Request{
			Data: &sqs.DeleteMessageBatchOutput{},
		},
	}
}
