package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/Toskosz/everythingreviewed/services"
	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) GetLog(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	abertura := req.PathParameters["log-abertura"]
	tokenStr, _ := req.Headers["jwt"]
	requesterUsername, err := h.userService.GetUserFromToken(tokenStr)

	log, err := h.logService.GetLog(requesterUsername, abertura)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	return services.ApiResponse(http.StatusOK, log)
}

func (h *Handler) GetMyLogs(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	// If the requests get here then we are sure the token exists
	// because of the middleware
	tokenStr, _ := req.Headers["jwt"]

	username, err := h.userService.GetUserFromToken(tokenStr)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	logs, err := h.logService.GetUserLogs(username)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	return services.ApiResponse(http.StatusOK, logs)
}

type createLogInput struct {
	TimestampAbertura     string  `json:"hora-abertura"`
	TimestampFechamento   string  `json:"hora-fechamento"`
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

func (i *createLogInput) sanitize() {
	i.Ativo = strings.TrimSpace(i.Ativo)
	i.TimestampAbertura = strings.TrimSpace(i.TimestampAbertura)
	i.TimestampFechamento = strings.TrimSpace(i.TimestampFechamento)
	i.Desc = strings.TrimSpace(i.Desc)
}

func (h *Handler) CreateLog(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	tokenStr, _ := req.Headers["jwt"]
	requesterUsername, err := h.userService.GetUserFromToken(tokenStr)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	var jsonData createLogInput

	if err := json.Unmarshal([]byte(req.Body), &jsonData); err != nil {
		return services.ApiResponse(http.StatusBadRequest,
			api_error.NewBadRequest("Bad payload"))
	}
	jsonData.sanitize()

	registerTradeLogPayload := &models.TradeLog{
		Username:              requesterUsername,
		Ativo:                 jsonData.Ativo,
		Resultado:             jsonData.Resultado,
		Contratos:             jsonData.Contratos,
		MEP:                   jsonData.MEP,
		MEN:                   jsonData.MEN,
		TempoOperacaoSegundos: jsonData.TempoOperacaoSegundos,
		PrecoCompra:           jsonData.PrecoCompra,
		PrecoVenda:            jsonData.PrecoVenda,
		TimestampAbertura:     jsonData.TimestampAbertura,
		TimestampFechamento:   jsonData.TimestampFechamento,
		Revisado:              jsonData.Revisado,
		Desc:                  jsonData.Desc,
	}

	log, err := h.logService.Create(registerTradeLogPayload)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	return services.ApiResponse(http.StatusCreated, log)
}

type updateLogInput struct {
	TimestampAbertura string `json:"hora-abertura"`
	Revisado          bool   `json:"revisado"`
	Desc              string `json:"descricao"`
}

func (i *updateLogInput) sanitize() {
	i.TimestampAbertura = strings.TrimSpace(i.TimestampAbertura)
	i.Desc = strings.TrimSpace(i.Desc)
}

func (h *Handler) UpdateLog(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	// get user
	tokenStr, _ := req.Headers["jwt"]
	requesterUsername, err := h.userService.GetUserFromToken(tokenStr)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	// map input data
	var jsonData updateLogInput
	if err := json.Unmarshal([]byte(req.Body), &jsonData); err != nil {
		return services.ApiResponse(http.StatusBadRequest,
			api_error.NewBadRequest("Bad payload"))
	}
	jsonData.sanitize()

	// check if log even exists
	currentLog, err := h.logService.GetLog(requesterUsername, jsonData.TimestampAbertura)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	currentLog.Revisado = jsonData.Revisado
	currentLog.Desc = jsonData.Desc

	newLog, err := h.logService.Update(currentLog)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	return services.ApiResponse(http.StatusOK, newLog)
}

func (h *Handler) DeleteLog(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	// get user
	tokenStr, _ := req.Headers["jwt"]
	requesterUsername, err := h.userService.GetUserFromToken(tokenStr)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	abertura := req.PathParameters["aberturaTs"]
	err = h.logService.Delete(requesterUsername, abertura)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	return services.ApiResponse(http.StatusOK, nil)
}
