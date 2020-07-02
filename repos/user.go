package repos

import (
	"errors"
	"log"

	"github.com/spectrum-team/tachartas/commons"
	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DbConfig   *models.DatabaseConfig
	Collection *mongo.Collection
}

func NewUserRepository(dbConfig *models.DatabaseConfig) *UserRepository {
	return &UserRepository{
		DbConfig:   dbConfig,
		Collection: dbConfig.MongoClient.Collection("user"),
	}
}

func (u *UserRepository) UpsertUser(user *models.User) error {

	if user.Email == "" || !commons.IsAnEmail(user.Email) {
		return errors.New("Invalid Email")
	}

	// Find a user
	ex := u.Collection.FindOne(u.DbConfig.Ctx, bson.M{"email": user.Email})
	if ex.Err() != nil {
		if ex.Err() == mongo.ErrNoDocuments {

			// If nothings found create a new user
			user.ID = primitive.NewObjectID()
			_, err := u.Collection.InsertOne(u.DbConfig.Ctx, user)
			if err != nil {
				log.Println("There was an error creating new user: ", err)
				return err
			}

			return nil
		}
		log.Println("There was an error looking for users: ", ex.Err())
		return ex.Err()
	}

	return nil
}

func (u *UserRepository) AssistToEvent(email string, event *models.Event, willAssist int) error {

	//Find the user
	res := u.Collection.FindOne(u.DbConfig.Ctx, bson.M{"email": email})
	if res.Err() != nil {
		log.Println("Error looking user: ", res.Err())
		return res.Err()
	}

	user := &models.User{}
	err := res.Decode(&user)
	if err != nil {
		log.Println("error decoding: ", err)
		return err
	}

	switch willAssist {
	case 1:
		e := &models.UpcomingEvents{
			EventID:    event.ID,
			EventName:  event.Name,
			EventDate:  event.Date,
			EventVenue: event.Venue,
		}

		user.UpcomingEvents = append(user.UpcomingEvents, e)
	case 0:
		indx := findIndex(user.UpcomingEvents, event.ID)
		if indx > -1 {
			user.UpcomingEvents = append(user.UpcomingEvents[:indx], user.UpcomingEvents[indx+1:]...)
		}
	}

	_, err = u.Collection.UpdateOne(
		u.DbConfig.Ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": user,
		},
	)
	if err != nil {
		log.Println("There was an error updating event: ", err)
		return err
	}

	return nil
}

func findIndex(list []*models.UpcomingEvents, eventID primitive.ObjectID) int {

	for i, item := range list {
		if item.EventID == eventID {
			return i
		}
	}

	return -1
}
