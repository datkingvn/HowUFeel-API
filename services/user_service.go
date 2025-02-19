package services

import (
	"HowUFeel-API-Prj/configs"
	"HowUFeel-API-Prj/helpers"
	"HowUFeel-API-Prj/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var userCollection = configs.GetCollection("users")

func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func RegisterUser(user *models.User) (*models.UserResponse, error) {
	ctx, cancel := createContext(10 * time.Second)
	defer cancel()

	filter := bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"phoneNumber": user.PhoneNumber},
	}}
	existingCount, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existingCount > 0 {
		return nil, errors.New("user already exists")
	}

	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	user.Password = helpers.HashAndSalt(user.Password)
	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
	defaultRole := "USER"
	user.Role = &defaultRole

	accessToken, refreshToken, err := helpers.GenerateToken(user.UserID, *user.Email, *user.Role)
	if err != nil {
		return nil, errors.New("token generation failed")
	}
	user.Token, user.RefreshToken = &accessToken, &refreshToken

	if _, err := userCollection.InsertOne(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user.ToResponse(), nil
}

func LoginUser(loginRequest *models.LoginRequest) (*models.UserResponse, string, string, error) {
	ctx, cancel := createContext(10 * time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
	if err != nil || user.Password == nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	passwordIsValid, err := helpers.VerifyPassword(*user.Password, loginRequest.Password)
	if err != nil || !passwordIsValid {
		return nil, "", "", errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := helpers.GenerateToken(user.UserID, *user.Email, *user.Role)
	if err != nil {
		return nil, "", "", errors.New("failed to generate token")
	}

	update := bson.M{"$set": bson.M{"token": accessToken, "refreshToken": refreshToken}}
	if _, err := userCollection.UpdateOne(ctx, bson.M{"userID": user.UserID}, update); err != nil {
		return nil, "", "", fmt.Errorf("failed to update tokens: %w", err)
	}

	return user.ToResponse(), accessToken, refreshToken, nil
}

func GetUserByID(userID, requesterID, requesterRole string) (*models.UserResponse, error) {
	if requesterID != userID && requesterRole != "ADMIN" {
		return nil, errors.New("unauthorized: access denied")
	}

	ctx, cancel := createContext(10 * time.Second)
	defer cancel()

	var user models.User
	objID, err := primitive.ObjectIDFromHex(userID)
	if err == nil {
		err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	} else {
		err = userCollection.FindOne(ctx, bson.M{"userID": userID}).Decode(&user)
	}

	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

func GetAllUsers(requesterRole string) ([]models.UserResponse, error) {
	if requesterRole != "ADMIN" {
		return nil, errors.New("forbidden: only admins can fetch users")
	}

	ctx, cancel := createContext(10 * time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToResponse()
	}
	return userResponses, nil
}
