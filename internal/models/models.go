package models

import "github.com/google/uuid"

type UserType string

const (
	UserTypeUnknown   UserType = "unknown"
	UserTypeUser      UserType = "user"
	UserTypeModerator UserType = "moderator"
)

type FlatStatus string

const (
	FlatStatusUnknown      FlatStatus = "unknown"
	FlatStatusCreated      FlatStatus = "created"
	FlatStatusApproved     FlatStatus = "approved"
	FlatStatusDeclined     FlatStatus = "declined"
	FlatStatusOnModeration FlatStatus = "on_moderation"
)

type House struct {
	ID        int    `json:"id" db:"id"`
	Address   string `json:"address" db:"address"`
	YearBuilt int    `json:"yearBuilt" db:"year_built"`
	Developer string `json:"developer" db:"developer"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type Flat struct {
	ID      int        `json:"id" db:"flat_number"` // используем номер квартиры относительно дома в качестве ID
	HouseID int        `json:"houseId" db:"house_id"`
	Price   int        `json:"price" db:"price"`
	Rooms   int        `json:"rooms" db:"rooms"`
	Status  FlatStatus `json:"status" db:"status"`
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	UserType     UserType  `json:"userType" db:"user_type"`
}
