package analytics

import (
  "testing"
  "encoding/json"

  "github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"
  "github.com/aws/aws-sdk-go-v2/service/sqs"
  "github.com/aws/aws-sdk-go-v2/aws"

  . "github.com/smartystreets/goconvey/convey"
)

type mockedReceiveMsgs struct {
  sqsiface.SQSAPI
}

func (m mockedReceiveMsgs) ReceiveMessageRequest(in *sqs.ReceiveMessageInput) sqs.ReceiveMessageRequest {
  // Only need to return mocked response output
  message := Message{
    Created: "Now",
    Url: "/test/url",
    Term: "test_term",
    ListType: "test_list_type",
    GaID: "testgaID",
    GID: "testgID",
    PageIndex: 0,
    LinkIndex: 1,
    PageSize: 2,
  }

  body, err := json.Marshal(message)

  if err != nil {
    panic(err)
  }

  output := sqs.ReceiveMessageOutput{
    Messages: []sqs.Message{
      {
        Body: aws.String(string(body)),
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
  message := messages[0]

  Convey("Given valid input parameters", t, func() {
    So(q, ShouldNotBeNil)
    So(messages, ShouldNotBeNil)
    So(err, ShouldBeNil)
    So(len(messages), ShouldEqual, 1)

    So(message.Created, ShouldEqual, "Now")
    So(message.Url, ShouldEqual, "/test/url")
    So(message.Term, ShouldEqual, "test_term")
    So(message.ListType, ShouldEqual, "test_list_type")
    So(message.GaID, ShouldEqual, "testgaID")
    So(message.GID, ShouldEqual, "testgID")
    So(message.PageIndex, ShouldEqual, 0)
    So(message.LinkIndex, ShouldEqual, 1)
    So(message.PageSize, ShouldEqual, 2)
    So(message.ReceiptHandle(), ShouldEqual, "testHandle")
  })
}
