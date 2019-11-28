package repos

import (
	"log"

	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepository struct {
	DbConfig   *models.DatabaseConfig
	Collection *mongo.Collection
}

func NewEventRepository(dbConfig *models.DatabaseConfig) *EventRepository {
	return &EventRepository{
		DbConfig:   dbConfig,
		Collection: dbConfig.MongoClient.Collection("event"),
	}
}

func (e *EventRepository) FindByID(id primitive.ObjectID) (*models.Event, error) {

	res := e.Collection.FindOne(e.DbConfig.Ctx, bson.M{"_id": id})
	if res.Err() != nil {
		log.Println("There was an error: ", res.Err())
		return nil, res.Err()
	}

	event := &models.Event{}
	err := res.Decode(&event)
	if err != nil {
		log.Println("There was an error decoding: ", err)
		return nil, err
	}

	return event, nil
}

func (e *EventRepository) Find(query bson.M) ([]*models.Event, error) {

	events := make([]*models.Event, 0)

	cursor, err := e.Collection.Find(e.DbConfig.Ctx, query)
	if err != nil {
		log.Println("There was an error: ", err)
		return nil, err
	}

	for cursor.Next(e.DbConfig.Ctx) {
		event := &models.Event{}
		err = cursor.Decode(&event)
		if err != nil {
			log.Println("Error decoding event: ", err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (e *EventRepository) Insert(event *models.Event) (primitive.ObjectID, error) {

	// Generate ID
	event.ID = primitive.NewObjectID()

	res, err := e.Collection.InsertOne(e.DbConfig.Ctx, event)
	if err != nil {
		log.Println("There was an error inserting eevent: ", err)
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (e *EventRepository) Update(id primitive.ObjectID, event *models.Event) error {

	_, err := e.Collection.UpdateOne(
		e.DbConfig.Ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": event,
		},
	)
	if err != nil {
		log.Println("There was an error updating event: ", err)
		return err
	}

	return nil
}
