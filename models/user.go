package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastName" validate:"required,min=2,max=100"`
	Password     *string            `json:"password" binding:"required,min=6"`
	Email        *string            `json:"email" validate:"email, required"`
	PhoneNumber  *string            `json:"phoneNumber" validate:"required"`
	Token        *string            `json:"token,omitempty"`
	Role         *string            `json:"role" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string            `json:"refreshToken,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserID       string             `json:"userId"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.UserID,
		FirstName:   *u.FirstName,
		LastName:    *u.LastName,
		Email:       *u.Email,
		PhoneNumber: *u.PhoneNumber,
		Role:        *u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
