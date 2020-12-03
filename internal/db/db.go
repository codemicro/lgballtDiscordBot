package db

import (
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var Conn *gorm.DB

func init() {

	fmt.Println("Setting up db")

	var err error
	Conn, err = gorm.Open(sqlite.Open(config.DbFileName), &gorm.Config{})
	Conn, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"lgballtUser", "K6Y9ljg6s4bV3LQY", "localhost:3306", "lgballtBot")), &gorm.Config{})
	if err != nil {
		logging.Error(err, fmt.Sprintf("Unable to open db file %s", config.DbFileName))
		os.Exit(1)
	}

	err = Conn.AutoMigrate(&UserBio{})
}

func RecordNotFound(g *gorm.DB) bool {
	return errors.Is(g.Error, gorm.ErrRecordNotFound)
}
