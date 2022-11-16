package db

import (
	"os"

	"github.com/Toskosz/serverless-tradelog/models"
	"github.com/Toskosz/serverless-tradelog/models/api_error"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type userRecords struct {
	DB        *dynamodb.DynamoDB
	tableName string
}

func NewUserDBConn(table string) models.InterfaceDBUser {
	return &userRecords{
		DB: dynamodb.New(
			session.New(),
			aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")),
		),
		tableName: table,
	}
}

func (r *userRecords) FetchUserByUsername(username string) (*models.User, error) {

	user := &models.User{}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
			},
		},
		TableName: aws.String(r.tableName),
	}

	userItem, err := r.DB.GetItem(input)
	if err != nil {
		return user, api_error.NewInternal()
	}

	if userItem.Item == nil {
		return user, api_error.NewNotFound("username", username)
	}

	return user, nil
}

func (r *userRecords) CreateUser(user *models.User) (*models.User, error) {

	currentUser, err := r.FetchUserByUsername(user.Username)
	if err == nil {
		if currentUser != nil && len(currentUser.Email) != 0 {
			return nil, api_error.NewBadRequest(api_error.DuplicateEmailError)
		}
	}
	if err.Error() == api_error.InternalError {
		return nil, api_error.NewInternal()
	}

	dynamoItem, err := dynamodbattribute.MarshalMap(user)
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

	return user, nil
}
