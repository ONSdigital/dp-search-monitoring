package mongo

import (
	"testing"
	"github.com/ONSdigital/dp-search-monitoring/analytics"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ONSdigital/dp-search-monitoring/importer"
	"github.com/ONSdigital/dp-search-monitoring/config"
)

var message = analytics.Message{
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

func GetSQSClient() (*analytics.SQSReaderImpl) {
	sqsClient := analytics.MockedReceiveMsgs{}

	q := analytics.SQSReaderImpl{
		Client: sqsClient,
		URL: "http://fake.url",
	}

	return &q
}

func TestMongoClientMock_Insert(t *testing.T) {
	c := &MongoClientMock{}
	c.InsertFunc = func(message *analytics.Message) error {
		return nil
	}

	Convey("Given valid input message", t, func() {
		err := c.Insert(&message)

		So(len(c.calls.Insert), ShouldEqual, 1)
		So(c.calls.Insert[0].Message, ShouldResemble, &message)

		So(err, ShouldBeNil)
	})
}

func TestPullMessages(t *testing.T) {
	if config.SQSDeleteEnabled {
		c := &MongoClientMock{}
		c.InsertFunc = func(message *analytics.Message) error {
			return nil
		}

		q := GetSQSClient()

		count, err := importer.ImportSQSMessages(q, c)

		Convey("Given a valid SQSReader and MongoClient", t, func() {
			So(len(c.calls.Insert), ShouldEqual, 1)

			So(count, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})
	} else {
		t.Errorf("Deletion of SQS messages must be enabled!")
	}
}
