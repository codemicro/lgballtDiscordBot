from os import system
import sqlite3
import requests
import sys

db_location = sys.argv[1]

createNewBiosTable = """CREATE TABLE "user_bios_temp" (
	"user_id"	text,
	"sys_member_id"	text,
	"raw_bio_data"	text,
	"system_id"	TEXT,
	PRIMARY KEY("user_id","sys_member_id")
)"""
insertSql = "INSERT INTO `user_bios_temp` (`user_id`, `sys_member_id`, `raw_bio_data`, `system_id`) VALUES(?, ?, ?, ?)"
selectAllFromOld = "SELECT * FROM `user_bios`"
dropOld = "DROP TABLE `user_bios`"
renameNew = "ALTER TABLE `user_bios_temp` RENAME TO `user_bios`"

id_cache = {}

with sqlite3.connect(db_location) as db:	
    cursor = db.cursor()
    
    # create new table
    cursor.execute(createNewBiosTable)
    db.commit()

    to_do = []

    for item in cursor.execute(selectAllFromOld):
        uid, sysmemid, biotext = item
        to_do.append((uid, sysmemid, biotext))

    # insert existing removals
    for args in to_do:

        system_id = ""
        if args[1] != "":  # system member ID

            csid = id_cache.get(args[0], None) 
            if csid is not None:
                system_id = csid
            else:
                print(f"Requesting {args[0]}")
                r = requests.get(f"https://api.pluralkit.me/v1/a/{args[0]}")
                if r.status_code != 200: 
                    print(f"Warning: unable to fetch system information for user {args[0]}")
                else: 
                    x = r.json()["id"]
                    system_id = x
                    id_cache[args[0]] = x

        cursor.execute(insertSql, (args[0], args[1], args[2], system_id))
    db.commit()

    # delete old table
    cursor.execute(dropOld)
    # rename temp
    cursor.execute(renameNew)

    db.commit()

print("migrated db successfully to v4 format")
