//+build migratev5

package main

import (
	"database/sql"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	getBiosQuery = "SELECT `user_id`, `sys_member_id`, `system_id` FROM `user_bios`"
	updateSystemIDQuery = "UPDATE `user_bios` SET `system_id` = ? WHERE `user_id` = ? AND `sys_member_id` = ?"
)

var (
	pendingUpdates [][3]string
	missingSystems []string
)

func enqueueUpdate(userID string, memberID string, systemID string) {
	pendingUpdates = append(pendingUpdates, [3]string{userID, memberID, systemID})
}

func main() {
	dbFilename := os.Args[1]

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot open DB file %s", dbFilename)
	}
	defer db.Close()

	bios, err := db.Query(getBiosQuery)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get bios from DB")
	}

	fmt.Println("Checking bio entries")

	var n int
	for bios.Next() {
		n += 1

		var (
			userID string
			memberID string
			systemID string
		)

		_ = bios.Scan(&userID, &memberID, &systemID)

		// the IDs we want to update are short `zoxfc` style IDs
		if len(systemID) != 5 {
			continue
		}

		fmt.Printf("Updating %d: User:%#v SystemMember:%#v SystemID:%#v\n", n, userID, memberID, systemID)

		system, err := pluralkit.SystemById(systemID)
		if err != nil {

			if e, ok := err.(*pluralkit.Error); ok {
				if e.Code == pluralkit.ErrorCodeSystemNotFound {
					missingSystems = append(missingSystems, systemID)
                    continue
				}
			}

			log.Fatal().Err(err).Msgf("cannot get PK system ID %s", systemID)
		}

		systemID = system.UUID

		enqueueUpdate(userID, memberID, systemID)
	}

	defer bios.Close() // cursor needs to be closed before applying any updates, otherwise the DB remains locked

	fmt.Println("\nApplying DB updates")

	for _, update := range pendingUpdates {
		_, err := db.Exec(updateSystemIDQuery, update[2], update[0], update[1])
		if err != nil {
			log.Fatal().Err(err).Msgf("cannot update bio %s:%s", update[0], update[1])
		}
	}

	fmt.Printf("\nMissing systems: %#v\n", missingSystems)

	fmt.Println("\nmigrated db successfully to v5 format")

}
