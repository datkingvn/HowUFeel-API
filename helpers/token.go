package helpers

import (
	"HowUFeel-API-Prj/configs"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

var jwtSecretKey []byte

func SetJWTSecretKey(secret string) {
	jwtSecretKey = []byte(secret)
}

func GetJWTSecretKey() []byte {
	return jwtSecretKey
}

func VerifyToken(tokenString string) (*Claims, error) {
	secretKey := GetJWTSecretKey()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid Token")
}

func GenerateToken(userId, email, role string) (string, string, error) {
	accessTokenExpiry := time.Now().Add(24 * time.Hour)
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)

	claims := Claims{
		UserID: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
		},
	}

	refreshClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiry),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func UpdateUserTokens(signedAccessToken, signedRefreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userCollection := configs.GetCollection("users")
	updateObj := bson.D{
		{"$set", bson.D{
			{"accessToken", signedAccessToken},
			{"refreshToken", signedRefreshToken},
			{"updatedAt", time.Now().UTC()},
		}},
	}

	filter := bson.M{"userId": userId}
	_, err := userCollection.UpdateOne(ctx, filter, updateObj)
	if err != nil {
		return err
	}
	return nil
}
