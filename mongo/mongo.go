package mongo

import (
	"github.com/globalsign/mgo"

	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"github.com/ONSdigital/dp-search-monitoring/config"
	"github.com/ONSdigital/dp-search-monitoring/importer"
	"github.com/ONSdigital/go-ns/log"
)

type MongoClientImpl struct {
	session *mgo.Session
}

func NewMongoClient() (*MongoClientImpl, error) {
	// Create mongo client
	session, err := mgo.Dial(config.MongoDBUrl)
	if err != nil {
		return nil, err
	}
	return &MongoClientImpl{session}, nil
}

func (client *MongoClientImpl) Insert(message *analytics.Message) error {
	c := client.session.DB(config.MongoDBDatabase).C(config.MongoDBCollection)
	return c.Insert(message)
}

func Import() error {
	q, err := analytics.GetReader()

	if err != nil {
		log.Error(err, nil)
		return err
	}

	// Wraps ImportSQSMessages and logs any errors raised
	log.Debug("Starting import.", nil)

	client, err := NewMongoClient()
	if err != nil {
		log.Error(err, nil)
		return err
	}

	count, err := importer.ImportSQSMessages(q, client)
	if err != nil {
		log.Error(err, nil)
		return err
	}

	log.Debug("Insert complete", log.Data{
		"total": count,
	})
	return nil
}
