package mongo

import (
	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"testing"

	"github.com/ONSdigital/dp-search-monitoring/config"
	"github.com/ONSdigital/dp-search-monitoring/importer"
	. "github.com/smartystreets/goconvey/convey"
)

var message = analytics.Message{
	Created:   "Now",
	Url:       "/test/url",
	Term:      "test_term",
	ListType:  "test_list_type",
	GaID:      "testgaID",
	GID:       "testgID",
	PageIndex: 0,
	LinkIndex: 1,
	PageSize:  2,
}

// GetSQSClient returns an instance of analytics.SQSReaderImpl with a mocked SQS client
func GetSQSClient() *analytics.SQSReaderImpl {
	sqsClient := analytics.MockedReceiveMsgs{}

	q := analytics.SQSReaderImpl{
		Client: sqsClient,
		URL:    "http://fake.url",
	}

	return &q
}

// TestMongoClientMock_Insert tests that our mocked MongoClient behaves properly on 'insert'
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

// TestImportSQSMessages tests the import process using our mocked SQS and MongoDB clients.
// Asserts that the single message defined at the top of this file is the only message
// received and 'inserted' into the database.
func TestImportSQSMessages(t *testing.T) {
	if config.SQSDeleteEnabled {
		// Create the mock mongo client
		c := &MongoClientMock{}
		c.InsertFunc = func(message *analytics.Message) error {
			return nil
		}

		// Create the SQSReaderImpl with a mocked SQS client
		q := GetSQSClient()

		// Call ImportSQSMessages and get the total number of messages received/inserted
		count, err := importer.ImportSQSMessages(q, c)

		// Assertions
		Convey("Given a valid SQSReader and MongoClient", t, func() {
			So(len(c.calls.Insert), ShouldEqual, 1)

			So(count, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})
	} else {
		t.Errorf("Deletion of SQS messages must be enabled!")
	}
}
