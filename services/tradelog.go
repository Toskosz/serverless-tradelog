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

func (s *logService) Delete(logId int) error {
	return s.dbConn.DeleteLog(logId)
}

func (s *logService) GetLog(id int) (*models.TradeLog, error) {
	return s.dbConn.GetLogById(id)
}

func (s *logService) GetUserLogs(userId int) ([]*models.TradeLog, error) {
	return s.dbConn.GetLogsByUserId(userId)
}
