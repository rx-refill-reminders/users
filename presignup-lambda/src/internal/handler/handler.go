package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/handler-utils/config"
	"github.com/rx-refill-reminders/users/handler-utils/model"
	"github.com/rx-refill-reminders/users/handler-utils/usersdb"
)

const triggerSourcePrefix = "PreSignUp_"

type Handler interface {
	Handle(
		ctx context.Context,
		event events.CognitoEventUserPoolsPreSignup,
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
	event events.CognitoEventUserPoolsPreSignup,
) error {
	if !strings.HasPrefix(event.TriggerSource, triggerSourcePrefix) {
		return nil
	}

	attrs := event.Request.UserAttributes
	sub, ok := attrs["sub"]
	if !ok || sub == "" {
		return fmt.Errorf("missing sub in event.request.userAttributes")
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
		if errors.Is(err, usersdb.ErrUserAlreadyExists) {
			return fmt.Errorf("user with sub %q already exists: %w", sub, err)
		}

		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}
