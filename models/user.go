package models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `gorm:"not null;default:null"`
}
