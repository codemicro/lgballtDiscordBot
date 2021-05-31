import sqlite3
import sys

db_location = sys.argv[1]

createNewToneTagTable = """CREATE TABLE IF NOT EXISTS `tone_tags` (`shorthand` text,`description` text,PRIMARY KEY (`shorthand`))"""
insertSql = "INSERT INTO `tone_tags` (`shorthand`, `description`) VALUES(?, ?)"

toneTags = {
    "j": "joking",
    "hj": "half joking",
    "s": "sarcastic/sarcasm",
    "sarc": "sarcastic/sarcasm",
    "srs": "serious",
    "nsrs": "not serious",
    "lh": "light hearted",
    "g": "genuine/genuine question",
    "gen": "genuine/genuine question",
    "ij": "inside joke",
    "ref": "reference",
    "t": "teasing",
    "nm": "not mad",
    "lu": "a little upset",
    "nbh": "nobody here",
    "nsb": "not subtweeting",
    "nay": "not at you",
    "ay": "at you",
    "nbr": "not being rude",
    "ot": "off topic",
    "th": "threat",
    "cb": "clickbait",
    "f": "fake",
    "q": "quote",
    "l": "lyrics",
    "ly": "lyrics",
    "c": "copypasta",
    "m": "metaphor/metaphorically",
    "li": "literal/literally",
    "rt": "rhetorical question",
    "rh": "rhetorical question",
    "hyp": "hyperbole",
    "p": "platonic",
    "r": "romantic",
    "a": "alterous",
    "pc": "positive connotation",
    "pos": "positive connotation",
    "nc": "negative connotation",
    "neg": "negative connotation",
    "neu": "neutral/neutral connotation"
}

with sqlite3.connect(db_location) as db:	
    cursor = db.cursor()
    
    # add system ID field to bios

    # create new table
    cursor.execute(createNewToneTagTable)
    db.commit()

    for shorthand in toneTags:
        description = toneTags[shorthand]

        try:
            cursor.execute(insertSql, (shorthand, description))
        except sqlite3.IntegrityError:
            print("Unable to insert {} (sqlite3.IntegrityError)".format(shorthand))

    db.commit()

print("migrated db successfully to v4.7.x format")
