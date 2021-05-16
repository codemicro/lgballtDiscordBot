# What are bios?

Think of Bios like ID cards. Input responses for pre-defined fields such as a Pronouns field or a Sexuality field, and then your responses get put into a nice Bio card which can be viewed by anyone using `$bio [@username]`.

If you just want to get your own, start by filling in any of the below fields!

# How do I fill out these fields?

Run `$bio set [field] [value]`. For example, `$bio set Pronouns She/Her` would set the Pronouns field of your Bio to "She/Her".

To remove a field, run `$bio clear [field]` with no other arguments.

# What fields can I fill in?

The current fields are:

```yml
{{ .Fields }}
```

# What about bios for systems?

I'm glad you asked! Run `$bio syshelp` for more information about bios for systems.

# I think a new field should be added. How can I request one?

Request new fields in <#698575463278313583>.

# How do I view someone else's bio without pinging them?

User IDs can be used instead of mentioning a user. To get a User ID, first enable Developer Mode by going to User Settings, Appearance, and toggling it to on. After that, right click a username on desktop or tap the 3 dots on a profile card on mobile then click "Copy ID".
Now just run `$bio [UserID]` to view their Bio. For example, `$bio 516962733497778176`

# Anything else I should know?

- You don't need to wipe a field to put in new info. Just run `$bio set [field] [text]` to overwrite it.
- If you end up in a situation where you have no fields left in your bio because you've removed them all, your entire bio is deleted.

# TL;DR/Commands

- View your own Bio with `$bio`, another user's with `$bio [user id or mention]`
- Fill in a field with `$bio set [field] [text]`. Fields can be overwritten with the same command.
- Wipe a field with `$bio clear [field]`