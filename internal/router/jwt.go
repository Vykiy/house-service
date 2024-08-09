package router

import (
	"github.com/Vykiy/house-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTIssuer struct {
	jwtSecret string
}

func NewJWTIssuer(jwtSecret string) *JWTIssuer {
	return &JWTIssuer{jwtSecret: jwtSecret}
}

func (j *JWTIssuer) IssueToken(userType models.UserType, userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_type": userType,
		"user_id":   userID.String(),
	})

	return token.SignedString([]byte(j.jwtSecret))
}

func (j *JWTIssuer) ParseToken(tokenString string) (models.UserType, uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.jwtSecret), nil
	})
	if err != nil {
		return models.UserTypeUnknown, uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.UserTypeUnknown, uuid.Nil, err
	}

	userType, ok := claims["user_type"].(string)
	if !ok {
		return models.UserTypeUnknown, uuid.Nil, err
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return models.UserTypeUnknown, uuid.Nil, err
	}

	return models.UserType(userType), userID, nil
}
