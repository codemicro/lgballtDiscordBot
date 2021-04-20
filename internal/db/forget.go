package db

func ForgetUserID(userID string) error {
	mute := new(UserMute)
	mute.UserId = userID
	if err := mute.Delete(); err != nil {
		return err
	}

	return Conn.Exec(
		"DELETE FROM `user_bios` WHERE `user_bios`.`user_id` = ?",
		userID,
	).Error
}
