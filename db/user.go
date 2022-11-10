package db

import (
	"os"
	"strconv"

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

func (r *userRecords) FindUserByEmail(email string) (*models.User, error) {

	user := &models.User{}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(r.tableName),
	}

	userItem, err := r.DB.GetItem(input)
	if err != nil {
		return user, api_error.NewInternal()
	}

	if userItem.Item == nil {
		return user, api_error.NewNotFound("email", email)
	}

	return user, nil
}

func (r *userRecords) GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	userId := strconv.Itoa(id)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(userId),
			},
		},
		TableName: aws.String(r.tableName),
	}

	userItem, err := r.DB.GetItem(input)
	if err != nil {
		return user, api_error.NewInternal()
	}

	if userItem.Item == nil {
		return user, api_error.NewNotFound("Id", userId)
	}

	return user, nil
}

func (r *userRecords) CreateUser(user *models.User) (*models.User, error) {

	currentUser, err := r.FindUserByEmail(user.Email)
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
