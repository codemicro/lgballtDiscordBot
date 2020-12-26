import sqlite3
import json

with sqlite3.connect("lgballtBot.db") as db:	
    cursor = db.cursor()

    biodat = None
    with open("biosData.json", encoding="utf8") as f:
        biodat = json.load(f)["userBios"]

    for userid in biodat:
        dat = json.dumps(biodat[userid])
        print(userid, dat)
        insertData = '''INSERT INTO user_bios(user_id,raw_bio_data) VALUES(?,?)'''
        cursor.execute(insertData,[(userid),(dat)])
        db.commit()
