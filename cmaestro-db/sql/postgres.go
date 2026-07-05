package sql

import (
	"fmt"

	"gorm.io/driver/mysql"
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
	/*
		HARD CODED VALUES INSIDE DB CONNECTION
		Charset   string (utf-8)
		ParseTime bool   (True)
		Local     string (Local)
	*/
}

type Client struct {
	dbConn *gorm.DB
}

func New(cfg Config) (*Client, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/gorm?charset=utf8&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Endpoint),
		DefaultStringSize:         defaultStringSize,             // default size for string fields
		DisableDatetimePrecision:  !defaultDateTimePrecision,     // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    !defaultSupportRenameIndex,    // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   !defaultSupportRenameColumn,   // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: !defaultInitializeWithVersion, // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Client{
		dbConn: db,
	}, nil
}
