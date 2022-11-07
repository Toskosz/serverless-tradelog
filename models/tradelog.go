package models

import "time"

type Timestamp time.Time

type InterfaceDBLog interface {
	GetLogById(id int) (*TradeLog, error)
	GetLogsByUserId(userId int) ([]*TradeLog, error)
	CreateLog(log *TradeLog) (*TradeLog, error)
	UpdateLog(log *TradeLog) (*TradeLog, error)
	DeleteLog(logId int) error
}

type InterfaceLogService interface {
	GetLog(id int) (*TradeLog, error)
	GetUserLogs(userId int) ([]*TradeLog, error)
	Create(log *TradeLog) (*TradeLog, error)
	Update(log *TradeLog) (*TradeLog, error)
	Delete(logId int) error
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
	Desc                string    `json:"descricao"`
}
