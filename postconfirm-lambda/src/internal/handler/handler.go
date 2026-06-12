package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/config"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/model"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/usersdb"
)

const triggerSourceConfirmSignUp = "PostConfirmation_ConfirmSignUp"

type Handler interface {
	Handle(
		ctx context.Context,
		event events.CognitoEventUserPoolsPostConfirmation,
	) error
}

type handler struct {
	config.Config

	usersdb usersdb.Client
}

func NewHandlerFromEnv(ctx context.Context) (Handler, error) {
	cfg, err := config.Load(ctx)
	if err != nil {
		return nil, err
	}

	return NewHandler(*cfg), nil
}

func NewHandler(
	cfg config.Config,
) Handler {
	h := &handler{
		Config: cfg,

		usersdb: usersdb.NewClient(usersdb.Config{
			AWSConfig:  cfg.AWSConfig,
			UsersTable: cfg.UsersTable,
		}),
	}

	return h
}

func (h *handler) Handle(
	ctx context.Context,
	event events.CognitoEventUserPoolsPostConfirmation,
) error {
	if event.TriggerSource != triggerSourceConfirmSignUp {
		return nil
	}

	attrs := event.Request.UserAttributes
	sub, ok := attrs["sub"]
	if !ok || sub == "" {
		return fmt.Errorf("missing sub in event.request.userAttributes, skipping")
	}

	now := time.Now().UTC()
	user := model.User{
		ID: sub,

		UserMetadata: model.UserMetadata{
			CreatedAt: now,
			UpdatedAt: now,
		},

		UserInfo: model.UserInfo{
			Email:     attrs["email"],
			FirstName: attrs["given_name"],
			LastName:  attrs["family_name"],
		},
	}

	err := h.usersdb.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}
