package entity

type User struct {
	ID             int     `gorm:"primary_key;type:serial"`
	Name           string  `gorm:"type:varchar(100)"`
	IdentityNumber string  `gorm:"unique;type:varchar(16)"`
	PhoneNumber    string  `gorm:"unique;type:varchar(15)"`
	Balance        float64 `gorm:"type:decimal(10,2);check:balance >= 0"`
	AccountNumber  string  `gorm:"unique;type:varchar(17)"`
}

func NewUser() *User {
	return &User{}
}
