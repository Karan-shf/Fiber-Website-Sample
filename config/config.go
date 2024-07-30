package config

import "os"

var JWT_SECRET_KEY string

func LoadSecretKey() {
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
}
