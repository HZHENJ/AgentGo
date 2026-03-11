package common

import (
	"fmt"
	"time"
	"agentgo/pkg/conf"
	"agentgo/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() error {
	c := conf.Config.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.DbName, c.Charset, c.ParseTime, c.Loc,
	)

	var dbLogger logger.Interface
	if conf.Config.Service.AppMode == "debug" {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, 
		},
	})

	if err != nil {
		return err 
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20) 
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	DB = db

	return migration()
}

func migration() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Message{},
	)
}