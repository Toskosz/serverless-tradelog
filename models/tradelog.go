package models

type InterfaceDBLog interface {
	GetLog(username string, aberturaTs string) (*TradeLog, error)
	GetLogsByUsername(username string) (*[]TradeLog, error)
	CreateLog(log *TradeLog) (*TradeLog, error)
	UpdateLog(log *TradeLog) (*TradeLog, error)
	DeleteLog(username string, aberturaTs string) error
}

type InterfaceLogService interface {
	GetLog(username string, aberturaTs string) (*TradeLog, error)
	GetUserLogs(username string) (*[]TradeLog, error)
	Create(log *TradeLog) (*TradeLog, error)
	Update(log *TradeLog) (*TradeLog, error)
	Delete(username string, aberturaTs string) error
}

type TradeLog struct {
	Username              string  `json:"user-id"`
	TimestampAbertura     string  `json:"abertura"`
	TimestampFechamento   string  `json:"fechamento"`
	Ativo                 string  `json:"ativo"`
	Resultado             float32 `json:"resultado"`
	Contratos             int     `json:"contratos"`
	MEP                   float32 `json:"mep"`
	MEN                   float32 `json:"men"`
	TempoOperacaoSegundos int     `json:"duracao"`
	PrecoCompra           float32 `json:"preco-compra"`
	PrecoVenda            float32 `json:"preco-venda"`
	Revisado              bool    `json:"revisado"`
	Desc                  string  `json:"descricao"`
}
