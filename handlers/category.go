package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spectrum-team/tachartas/interfaces"
	"github.com/spectrum-team/tachartas/models"
	"github.com/spectrum-team/tachartas/repos"
)

type CategoryHandler struct {
	dbConfig     *models.DatabaseConfig
	categoryRepo interfaces.CategoryRepository
}

func NewCategoryHandler(dbConfig *models.DatabaseConfig) *CategoryHandler {
	return &CategoryHandler{
		dbConfig:     dbConfig,
		categoryRepo: repos.NewCategoryRepository(dbConfig),
	}
}

func (c *CategoryHandler) FindAll(w http.ResponseWriter, r *http.Request) {

	res, err := c.categoryRepo.FindAll()
	if err != nil {
		log.Println("There was an error looking for categories: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
