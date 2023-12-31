package dao

import (
	"downloader/pkg/model"
	"downloader/pkg/myLog"
	"downloader/pkg/ttviper"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func init() {
	config := ttviper.ReadConfig("./etc", "conf.yaml")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Viper.GetString("Mysql.User"),
		config.Viper.GetString("Mysql.Pass"),
		config.Viper.GetString("Mysql.Ip"),
		config.Viper.GetString("Mysql.Port"),
		config.Viper.GetString("Mysql.Database"),
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		myLog.Log.Error("err: init mysql failed...", err)
		return
	}
	db.AutoMigrate(&model.User{})

	DB = db
	myLog.Log.Info("mysql init successfully!")

}
