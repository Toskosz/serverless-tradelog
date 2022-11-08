package db

import (
	"os"

	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type logRecords struct {
	DB        *dynamodb.DynamoDB
	tableName string
}

func NewTradeLogDBConn(table string) models.InterfaceDBLog {
	return &logRecords{
		DB: dynamodb.New(
			session.New(),
			aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")),
		),
		tableName: table,
	}
}

func (r *logRecords) GetLog(username string, aberturaTs string) (
	*models.TradeLog, error) {
	log := &models.TradeLog{}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				N: aws.String(username),
			},
			"Abertura": {
				N: aws.String(aberturaTs),
			},
		},
		TableName: aws.String(r.tableName),
	}

	logItem, err := r.DB.GetItem(input)
	if err != nil {
		return log, api_error.NewInternal()
	}

	if logItem.Item == nil {
		return log, api_error.NewNotFound("log abertura", aberturaTs)
	}

	return log, nil
}

func (r *logRecords) GetLogsByUsername(username string) (
	*[]models.TradeLog, error) {

	filt := expression.Name("Username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, api_error.NewInternal()
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(r.tableName),
	}

	// Make the DynamoDB Query API call
	result, err := r.DB.Scan(params)
	if err != nil {
		return nil, api_error.NewInternal()
	}

	item := new([]models.TradeLog)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func (r *logRecords) CreateLog(log *models.TradeLog) (*models.TradeLog, error) {

	dynamoItem, err := dynamodbattribute.MarshalMap(log)
	if err != nil {
		return nil, api_error.NewInternal()
	}
	dynamoInput := &dynamodb.PutItemInput{
		Item:      dynamoItem,
		TableName: aws.String(r.tableName),
	}
	_, err = r.DB.PutItem(dynamoInput)
	if err != nil {
		return nil, api_error.NewInternal()
	}

	return log, nil
}

func (r *logRecords) UpdateLog(log *models.TradeLog) (*models.TradeLog, error) {
	currentLog, _ := r.GetLog(log.Username, log.TimestampAbertura)
	if currentLog != nil && len(currentLog.Ativo) == 0 {
		return nil,
			api_error.NewNotFound("Log abertura", log.TimestampAbertura)
	}

	av, err := dynamodbattribute.MarshalMap(log)
	if err != nil {
		return nil, api_error.NewBadRequest("Unable to marshall")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.tableName),
	}

	_, err = r.DB.PutItem(input)
	if err != nil {
		return nil, api_error.NewInternal()
	}
	return log, nil
}

func (r *logRecords) DeleteLog(aberturaTs string) error {

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				N: aws.String(aberturaTs),
			},
			"Abertura": {
				N: aws.String(aberturaTs),
			},
		},
		TableName: aws.String(r.tableName),
	}
	itemDeleted, err := r.DB.DeleteItem(input)

	// Failed to delete
	if err != nil {
		return api_error.NewInternal()
	}

	// Tried to delete item not present
	if itemDeleted == nil {
		return api_error.NewNotFound("Log abertura", aberturaTs)
	}

	return nil
}
