package models

type InterfaceDBLog interface {
	GetLog(username string, aberturaTs string) (*TradeLog, error)
	GetLogsByUsername(username string) (*[]TradeLog, error)
	CreateLog(log *TradeLog) (*TradeLog, error)
	UpdateLog(log *TradeLog) (*TradeLog, error)
	DeleteLog(aberturaTs string) error
}

type InterfaceLogService interface {
	GetLog(username string, aberturaTs string) (*TradeLog, error)
	GetUserLogs(username string) (*[]TradeLog, error)
	Create(log *TradeLog) (*TradeLog, error)
	Update(log *TradeLog) (*TradeLog, error)
	Delete(aberturaTs string) error
}

type TradeLog struct {
	Username            string  `json:"user-id"`
	TimestampAbertura   string  `json:"hora-abertura"`
	TimestampFechamento string  `json:"hora-fechamento"`
	Ativo               string  `json:"ativo"`
	Resultado           string  `json:"resultado"`
	Contratos           int     `json:"contratos"`
	MEP                 string  `json:"mep"`
	MEN                 string  `json:"men"`
	TempoOperacao       int     `json:"duracao"`
	PrecoCompra         float32 `json:"preco-compra"`
	PrecoVenda          float32 `json:"preco-venda"`
	Desc                string  `json:"descricao"`
}
