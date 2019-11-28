package interfaces

import (
	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventRepository interface {
	FindByID(primitive.ObjectID) (*models.Event, error)
	Find(bson.M) ([]*models.Event, error)
	Insert(*models.Event) (primitive.ObjectID, error)
	Update(primitive.ObjectID, *models.Event) error
}
