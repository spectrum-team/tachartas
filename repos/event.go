package repos

import (
	"bytes"
	"encoding/base64"
	"log"

	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

	event.ApiID = event.ID.Hex()

	e.addImageToEvent(event)

	return event, nil
}

func (e *EventRepository) Find(query bson.M) ([]*models.Event, error) {

	events := make([]*models.Event, 0)

	// opts := &options.FindOptions{
	// 	Limit: query.Limit,
	// 	Skip:  query.Skip,
	// }

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

		event.ApiID = event.ID.Hex()

		e.addImageToEvent(event)

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

func (e *EventRepository) Assist(id primitive.ObjectID, willAssist int) error {

	// First I will look for the existing event
	event, err := e.FindByID(id)
	if err != nil {
		log.Println("there was an error looking for existing event: ", err)
		return err
	}

	switch willAssist {
	case 1:
		event.Assistants++
	case 0:
		if event.Assistants > 0 {
			event.Assistants--
		}
	}

	err = e.Update(id, event)
	if err != nil {
		log.Println("Could not update assists on the event: ", err)
		return err
	}

	return nil
}

func (e *EventRepository) AddEventImage(id primitive.ObjectID, imgName string, img []byte) error {

	// First I will look for the existing event
	event, err := e.FindByID(id)
	if err != nil {
		log.Println("there was an error looking for existing event: ", err)
		return err
	}

	log.Println(event)

	bucket, err := gridfs.NewBucket(e.DbConfig.MongoClient, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	uploadStream, err := bucket.OpenUploadStream(imgName)
	if err != nil {
		log.Println("Error creating upload stream: ", err)
		return err
	}

	defer uploadStream.Close()

	_, err = uploadStream.Write(img)
	if err != nil {
		log.Println("There was an error storing the image: ", err)
		return err
	}

	fileId := uploadStream.FileID.(primitive.ObjectID)

	log.Println("The file ID hopefully: ", fileId)

	log.Println("images")
	event.Image = &fileId

	_, err = e.Collection.UpdateOne(
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

func (e *EventRepository) addImageToEvent(event *models.Event) {
	if event.Image != nil {

		bucket, err := gridfs.NewBucket(e.DbConfig.MongoClient, nil)
		if err != nil {
			log.Println(err)
			// return err
		}

		fileBuffer := bytes.NewBuffer(nil)
		if _, err := bucket.DownloadToStream(event.Image, fileBuffer); err != nil {
			log.Println("There was an error downloading file: ", err.Error())
		}

		contentStr := base64.StdEncoding.EncodeToString(fileBuffer.Bytes())
		event.ImageContet = &contentStr
	}

	return
}
