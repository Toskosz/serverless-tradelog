package models

type InterfaceDBUser interface {
	CreateUser(user *User) (*User, error)
	FetchUserByUsername(email string) (*User, error)
}

type InterfaceUserService interface {
	GetUserByUsername(username string) (*User, error)
	Register(user *User) (*User, error)
	Login(username string, password string) (*User, error)
	GetUserFromToken(tokenString string) (string, error)
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string
}
