package config

import "os"

type Config struct {
	Port           string
	UserHost       string
	UserPort       string
	ConnectionHost string
	ConnectionPort string
	PostHost       string
	PostPort       string
	CompanyHost    string
	CompanyPort    string
	AuthHost       string
	AuthPort       string
	RBACConfig     string
	RBACPolicy     string
}

func NewConfig() *Config {
	return &Config{
		Port:           os.Getenv("GATEWAY_PORT"),
		UserHost:       os.Getenv("USER_SERVICE_HOST"),
		UserPort:       os.Getenv("USER_SERVICE_PORT"),
		ConnectionHost: os.Getenv("USER_SERVICE_HOST"),
		ConnectionPort: os.Getenv("USER_SERVICE_PORT"),
		PostHost:       os.Getenv("POST_SERVICE_HOST"),
		PostPort:       os.Getenv("POST_SERVICE_PORT"),
		CompanyHost:    os.Getenv("COMPANY_SERVICE_HOST"),
		CompanyPort:    os.Getenv("COMPANY_SERVICE_PORT"),
		AuthHost:       os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort:       os.Getenv("AUTH_SERVICE_PORT"),
		RBACConfig:     os.Getenv("RBAC_CONFIG"),
		RBACPolicy:     os.Getenv("RBAC_POLICY"),
	}
}
