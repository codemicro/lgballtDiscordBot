# Creating a bio

Unlike with regular bios for singlets, you need to explicitly create a bio for a system member before you can edit bio fields. This can be done using `$bio import <member ID>`, where `member ID` is the PluralKit member ID that you would like to import. If anything is in the birthday, pronouns or description fields, they will be automatically copied into the new bio.

# Updating bio fields

You can update and remove bio fields in a very similar fashion to bios for singlets, namely using the following two commands: `$bio clear <member ID> <field>` to remove a field and `$bio set <member ID> <field> <new contents>` to update a field's value. Field names are the same as bios for singlets, which can be found in `$bio help`.

# Deleting a system member's bio

To delete a system member bio, simply remove every field that's present in it using `$bio clear <member ID> <field>`. This will trigger the bot to automatically delete the bio entry from the database.

# Viewing system member bios

System member bios can be viewed using the same command as you would use to view singlet bios, namely `$bio [ping or user ID]`. If a user account has multiple bios associated with it, it will ask which bio you would like to view first, afterwards showing a carousel-type interface that will allowing you to scroll between bios.

# Anything else?

Be aware that any changes you make in a bio for a system member will not be replicated in the bio that PluralKit stores for that member. Using bios for systems doesn't affect your ability to use regular bios.

If you have questions, encounter an issue or have a suggestion, feel free to ping Abi! ({{ .Ping }})