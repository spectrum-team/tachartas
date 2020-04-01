package interfaces

import "github.com/spectrum-team/tachartas/models"

type CategoryRepository interface {
	FindAll() ([]*models.Category, error)
}
