package services

import (
	"time"

	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
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

	if notValidTimeStamp(log.TimestampAbertura) ||
		notValidTimeStamp(log.TimestampFechamento) {
		return nil, api_error.NewBadRequest(api_error.InvalidDateTime)
	}
	if notValidMEN(log.MEN) || notValidMEP(log.MEP) {
		return nil, api_error.NewBadRequest(api_error.InvalidMepMen)
	}

	return s.dbConn.CreateLog(log)
}

func (s *logService) Update(log *models.TradeLog) (*models.TradeLog, error) {
	return s.dbConn.UpdateLog(log)
}

func (s *logService) Delete(username string, aberturaTs string) error {
	if notValidTimeStamp(aberturaTs) {
		return api_error.NewBadRequest(api_error.InvalidDateTime)
	}
	return s.dbConn.DeleteLog(username, aberturaTs)
}

func (s *logService) GetLog(username string, aberturaTs string) (
	*models.TradeLog, error) {
	if notValidTimeStamp(aberturaTs) {
		return nil, api_error.NewBadRequest(api_error.InvalidDateTime)
	}
	return s.dbConn.GetLog(username, aberturaTs)
}

func (s *logService) GetUserLogs(username string) (*[]models.TradeLog, error) {
	return s.dbConn.GetLogsByUsername(username)
}

func notValidTimeStamp(datetime string) bool {
	_, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return true
	}
	return false
}

func notValidMEP(value float32) bool {
	if value < 0 {
		return true
	} else {
		return false
	}
}

func notValidMEN(value float32) bool {
	if value > 0 {
		return true
	} else {
		return false
	}
}
