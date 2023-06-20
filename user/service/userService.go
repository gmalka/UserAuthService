package service

import (
	"time"
	"userService/internal/auth/tokenHandler"
	"userService/user/model"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	GetUser(email string) (model.User, error)
	GetAll() ([]model.User, error)
	Update(u model.User) error
	Create(u model.User) error
	Delete(username string) error
}

type TokenManager interface {
	CreateToken(username string, ttl time.Duration, kind int) (string, error)
	ParseToken(inputToken string, kind int) (UserClaims, error)
}

type UserService interface {
	GetUser(email string) (model.User, error)
	GetAll() ([]model.User, error)
	Update(u model.User) error
	Create(u model.User) error
	Delete(username string) error
	Register(user model.User) error
	AuthorizeEmail(user LoginRequest) (string, string, error)
	ParseToken(inputToken string, kind int) (UserClaims, error)
	GenerateTokens(email string) (string, string, error)
}

func NewAuthService(us UserStorage) UserService {
	return userService{us: us}
}

func (a userService) GetUser(email string) (model.User, error) {
	return  a.us.GetUser(email)
}

func (a userService) GetAll() ([]model.User, error) {
	return a.us.GetAll()
}

func (a userService) Update(u model.User) error {
	return a.us.Update(u)
}

func (a userService) Create(u model.User) error {
	return a.us.Create(u)
}

func (a userService) Delete(username string) error {
	return a.us.Delete(username)
}

func (a userService) Register(user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	err = a.us.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (a userService) AuthorizeEmail(user LoginRequest) (string, string, error) {
	wanteduser, err := a.us.GetUser(user.Email)
	if err != nil {
		return "", "", err
	}
 
	if err := bcrypt.CompareHashAndPassword([]byte(wanteduser.Password), []byte(user.Password)); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := a.GenerateTokens(user.Email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a userService) GenerateTokens(email string) (string, string, error) {
	accessToken, err := a.tokenHandler.CreateToken(email, 22, tokenHandler.AccessToken)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.tokenHandler.CreateToken(email, 222, tokenHandler.RefreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a userService) ParseToken(inputToken string, kind int) (UserClaims, error) {
	return a.tokenHandler.ParseToken(inputToken, kind)
}


// func (a userService) CreateToken(username string, ttl time.Duration, kind int) (string, error) {
// 	return a.tokenHandler.CreateToken(username, ttl, kind)
// }

// func (a userService) ParseToken(inputToken string, kind int) (UserClaims, error) {
// 	return a.tokenHandler.ParseToken(inputToken, kind)
// }

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type userService struct {
	us					UserStorage
	tokenHandler		TokenManager
}

type AuthorizeOut struct {
	UserID       int
	AccessToken  string
	RefreshToken string
}

type UserClaims struct {
	//	ID     		string `json:"id"`
	Username	string `json:"username"`
	Email		string `json:"email"`
	jwt.RegisteredClaims
}
