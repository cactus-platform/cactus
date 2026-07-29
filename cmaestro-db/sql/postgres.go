package sql

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultStringSize            = 256
	defaultDateTimePrecision     = true
	defaultSupportRenameIndex    = true
	defaultSupportRenameColumn   = true
	defaultInitializeWithVersion = true
)

type Config struct {
	Endpoint string
	Username string
	Password string
	Database string
	Port     int
	/*
		HARD CODED VALUES INSIDE DB CONNECTION
		Charset   string (utf-8)
		ParseTime bool   (True)
		Local     string (Local)
	*/
}

type Client struct {
	*gorm.DB
}

func New(cfg Config) (*Client, error) {
	credentials := url.UserPassword(cfg.Username, cfg.Password)
	password, _ := credentials.Password()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Europe/Paris",
		cfg.Endpoint, credentials.Username(), password, cfg.Database, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
