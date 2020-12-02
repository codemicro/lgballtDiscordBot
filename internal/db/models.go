package db

type UserBio struct {
	UserId string `gorm:"primarykey"`
	RawBioData string
	BioData map[string]string `gorm:"-"`
}
