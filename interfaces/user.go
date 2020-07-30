package interfaces

import "github.com/spectrum-team/tachartas/models"

type UserRepository interface {
	UpsertUser(*models.User) error
	AssistToEvent(string, *models.Event, int) (bool, error)
}
