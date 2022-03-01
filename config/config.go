package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"irisblog/model"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type configData struct {
	DB     mysqlConfig  `json:"mysql"`
	Server serverConfig `json:"server"`
}

var (
	ExecPath    string
	JsonData    configData
	SeverConfig serverConfig
	DB          *gorm.DB
)

func initPath() {
	root := filepath.Dir(os.Args[0])
	ExecPath, _ = filepath.Abs(root)
}

func initJson() {
	jsonPath := ExecPath + "/config/config.json"
	rawConfig, err := os.ReadFile(jsonPath)
	if err != nil {
		rawConfig = []byte(`{"mysql":{"host":"localhost","port": 3306},"server":{"site_name": "博客","env":"development","port": 8080,"log_level": "debug"}}`)
	}
	if errJson := json.Unmarshal(rawConfig, &JsonData); errJson != nil {
		log.Println("Invalid config.json", errJson.Error())
		os.Exit(-1)
	}
}

func initServer() {
	SeverConfig = JsonData.Server
}

func InitDB(setting *mysqlConfig) error {
	var (
		db  *gorm.DB
		err error
	)
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", setting.User, setting.Password, setting.Host, setting.Port, setting.Database)
	setting.Url = dns
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		if !strings.Contains(err.Error(), "1049") {
			return err
		}
		tempDns := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", setting.User, setting.Password, setting.Host, setting.Port)
		db, err = gorm.Open(mysql.Open(tempDns), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		if err != nil {
			return err
		}
		err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s charset=utf8", setting.Database)).Error
		if err != nil {
			return err
		}
		db, err = gorm.Open(mysql.Open(dns), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		if err != nil {
			return err
		}

	}
	sqlDB, errDB := db.DB()
	if errDB != nil {
		log.Fatalln(errDB.Error())
		os.Exit(-1)
	}
	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetConnMaxLifetime(-1)
	db.AutoMigrate(&model.Admin{}, &model.Article{}, &model.ArticleData{}, &model.Category{}, &model.Attachment{})
	DB = db
	return nil
}

func init() {
	initPath()
	initJson()
	initServer()
	if JsonData.DB.Database != "" {
		err := InitDB(&JsonData.DB)
		if err != nil {
			log.Fatalln("Failed To Connect Database: ", err.Error())
			os.Exit(-1)
		}
	}
}

func WriteConfig() error {
	jsonPath := ExecPath + "/config/config.json"
	configFile, err := os.OpenFile(jsonPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer configFile.Close()
	buff := &bytes.Buffer{}
	buf, errJson := json.MarshalIndent(JsonData, "", "\t")
	if errJson != nil {
		return errJson
	}
	buff.Write(buf)
	_, errFile := io.Copy(configFile, buff)
	if errFile != nil {
		return errFile
	}
	return nil
}
