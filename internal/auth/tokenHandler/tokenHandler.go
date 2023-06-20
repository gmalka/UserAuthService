package tokenHandler

import (
	"errors"
	"fmt"
	"time"
	"userService/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessToken = iota
	RefreshToken
)

type UserClaims struct {
	//	ID     		string `json:"id"`
	Username	string `json:"username"`
	Email		string `json:"email"`
	jwt.RegisteredClaims
}

type TokenManager interface {
	CreateToken(username, email string, ttl time.Duration, kind int) (string, error)
	ParseToken(inputToken string, kind int) (UserClaims, error)
}

type TokenJWT struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

func NewTokenJWT(token config.Token) TokenManager {
	return &TokenJWT{AccessSecret: []byte(token.AccessSecret), RefreshSecret: []byte(token.RefreshSecret)}
}

func (o *TokenJWT) CreateToken(username, email string, ttl time.Duration, kind int) (string, error) {
	claims := UserClaims{
		Username:			username,
		Email:				email,
		RegisteredClaims:	jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var secret []byte
	switch kind {
	case AccessToken:
		secret = o.AccessSecret
	case RefreshToken:
		secret = o.RefreshSecret
	default:
		return "", errors.New("unknown type of token")
	}

	return token.SignedString(secret)
}

func (o *TokenJWT) ParseToken(inputToken string, kind int) (UserClaims, error) {
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		var secret []byte
		switch kind {
		case AccessToken:
			secret = o.AccessSecret
		case RefreshToken:
			secret = o.RefreshSecret
		default:
			return "", errors.New("unknown type of token")
		}
		_ = secret

		return secret, nil
	})

	if err != nil {
		return UserClaims{}, err
	}

	if !token.Valid {
		return UserClaims{}, fmt.Errorf("not valid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserClaims{}, fmt.Errorf("error get user claims from token")
	}

	return UserClaims{
		Username:           claims["username"].(string),
		RegisteredClaims: 	jwt.RegisteredClaims{},
	}, nil
}