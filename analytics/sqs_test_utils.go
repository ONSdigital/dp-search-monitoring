package analytics

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
)

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

func (m MockedReceiveMsgs) ReceiveMessageRequest(in *sqs.ReceiveMessageInput) sqs.ReceiveMessageRequest {
	// Only need to return mocked response output

	var output sqs.ReceiveMessageOutput
    if len(messages) > 0 {
		message := messages[0]

		body, err := json.Marshal(message)

		if err != nil {
			panic(err)
		}

		output = sqs.ReceiveMessageOutput{
			Messages: []sqs.Message{
				{
					Body:          aws.String(string(body)),
					ReceiptHandle: aws.String("testHandle"),
				},
			},
		}
	} else {
		output = sqs.ReceiveMessageOutput{
			Messages: []sqs.Message{},
		}
	}
	return sqs.ReceiveMessageRequest{
		Request: &aws.Request{
			Data: &output,
		},
	}
}

func (m MockedReceiveMsgs) DeleteMessageRequest(in *sqs.DeleteMessageInput) sqs.DeleteMessageRequest {
	messages = nil
	return sqs.DeleteMessageRequest{
		Request: &aws.Request{
			Data: &sqs.DeleteMessageOutput{},
		},
	}
}

func (m MockedReceiveMsgs) DeleteMessageBatchRequest(in *sqs.DeleteMessageBatchInput) sqs.DeleteMessageBatchRequest {
	messages = nil
	return sqs.DeleteMessageBatchRequest{
		Request: &aws.Request{
			Data: &sqs.DeleteMessageBatchOutput{},
		},
	}
}