package db

import (
	"os"

	"github.com/Toskosz/serverless-tradelog/models"
	"github.com/Toskosz/serverless-tradelog/models/api_error"
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
			"username": {
				N: aws.String(username),
			},
			"abertura": {
				N: aws.String(aberturaTs),
			},
		},
		TableName: aws.String(r.tableName),
	}

	logItem, err := r.DB.GetItem(input)
	if err != nil {
		return log, api_error.NewUnsupportedMediaType(err.Error())
	}

	if logItem.Item == nil {
		return log, api_error.NewNotFound("log abertura", aberturaTs)
	}

	return log, nil
}

func (r *logRecords) GetLogsByUsername(username string) (
	*[]models.TradeLog, error) {

	filt := expression.Name("username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, api_error.NewUnsupportedMediaType(err.Error())
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
		return nil, api_error.NewUnsupportedMediaType(err.Error())
	}

	item := new([]models.TradeLog)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func (r *logRecords) CreateLog(log *models.TradeLog) (*models.TradeLog, error) {

	dynamoItem, err := dynamodbattribute.MarshalMap(log)
	if err != nil {
		return nil, api_error.NewUnsupportedMediaType(err.Error())
	}
	dynamoInput := &dynamodb.PutItemInput{
		Item:      dynamoItem,
		TableName: aws.String(r.tableName),
	}
	_, err = r.DB.PutItem(dynamoInput)
	if err != nil {
		return nil, api_error.NewUnsupportedMediaType(err.Error())
	}

	return log, nil
}

func (r *logRecords) UpdateLog(log *models.TradeLog) (*models.TradeLog, error) {

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				BOOL: aws.Bool(log.Revisado),
			},
			":d": {
				S: aws.String(log.Desc),
			},
		},
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(log.Username),
			},
			"abertura": {
				S: aws.String(log.TimestampAbertura),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET Revisado = :r, Desc = :d"),
	}

	_, err := r.DB.UpdateItem(input)
	if err != nil {
		return nil, api_error.NewUnsupportedMediaType(err.Error())
	}

	return log, nil
}

func (r *logRecords) DeleteLog(username string, aberturaTs string) error {

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				N: aws.String(username),
			},
			"abertura": {
				N: aws.String(aberturaTs),
			},
		},
		TableName: aws.String(r.tableName),
	}
	itemDeleted, err := r.DB.DeleteItem(input)

	// Failed to delete
	if err != nil {
		return api_error.NewUnsupportedMediaType(err.Error())
	}

	// Tried to delete item not present
	if itemDeleted == nil {
		return api_error.NewNotFound("Log abertura", aberturaTs)
	}

	return nil
}
