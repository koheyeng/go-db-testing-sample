package db

import (
	"go-db-testing-sample/sample/db/models"
	"log"

	"github.com/jinzhu/gorm"
)

type DBHandler struct {
	db *gorm.DB
}

func NewDBHandler(db *gorm.DB) *DBHandler {
	return &DBHandler{db: db}
}

func (s *DBHandler) StoreSample(sample *models.Sample) error {

	if err := s.db.Save(sample).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *DBHandler) FetchSample() (*models.Sample, error) {
	sample := &models.Sample{}

	if err := s.db.Find(sample).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return sample, nil
}
