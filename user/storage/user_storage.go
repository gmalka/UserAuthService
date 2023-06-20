package storage

import (
	"userService/user/model"

	"github.com/jmoiron/sqlx"
)

type UserStorage interface {
	GetUser(email string) (model.User, error)
	GetAll() ([]model.User, error)
	Update(u model.User) error
	Create(u model.User) error
	Delete(username string) error
}

func NewUserStorage(db *sqlx.DB) UserStorage {
	return &userStorage{db: db}
}

type userStorage struct {
	db *sqlx.DB
}

func (u userStorage) GetUser(email string) (model.User, error) {
	var user model.User
	err := u.db.QueryRow("SELECT * FROM users WHERE email=$1", email).Scan(&u)

	return user, err
}

func (u userStorage) GetAll() ([]model.User, error) {
	var users []model.User
	rows, err := u.db.Query("SELECT * FROM users")
	if err != nil{
		return nil, err
	}

	err = rows.Scan(&u)
	if err != nil{
		return nil, err
	}

	return users, nil
}

func (u userStorage) Update(user model.User) error {
	_, err := u.db.Exec("UPDATE users SET username=$1,password=$2,email=$3 WHERE id=$4", user.Username, user.Password, user.Email, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (u userStorage) Create(user model.User) error {
	_, err := u.db.Exec("INSERT INTO users(username,password,email) VALUES($1,$2,$3) ", user.Username, user.Password, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (u userStorage) Delete(username string) error {
	_, err := u.db.Exec("DELETE FROM users WHERE username=$1", username)

	return err
}