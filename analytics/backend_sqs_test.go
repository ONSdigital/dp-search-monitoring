package analytics

import (
  "errors"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
)

func TestGetMessages(t *testing.T) {
  q := &SQSReaderMock{}
  q.GetMessagesFunc = func(waitTimeout int64, maxNumberOfMessages int64) ([]Message, error) {
    return nil, errors.New("Unable to get messages")
  }

  Convey("Given valid input parameters", t, func() {
    messages, err := q.GetMessages(20, 10)

    So(len(q.calls.GetMessages), ShouldEqual, 1)
    So(q.calls.GetMessages[0].WaitTimeout, ShouldEqual, 20)
    So(q.calls.GetMessages[0].MaxNumberOfMessages, ShouldEqual, 10)

    So(messages, ShouldBeNil)
    So(err, ShouldNotBeNil)
  })
}
