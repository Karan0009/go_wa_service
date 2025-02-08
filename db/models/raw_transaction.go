package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var RawTransactionStatuses = struct {
	PendingTextExtraction string
	ExtractingText        string
	Pending               string
	Processing            string
	Processed             string
	Failed                string
	Invalid               string
}{
	PendingTextExtraction: "PENDING_TEXT_EXTRACTION",
	ExtractingText:        "EXTRACTING_TEXT",
	Pending:               "PENDING",
	Processing:            "PROCESSING",
	Processed:             "PROCESSED",
	Failed:                "FAILED",
	Invalid:               "INVALID",
}

var RawTransactionType = struct {
	WAImage string
	WAText  string
}{
	WAImage: "WA_IMAGE",
	WAText:  "WA_TEXT",
}

type RawTransaction struct {
	*gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID             uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RawTransactionType string    `gorm:"type:varchar(50);not null" json:"raw_transaction_type"`
	RawTransactionData string    `gorm:"type:text;not null" json:"raw_transaction_data"`
	ExtractedText      *string   `gorm:"type:text;null" json:"extracted_text"`
	TransactionID      *int64    `gorm:"type:bigint;null" json:"transaction_id"`
	Remark             *string   `gorm:"type:text;null" json:"remark"`
	Meta               *string   `gorm:"type:jsonb;null" json:"meta"`
	Status             string    `gorm:"type:varchar(50);not null;default:PENDING" json:"status"`
}

// BeforeCreate sets a UUID before inserting a new record
func (rt *RawTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	rt.ID = uuid.New()
	return nil
}
