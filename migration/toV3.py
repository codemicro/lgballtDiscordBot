import sqlite3
import sys

db_location = sys.argv[1]

createNewBiosTable = "CREATE TABLE `user_bios_temp` (`user_id` text,`sys_member_id` text,`raw_bio_data` text,PRIMARY KEY (`user_id`,`sys_member_id`))"
insertSql = "INSERT INTO `user_bios_temp` (`user_id`, `sys_member_id`, `raw_bio_data`) VALUES(?, ?, ?)"
selectAllFromOld = "SELECT * FROM `user_bios`"
dropOld = "DROP TABLE `user_bios`"
renameNew = "ALTER TABLE `user_bios_temp` RENAME TO `user_bios`"

with sqlite3.connect(db_location) as db:	
    cursor = db.cursor()
    
    # create new table
    cursor.execute(createNewBiosTable)
    db.commit()

    to_do = []

    # recategorise existing removals
    for item in cursor.execute(selectAllFromOld):
        uid, biotext = item
        to_do.append((uid, biotext))

    # insert existing removals
    for args in to_do:
        cursor.execute(insertSql, (args[0], "", args[1]))
    db.commit()

    # delete old table
    cursor.execute(dropOld)
    # rename temp
    cursor.execute(renameNew)

    db.commit()

print("migrated db successfully to v3 format")
