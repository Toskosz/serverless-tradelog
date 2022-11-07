package models

import "time"

type Timestamp time.Time

type InterfaceDBLog interface {
	GetUserById(id int) (*User, error)
	CreateUser(user *User) (*User, error)
	FindUserByEmail(email string) (*User, error)
}

type InterfaceLogService interface {
	GetUserById(id int) (*User, error)
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
}

type TradeLog struct {
	Id                  uint      `json:"id"`
	Userid              string    `json:"user-id" gorm:"unique"`
	Ativo               string    `json:"ativo"`
	Resultado           string    `json:"resultado"`
	Contratos           int       `json:"contratos"`
	MEP                 string    `json:"mep"`
	MEN                 string    `json:"men"`
	TempoOperacao       float64   `json:"duracao"`
	PrecoCompra         float32   `json:"preco-compra"`
	PrecoVenda          float32   `json:"preco-venda"`
	TimestampAbertura   Timestamp `json:"hora-abertura"`
	TimestampFechamento Timestamp `json:"hora-fechamento"`
}
