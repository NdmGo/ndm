package db

import (
	"fmt"
	stdlog "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	_ "modernc.org/sqlite"

	"ndm/internal/conf"
	"ndm/internal/model"
	"ndm/internal/utils"
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
			SlowThreshold:             time.Second, // slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,
			Colorful:                  true, // disable colorful printing
		},
	)

	database := conf.Database

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: database.TablePrefix,
		},
		Logger: newLogger,
		// performance optimization: enable prepared statement cache
		PrepareStmt: true,
		// performance optimization: disable auto ping
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var dB *gorm.DB
	var err error

	switch database.Type {
	case "sqlite3":
		{
			if !(strings.HasSuffix(database.Path, ".db") && len(database.Path) > 3) {
				log.Fatalf("db name error.")
			}

			if strings.HasPrefix(database.Path, "/") {
				dB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental", database.Path)), gormConfig)
			} else {
				conf_path := conf.WorkDir()
				custom_dir := fmt.Sprintf("%s/custom", conf_path)

				db_file := fmt.Sprintf("%s/%s", custom_dir, database.Path)
				db_dir := filepath.Dir(db_file)

				if !utils.IsExist(db_dir) {
					os.MkdirAll(db_dir, os.ModePerm)
				}

				dB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental", db_file)), gormConfig)
			}

		}
	case "mysql":
		{
			dsn := database.DSN
			if dsn == "" {
				//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
					database.User, database.Password, database.Hostname, database.Hostport, database.Name, database.SSLMode)
			}
			dB, err = gorm.Open(mysql.Open(dsn), gormConfig)
		}
	case "postgres":
		{
			dsn := database.DSN
			if dsn == "" {
				dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
					database.Hostname, database.User, database.Password, database.Name, database.Hostport, database.SSLMode)
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

	// performance optimization: configure connection pool
	sqlDB, err := db.DB()
	if err == nil {
		// set maximum number of idle connections in the connection pool
		sqlDB.SetMaxIdleConns(10)
		// set maximum number of open database connections
		sqlDB.SetMaxOpenConns(100)
		// set maximum time a connection can be reused
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	err = AutoMigrate(new(model.User), new(model.Storage), new(model.SettingItem), new(model.Logs), new(model.Tasks))
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

func CheckDbConnnect(data map[string]string) error {
	if strings.EqualFold(data["type"], "mysql") {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", data["username"], data["password"], data["hostname"], data["hostport"], data["dbname"])
		_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
	}
	return nil
}
