package mongo

import (
	"github.com/globalsign/mgo"

	"github.com/ONSdigital/dp-search-monitoring/analytics"
	"github.com/ONSdigital/dp-search-monitoring/config"
)

type MongoClientImpl struct {
	session *mgo.Session
}

func New() (*MongoClientImpl, error) {
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

func (client *MongoClientImpl) Close() {
	client.session.Close()
}
