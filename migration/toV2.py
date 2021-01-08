import sqlite3
import sys

db_location = sys.argv[1]

userRemoveTables = "CREATE TABLE IF NOT EXISTS `{}` (`user_id` text, `reason` text, PRIMARY KEY (`user_id`))"
insertSql = "INSERT INTO `user_{}` (`user_id`, `reason`) VALUES(?, ?)"

with sqlite3.connect(db_location) as db:	
    cursor = db.cursor()
    
    # create new tables
    for tn in ["user_kicks", "user_bans"]:
        cursor.execute(userRemoveTables.format(tn))
    db.commit()

    to_do = []

    # recategorise existing removals
    for item in cursor.execute("SELECT * FROM `user_removes`"):
        uid, reason, action = item
        fib = ""
        if action.lower() == "banned":
            fib = "bans"
        elif action.lower() == "kicked":
            fib = "kicks"
        to_do.append((insertSql.format(fib), (uid, reason)))

    # insert existing removals
    for item in to_do:
        sql, args = item
        cursor.execute(sql, args)

    # delete old table
    cursor.execute("DROP TABLE `user_removes`")

    db.commit()

print("migrated db successfully to v2 format")
