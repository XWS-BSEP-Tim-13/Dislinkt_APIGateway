package config

import "os"

type Config struct {
	Port        string
	UserHost    string
	UserPort    string
	PostHost    string
	PostPort    string
	CompanyHost string
	CompanyPort string
	AuthHost    string
	AuthPort    string
}

func NewConfig() *Config {
	return &Config{
		Port:        os.Getenv("GATEWAY_PORT"),
		UserHost:    os.Getenv("USER_SERVICE_HOST"),
		UserPort:    os.Getenv("USER_SERVICE_PORT"),
		PostHost:    os.Getenv("POST_SERVICE_HOST"),
		PostPort:    os.Getenv("POST_SERVICE_PORT"),
		CompanyHost: os.Getenv("COMPANY_SERVICE_HOST"),
		CompanyPort: os.Getenv("COMPANY_SERVICE_PORT"),
		AuthHost:    os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort:    os.Getenv("AUTH_SERVICE_PORT"),
	}
}
