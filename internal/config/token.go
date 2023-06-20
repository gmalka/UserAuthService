package config

import "os"

type Token struct {
	AccessSecret string
	RefreshSecret string
}

func InitTokenConfig() Token {
	return Token{
		AccessSecret: os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
	}
}