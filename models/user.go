package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	DisplayName     string             `json:"displayName"`
	Email           string             `json:"email"`
	IsEmailVerified bool               `json:"isEmailVerified"`
	PhoneNumber     *string            `json:"phoneNumber"`
	PhotoURL        *string            `json:"photoUrl"`
	ProviderID      string             `json:"providerId"`
	UserID          string             `json:"userId"`
	DeviceID        *string            `json:"deviceId"`
	UpcomingEvents  []*UpcomingEvents  `json:"upcoming_events"`
}

type UpcomingEvents struct {
	EventID    primitive.ObjectID `json:"event_id"`
	EventName  string             `json:"event_name"`
	EventVenue string             `bson:"event_venue" json:"event_venue,omitempty"`
	EventDate  time.Time          `bson:"event_date" json:"event_date,omitempty"`
}
