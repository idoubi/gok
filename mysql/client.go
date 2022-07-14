package mysql

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 数据库配置
type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	Charset      string `json:"charset"`
	Debug        bool   `json:"debug"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxLifetime  int    `json:"maxLifetime"`
}

var dbch = make(chan map[string]*gorm.DB)

// InitWithName 初始化数据库连接
func InitWithName(name string) error {
	var conf Config
	sub := viper.Sub("mysql." + name)
	if sub == nil {
		return fmt.Errorf("invalid mysql config under %s", name)
	}
	if err := sub.Unmarshal(&conf); err != nil {
		return err
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}
	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database, conf.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if conf.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	}
	if conf.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	}
	if conf.Debug {
		db = db.Debug()
	}

	// set db
	dbch <- map[string]*gorm.DB{name: db}

	return nil
}

// GetClient 获取数据库连接
func GetClient(name string) *gorm.DB {
	dbm := <-dbch
	if db, ok := dbm[name]; ok {
		return db
	}

	return nil
}

// GetDB 获取数据库连接
func GetDB(name string) *gorm.DB {
	return GetClient(name)
}

func dbPool() {
	var dbs = make(map[string]*gorm.DB)
	for {
		select {
		case dbm := <-dbch:
			for name, db := range dbm {
				dbs[name] = db
			}
		case dbch <- dbs:
		}
	}
}

func init() {
	go dbPool()
}
