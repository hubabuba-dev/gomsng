package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	RefreshTokenSecret string
	JWTTokenSecret     string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	PGConnDSN          string
}

func InitConfiguration() *Config {
	var err error
	var access_ttl_time time.Duration
	var refresh_ttl_time time.Duration
	access_ttl := os.Getenv("ACCESS_TOKEN_TTL")
	if access_ttl == "" {
		log.Fatal("No ACCESS token ttl")
	} else {
		access_ttl_time, err = time.ParseDuration(access_ttl)
		if err != nil {
			log.Fatal("Wrong format of access token ttl")
		}
	}
	refresh_ttl := os.Getenv("REFRESH_TOKEN_TTL")
	if refresh_ttl == "" {
		log.Fatal("No Refresh token ttl")
	} else {
		refresh_ttl_time, err = time.ParseDuration(refresh_ttl)
		if err != nil {
			log.Fatal("Wrong format of access token ttl")
		}
	}
	refresh_token_secret := os.Getenv("TOKEN_SECRET")
	if refresh_token_secret == "" {
		log.Fatal("No Refresh token secret")
	}
	jwt_token_secret := os.Getenv("JWT_TOKEN_SECRET")
	if jwt_token_secret == "" {
		log.Fatal("No jwt token secret")
	}
	pg_conn_dsn := os.Getenv("PG_CONN_DSN")
	if pg_conn_dsn == "" {
		log.Fatal("No pg conn dsn")
	}
	return &Config{
		RefreshTokenSecret: refresh_token_secret,
		JWTTokenSecret:     jwt_token_secret,
		AccessTokenTTL:     access_ttl_time,
		RefreshTokenTTL:    refresh_ttl_time,
		PGConnDSN:          pg_conn_dsn,
	}
}
