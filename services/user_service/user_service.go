package user_service

import (
	"fmt"
	"wa_bot_service/db/models"
	db_service "wa_bot_service/modules/db"

	"github.com/google/uuid"
)

// CreateUser function that uses the global dbClient directly from user_service
func CreateUser(phoneNumber, countryCode string) (*models.User, error) {
	// Proceed with user creation logic using the global dbClient
	dbClient := db_service.GetDBClient()
	user := &models.User{PhoneNumber: phoneNumber, CountryCode: countryCode}
	result := dbClient.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %v", result.Error)
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	dbClient := db_service.GetDBClient()
	result := dbClient.First(&user, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user by ID: %v", result.Error)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	dbClient := db_service.GetDBClient()
	result := dbClient.Unscoped().Where("phone_number = ?", phoneNumber).
		Where("country_code = ?", "+91").
		First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user by phoneNumber: %v", result.Error)
	}
	return &user, nil
}

// UpdateUser updates a user's details
func UpdateUser(id uuid.UUID, PhoneNumber string) (*models.User, error) {
	// Find the user by ID
	user := &models.User{ID: id}
	dbClient := db_service.GetDBClient()
	result := dbClient.First(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %v", result.Error)
	}

	// Update the user's fields
	user.PhoneNumber = PhoneNumber
	updateResult := dbClient.Save(user)
	if updateResult.Error != nil {
		return nil, fmt.Errorf("failed to update user: %v", updateResult.Error)
	}

	return user, nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(id uuid.UUID) error {
	// Find the user by ID
	user := &models.User{ID: id}
	dbClient := db_service.GetDBClient()
	result := dbClient.First(user)
	if result.Error != nil {
		return fmt.Errorf("failed to find user: %v", result.Error)
	}

	// Delete the user
	deleteResult := dbClient.Delete(user)
	if deleteResult.Error != nil {
		return fmt.Errorf("failed to delete user: %v", deleteResult.Error)
	}

	return nil
}
