package main

import (
	"os"
	"strings"

	"github.com/bcc-code/brunstadtv/backend/auth0"
	"github.com/samber/lo"
)

type postgres struct {
	ConnectionString string
}

type algolia struct {
	AppId            string
	ApiKey           string
	SearchOnlyApiKey string
}

type envConfig struct {
	DB        postgres
	Algolia   algolia
	Port      string
	JWTConfig auth0.JWTConfig
}

func getEnvConfig() envConfig {
	aud := lo.Map(strings.Split(os.Getenv("JWT_AUDIENCES"), ","),
		func(s string, _ int) string {
			return strings.TrimSpace(s)
		},
	)

	return envConfig{
		DB: postgres{
			ConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		},
		JWTConfig: auth0.JWTConfig{
			Domain:    os.Getenv("JWT_DOMAIN"),
			Issuer:    os.Getenv("JWT_ISSUER"),
			Audiences: aud,
		},
		Port: os.Getenv("PORT"),
		Algolia: algolia{
			AppId:            os.Getenv("ALGOLIA_APP_ID"),
			ApiKey:           os.Getenv("ALGOLIA_API_KEY"),
			SearchOnlyApiKey: os.Getenv("ALGOLIA_SEARCH_ONLY_API_KEY"),
		},
	}
}