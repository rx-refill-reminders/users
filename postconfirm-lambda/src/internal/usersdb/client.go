package usersdb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/dynamoiface"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/model"
)

var (
	marshalMap = attributevalue.MarshalMap
)

type Config struct {
	AWSConfig  aws.Config
	UsersTable string
}

type Client interface {
	CreateUser(
		ctx context.Context,
		user model.User,
	) error
}

type client struct {
	Config

	dynamo dynamoiface.Client
}

func NewClient(cfg Config) Client {
	return &client{
		Config: cfg,

		dynamo: dynamodb.NewFromConfig(cfg.AWSConfig),
	}
}

func (c *client) CreateUser(
	ctx context.Context,
	user model.User,
) error {
	// Convert the user model to a DynamoDB attributevalue map
	item, err := marshalMap(user)
	if err != nil {
		return fmt.Errorf("error marshaling user: %w", err)
	}

	// Construct an expression to ensure the new user will be unique
	requireNewCondition := expression.AttributeNotExists(expression.Name("id"))
	expr, err := expression.
		NewBuilder().
		WithCondition(requireNewCondition).
		Build()

	if err != nil {
		return fmt.Errorf("error constructing condition: %w", err)
	}

	// Try to write the new user to the database
	_, err = c.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 &c.UsersTable,
		Item:                      item,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		return fmt.Errorf("error writing new user to the database: %w", err)
	}

	return nil
}
