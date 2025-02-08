package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// USER_STATUSES represents the possible user statuses
var USER_STATUSES = struct {
	ACTIVE   string
	INACTIVE string
}{
	ACTIVE:   "ACTIVE",
	INACTIVE: "INACTIVE",
}

// User represents the User model in the database.
type User struct {
	*gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PhoneNumber string    `gorm:"type:varchar(255);not null"`
	CountryCode string    `gorm:"type:varchar(10);not null"`
	Status      string    `gorm:"type:varchar(50);not null;default:ACTIVE"`
}

// BeforeCreate hook to generate UUID before inserting a new record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}
