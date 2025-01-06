package db

import (
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"ndm/internal/conf"
	"ndm/internal/model"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	return db
}

func InitDb() {
	logLevel := logger.Silent
	newLogger := logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,
			Colorful:                  true, // 禁用彩色打印
		},
	)

	database := conf.Database

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: database.TablePrefix,
		},
		Logger: newLogger,
	}

	var dB *gorm.DB
	var err error

	switch database.Type {
	case "sqlite3":
		{
			if !(strings.HasSuffix(database.Path, ".db") && len(database.Path) > 3) {
				log.Fatalf("db name error.")
			}
			dB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
				database.Path)), gormConfig)
		}
	case "mysql":
		{
			dsn := database.DSN
			if dsn == "" {
				//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
					database.User, database.Password, database.Host, database.Port, database.Name, database.SSLMode)
			}
			dB, err = gorm.Open(mysql.Open(dsn), gormConfig)
		}
	case "postgres":
		{
			dsn := database.DSN
			if dsn == "" {
				dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
					database.Host, database.User, database.Password, database.Name, database.Port, database.SSLMode)
			}
			dB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		}
	default:
		log.Fatalf("not supported database type: %s", database.Type)
	}

	if err != nil {
		log.Fatalf("failed to connect database:%s", err.Error())
	}

	Init(dB)
}

func Init(d *gorm.DB) {
	db = d
	err := AutoMigrate(new(model.User))
	if err != nil {
		log.Fatalf("failed migrate database: %s", err.Error())
	}
}

func AutoMigrate(dst ...interface{}) error {
	var err error
	if conf.Database.Type == "mysql" {
		err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(dst...)
	} else {
		err = db.AutoMigrate(dst...)
	}
	return err
}
