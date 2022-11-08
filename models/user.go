package models

type InterfaceDBUser interface {
	GetUserById(id int) (*User, error)
	CreateUser(user *User) (*User, error)
	FindUserByEmail(email string) (*User, error)
}

type InterfaceUserService interface {
	GetUserById(id int) (*User, error)
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
	GetUserFromToken(tokenString string) (string, error)
}

type User struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string
}
