package entity

type User struct {
	ID             int64   `gorm:"primary_key;type:bigserial"`
	Name           string  `gorm:"type:varchar(100)"`
	IdentityNumber string  `gorm:"unique;type:varchar(16)"`
	PhoneNumber    string  `gorm:"unique;type:varchar(15)"`
	Balance        float64 `gorm:"type:decimal(10,2);check:balance >= 0"`
	AccountNumber  string  `gorm:"unique;type:varchar(17)"`
}

func NewUser() *User {
	return &User{}
}

func (User) GetFieldFromConstraint(constraintName string) string {
	constraintToField := map[string]string{
		"uni_users_identity_number": "identity_number",
		"uni_users_phone_number":    "phone_number",
		"uni_users_account_number":  "account_number",
	}

	if field, exists := constraintToField[constraintName]; exists {
		return field
	}
	return "Unknown"
}
