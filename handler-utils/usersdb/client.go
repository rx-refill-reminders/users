package usersdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rx-refill-reminders/users/handler-utils/model"
	"github.com/rx-refill-reminders/users/handler-utils/usersdb/dynamoiface"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")

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

	ConfirmUser(
		ctx context.Context,
		id string,
		confirmedAt time.Time,
	) error

	RecordLogin(
		ctx context.Context,
		id string,
		lastLogin time.Time,
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
	item, err := marshalMap(user)
	if err != nil {
		return fmt.Errorf("error marshaling user: %w", err)
	}

	requireNewCondition := expression.AttributeNotExists(expression.Name("id"))
	expr, err := expression.
		NewBuilder().
		WithCondition(requireNewCondition).
		Build()
	if err != nil {
		return fmt.Errorf("error constructing condition: %w", err)
	}

	_, err = c.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 &c.UsersTable,
		Item:                      item,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		var conditionalCheckFailed *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailed) {
			return fmt.Errorf("%w: %w", ErrUserAlreadyExists, err)
		}

		return fmt.Errorf("error writing new user to the database: %w", err)
	}

	return nil
}

func (c *client) ConfirmUser(
	ctx context.Context,
	id string,
	confirmedAt time.Time,
) error {
	ttl := model.UserRetentionTTL(confirmedAt)

	update := expression.
		Set(expression.Name("hasConfirmed"), expression.Value(true)).
		Set(expression.Name("updatedAt"), expression.Value(confirmedAt)).
		Set(expression.Name("ttl"), expression.Value(ttl))

	requireExistingCondition := expression.AttributeExists(expression.Name("id"))
	expr, err := expression.
		NewBuilder().
		WithUpdate(update).
		WithCondition(requireExistingCondition).
		Build()
	if err != nil {
		return fmt.Errorf("error constructing update: %w", err)
	}

	_, err = c.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &c.UsersTable,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("error confirming user: %w", err)
	}

	return nil
}

func (c *client) RecordLogin(
	ctx context.Context,
	id string,
	lastLogin time.Time,
) error {
	ttl := model.UserRetentionTTL(lastLogin)

	update := expression.
		Set(expression.Name("hasLoggedIn"), expression.Value(true)).
		Set(expression.Name("lastLogin"), expression.Value(lastLogin)).
		Set(expression.Name("updatedAt"), expression.Value(lastLogin)).
		Set(expression.Name("ttl"), expression.Value(ttl))

	requireExistingCondition := expression.AttributeExists(expression.Name("id"))
	expr, err := expression.
		NewBuilder().
		WithUpdate(update).
		WithCondition(requireExistingCondition).
		Build()
	if err != nil {
		return fmt.Errorf("error constructing update: %w", err)
	}

	_, err = c.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &c.UsersTable,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("error recording login: %w", err)
	}

	return nil
}
