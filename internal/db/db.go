package db

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var Conn *gorm.DB

func init() {

	fmt.Println("Setting up db")

	var err error
	Conn, err = gorm.Open(sqlite.Open(config.DbFileName), &gorm.Config{})
	if err != nil {
		logging.Error(err, fmt.Sprintf("Unable to open db file %s", config.DbFileName))
		os.Exit(1)
	}

	err = Conn.AutoMigrate(&UserBio{})
}
