package db

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var Conn *gorm.DB

func init() {

	fmt.Println("Setting up db")

	dbConfig := new(gorm.Config)

	if config.DebugMode {
		dbConfig.Logger = logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: true,
			},
		)
	}

	var err error
	Conn, err = gorm.Open(sqlite.Open(config.DbFileName), dbConfig)
	if err != nil {
		logging.Error(err, fmt.Sprintf("Unable to open db file %s", config.DbFileName))
		os.Exit(1)
	}

	err = Conn.AutoMigrate(&UserBio{}, &ReactionRole{}, &userBan{}, &userKick{}, &VerificationFail{}, &UserMute{}, &AnalyticsEvent{}, &ToneTag{})
	if err != nil {
		logging.Error(err, "Failed to run database migrations")
		os.Exit(1)
	}
}
