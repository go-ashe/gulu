package postgres

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(conf Config) (*gorm.DB, error) {
	var opts []gorm.Option
	if conf.IsDebug {
		opts = append(opts, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	} else {
		opts = append(opts, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // 设置日志模式为 Silent
		})
	}

	db, err := gorm.Open(postgres.Open(conf.DSN), opts...)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if conf.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpen)
	}
	if conf.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdle)
	}
	if conf.MaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(int64(conf.MaxLifeTime)))
	}
	if conf.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(int64(conf.MaxIdleTime)))
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
