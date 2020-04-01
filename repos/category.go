package repos

import (
	"log"

	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository struct {
	DbConfig   *models.DatabaseConfig
	Collection *mongo.Collection
}

func NewCategoryRepository(dbConfig *models.DatabaseConfig) *CategoryRepository {
	return &CategoryRepository{
		DbConfig:   dbConfig,
		Collection: dbConfig.MongoClient.Collection("category"),
	}
}

func (c *CategoryRepository) FindAll() ([]*models.Category, error) {

	cursor, err := c.Collection.Find(c.DbConfig.Ctx, primitive.M{})
	if err != nil {
		log.Println("There was an error getting all categories: ", err)
		return nil, err
	}

	categories := make([]*models.Category, 0)
	err = cursor.All(c.DbConfig.Ctx, &categories)
	if err != nil {
		log.Println("There was an error decoding categories: ", err)
		return nil, err
	}

	return categories, nil
}
