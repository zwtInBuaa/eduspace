package database

import (
	"EDU_TH_2_backend/gin/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"

	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/models"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error

	switch config.GetString("db.type") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.GetString("db.mysql.user"),
			config.GetString("db.mysql.password"),
			config.GetString("db.mysql.host"),
			config.GetString("db.mysql.port"),
			config.GetString("db.mysql.database"),
		)
		logger.Info(dsn)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			config.GetString("db.postgresql.host"),
			config.GetString("db.postgresql.user"),
			config.GetString("db.postgresql.password"),
			config.GetString("db.postgresql.database"),
			config.GetString("db.postgresql.port"),
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		// 检查 SQLite 数据库文件是否存在
		if _, err := os.Stat(config.GetString("db.sqlite.file")); os.IsNotExist(err) {
			// 如果数据库文件不存在，则创建一个空的数据库文件
			file, err := os.Create(config.GetString("db.sqlite.file"))
			if err != nil {
				panic(err)
			}
			file.Close()
		}
		db, err = gorm.Open(sqlite.Open(config.GetString("db.sqlite.file")), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", config.GetString("db.type"))
	}

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Post{},
		&models.Comment{},
		&models.Question{},
		&models.Exam{},
		&models.UserSubmitHistory{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
