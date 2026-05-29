package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/envconfig"
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
	envconfig.Config

	usersdb usersdb.Client
}

func NewHandlerFromEnv() (Handler, error) {
	cfg, err := envconfig.Load()
	if err != nil {
		return nil, fmt.Errorf("error parsing environment: %w", err)
	}

	return NewHandler(*cfg), nil
}

func NewHandler(cfg envconfig.Config) Handler {
	h := &handler{
		Config: cfg,

		usersdb: usersdb.NewClient(usersdb.Config{
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

		UserInfo: model.UserInfo{
			Email:     attrs["email"],
			FirstName: attrs["given_name"],
			LastName:  attrs["family_name"],
		},

		UserServerDrivenData: model.UserServerDrivenData{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	err := h.usersdb.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}
