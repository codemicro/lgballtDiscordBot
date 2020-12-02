package db

import (
	"errors"
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
	Conn, err = gorm.Open(sqlite.Open(config.Config.DbFileName), &gorm.Config{})
	if err != nil {
		logging.Error(err, fmt.Sprintf("Unable to open db file %s", config.Config.DbFileName))
		os.Exit(1)
	}

	err = Conn.AutoMigrate(&UserBio{})
}

func RecordNotFound(g *gorm.DB) bool {
	return errors.Is(g.Error, gorm.ErrRecordNotFound)
}
