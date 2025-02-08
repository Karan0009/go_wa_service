package test_series_raw_question_service

import (
	"errors"
	"wa_bot_service/db/models"

	db_service "wa_bot_service/modules/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestSeriesRawQuestionService struct {
	db *gorm.DB
}

// NewTestSeriesRawQuestionService creates a new instance of the service
func NewTestSeriesRawQuestionService() *TestSeriesRawQuestionService {
	return &TestSeriesRawQuestionService{db: db_service.GetDBClient()}
}

// CreateTestSeriesRawQuestion inserts a new raw question into the database
func (s *TestSeriesRawQuestionService) CreateTestSeriesRawQuestion(rawQuestionData string) (*models.TestSeriesRawQuestion, error) {
	rawQuestion := &models.TestSeriesRawQuestion{
		RawQuestionData: rawQuestionData,
		Status:          models.TestSeriesRawQuestionStatuses.Pending,
	}

	if err := s.db.Omit("deleted_at").Create(rawQuestion).Error; err != nil {
		return nil, err
	}

	return rawQuestion, nil
}

// GetTestSeriesRawQuestionByID fetches a raw question by its ID
func (s *TestSeriesRawQuestionService) GetTestSeriesRawQuestionByID(id uuid.UUID) (*models.TestSeriesRawQuestion, error) {
	var rawQuestion models.TestSeriesRawQuestion
	if err := s.db.First(&rawQuestion, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &rawQuestion, nil
}

// UpdateQuestionStatus updates the status of a raw question
func (s *TestSeriesRawQuestionService) UpdateQuestionStatus(id uuid.UUID, newStatus string) error {
	if err := s.db.Model(&models.TestSeriesRawQuestion{}).
		Where("id = ?", id).
		Update("status", newStatus).Error; err != nil {
		return err
	}
	return nil
}

// AssignQuestionID updates the question_id field for a raw question
func (s *TestSeriesRawQuestionService) AssignQuestionID(id uuid.UUID, questionID int64) error {
	if err := s.db.Model(&models.TestSeriesRawQuestion{}).
		Where("id = ?", id).
		Update("question_id", questionID).Error; err != nil {
		return err
	}
	return nil
}

// GetPendingRawQuestions fetches all raw questions that are still pending
func (s *TestSeriesRawQuestionService) GetPendingRawQuestions() ([]models.TestSeriesRawQuestion, error) {
	var questions []models.TestSeriesRawQuestion
	err := s.db.Where("status = ?", models.TestSeriesRawQuestionStatuses.Pending).Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// DeleteTestSeriesRawQuestion removes a raw question from the database
func (s *TestSeriesRawQuestionService) DeleteTestSeriesRawQuestion(id uuid.UUID) error {
	if err := s.db.Delete(&models.TestSeriesRawQuestion{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
