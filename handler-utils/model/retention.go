package model

import (
	"time"

	dynamoutils "github.com/lucaspopp0/go-dynamo-utils"
)

const UserRetention = 20 * 365 * 24 * time.Hour

func UserRetentionTTL(from time.Time) dynamoutils.TTL {
	return dynamoutils.NewTTL(from.Add(UserRetention))
}
