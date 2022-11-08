package services

import (
	"github.com/Toskosz/everythingreviewed/models"
)

type logService struct {
	dbConn models.InterfaceDBLog
}

func NewLogService(conn models.InterfaceDBLog) models.InterfaceLogService {
	return &logService{
		dbConn: conn,
	}
}

func (s *logService) Create(log *models.TradeLog) (*models.TradeLog, error) {
	return s.dbConn.CreateLog(log)
}

func (s *logService) Update(log *models.TradeLog) (*models.TradeLog, error) {
	return s.dbConn.UpdateLog(log)
}

func (s *logService) Delete(aberturaTs string) error {
	return s.dbConn.DeleteLog(aberturaTs)
}

func (s *logService) GetLog(username string, aberturaTs string) (*models.TradeLog, error) {
	return s.dbConn.GetLog(username, aberturaTs)
}

func (s *logService) GetUserLogs(username string) (*[]models.TradeLog, error) {
	return s.dbConn.GetLogsByUsername(username)
}
