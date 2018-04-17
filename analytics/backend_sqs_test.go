package analytics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestSQSReaderImpl_GetMessages tests that the SQSReaderImpl struct correctly interfaces with the (mocked) SQS
// client
func TestSQSReaderImpl_GetMessages(t *testing.T) {
	// Get a handle on the Mocked SQS client interface
	client := MockedReceiveMsgs{}

	// Initialise our SQSReaderImpl with a fake URL
	q := SQSReaderImpl{
		Client: client,
		URL:    "http://fake.url",
	}

	// Get the list of messages (should have length of 1)
	messages, err := q.GetMessages(20, 10)
	message := messages[0]

	// Assertions
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
