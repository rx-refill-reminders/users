package model

import "time"

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
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
