package db

import (
	"os"
	"strconv"

	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))

const tableName = ""

func FindUserByEmail(email string) (*models.User, error) {

	user := &models.User{}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	userItem, err := db.GetItem(input)
	if err != nil {
		return user, api_error.NewInternal()
	}

	if userItem.Item == nil {
		return user, api_error.NewNotFound("email", email)
	}

	return user, nil
}

func GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	userId := strconv.Itoa(id)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(userId),
			},
		},
		TableName: aws.String(tableName),
	}

	userItem, err := db.GetItem(input)
	if err != nil {
		return user, api_error.NewInternal()
	}

	if userItem.Item == nil {
		return user, api_error.NewNotFound("Id", userId)
	}

	return user, nil
}

func CreateUser(user *models.User) (*models.User, error) {

	currentUser, err := FindUserByEmail(user.Email)
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
		TableName: aws.String(tableName),
	}
	_, err = db.PutItem(dynamoInput)
	if err != nil {
		return nil, api_error.NewInternal()
	}

	return user, nil
}
