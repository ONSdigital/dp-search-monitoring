package analytics

import (
  "testing"

  "github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"

  . "github.com/smartystreets/goconvey/convey"
  "github.com/aws/aws-sdk-go-v2/service/sqs"
  "github.com/aws/aws-sdk-go-v2/aws"
)

type mockedReceiveMsgs struct {
  sqsiface.SQSAPI
  Resp sqs.ReceiveMessageOutput
}

func (m mockedReceiveMsgs) ReceiveMessageRequest(in *sqs.ReceiveMessageInput) sqs.ReceiveMessageRequest {
  // Only need to return mocked response output
  output := sqs.ReceiveMessageOutput{
    Messages: []sqs.Message{
      {
        Body: aws.String(`{"created":"1","url":"test_url","term":"test","listType":"test","gaID":"testgaID","gID":"testgID","pageIndex":0,"linkIndex":0,"pageSize":0}`),
        ReceiptHandle: aws.String("testHandle"),
      },
    },
  }
  return sqs.ReceiveMessageRequest{
    Request: &aws.Request{
      Data: &output,
    },
  }
}

func TestSQSReaderImpl_GetMessages(t *testing.T) {
  client := mockedReceiveMsgs{}

  q := SQSReaderImpl{
    Client: client,
    URL: "http://fake.url",
  }

  messages, err := q.GetMessages(20, 10)

  Convey("Given valid input parameters", t, func() {
    So(q, ShouldNotBeNil)
    So(messages, ShouldNotBeNil)
    So(err, ShouldBeNil)
    So(len(messages), ShouldEqual, 1)
  })
}
