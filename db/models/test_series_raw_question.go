package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestSeriesRawQuestionStatuses holds the possible statuses for a raw question
var TestSeriesRawQuestionStatuses = struct {
	Pending    string
	Processing string
	Processed  string
	Failed     string
}{
	Pending:    "PENDING",
	Processing: "PROCESSING",
	Processed:  "PROCESSED",
	Failed:     "FAILED",
}

// TestSeriesRawQuestion represents the model for test series raw questions
type TestSeriesRawQuestion struct {
	*gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	RawQuestionData string    `gorm:"type:text;not null" json:"raw_question_data"`
	QuestionID      *int64    `gorm:"type:bigint;null" json:"question_id"`
	Remark          *string   `gorm:"type:text;null" json:"remark"`
	Meta            *string   `gorm:"type:jsonb;null" json:"meta"`
	Status          string    `gorm:"type:varchar(20);not null;default:PENDING" json:"status"`
}

// BeforeCreate sets a UUID before inserting a new record
func (q *TestSeriesRawQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	return nil
}
