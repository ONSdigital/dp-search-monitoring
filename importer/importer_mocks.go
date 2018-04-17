// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package importer

import (
	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"sync"
)

var (
	lockImportClientMockInsert sync.RWMutex
)

// ImportClientMock is a mock implementation of ImportClient.
//
//     func TestSomethingThatUsesImportClient(t *testing.T) {
//
//         // make and configure a mocked ImportClient
//         mockedImportClient := &ImportClientMock{
//             InsertFunc: func(message *analytics.Message) error {
// 	               panic("TODO: mock out the Insert method")
//             },
//         }
//
//         // TODO: use mockedImportClient in code that requires ImportClient
//         //       and then make assertions.
//
//     }
type ImportClientMock struct {
	// InsertFunc mocks the Insert method.
	InsertFunc func(message *analytics.Message) error

	// calls tracks calls to the methods.
	calls struct {
		// Insert holds details about calls to the Insert method.
		Insert []struct {
			// Message is the message argument value.
			Message *analytics.Message
		}
	}
}

// Insert calls InsertFunc.
func (mock *ImportClientMock) Insert(message *analytics.Message) error {
	if mock.InsertFunc == nil {
		panic("moq: ImportClientMock.InsertFunc is nil but ImportClient.Insert was just called")
	}
	callInfo := struct {
		Message *analytics.Message
	}{
		Message: message,
	}
	lockImportClientMockInsert.Lock()
	mock.calls.Insert = append(mock.calls.Insert, callInfo)
	lockImportClientMockInsert.Unlock()
	return mock.InsertFunc(message)
}

// InsertCalls gets all the calls that were made to Insert.
// Check the length with:
//     len(mockedImportClient.InsertCalls())
func (mock *ImportClientMock) InsertCalls() []struct {
	Message *analytics.Message
} {
	var calls []struct {
		Message *analytics.Message
	}
	lockImportClientMockInsert.RLock()
	calls = mock.calls.Insert
	lockImportClientMockInsert.RUnlock()
	return calls
}
