package model

import (
	"time"

	dynamoutils "github.com/lucaspopp0/go-dynamo-utils"
)

type User struct {
	ID string `json:"id" dynamodbav:"id"`

	UserInfo

	UserMetadata
}

type UserInfo struct {
	Email string `json:"email" dynamodbav:"email"`

	FirstName string `json:"firstName" dynamodbav:"firstName"`
	LastName  string `json:"lastName" dynamodbav:"lastName"`
}

type UserMetadata struct {
	HasConfirmed bool `json:"hasConfirmed" dynamodbav:"hasConfirmed"`
	HasLoggedIn  bool `json:"hasLoggedIn" dynamodbav:"hasLoggedIn"`

	LastLogin time.Time `json:"lastLogin" dynamodbav:"lastLogin"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`

	TTL dynamoutils.TTL `json:"ttl" dynamodbav:"ttl"`
}
