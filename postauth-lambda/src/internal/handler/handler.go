package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/handler-utils/config"
	"github.com/rx-refill-reminders/users/handler-utils/usersdb"
)

const triggerSourcePrefix = "PostAuthentication_"

type Handler interface {
	Handle(
		ctx context.Context,
		event events.CognitoEventUserPoolsPostAuthentication,
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
	event events.CognitoEventUserPoolsPostAuthentication,
) error {
	if !strings.HasPrefix(event.TriggerSource, triggerSourcePrefix) {
		return nil
	}

	attrs := event.Request.UserAttributes
	sub, ok := attrs["sub"]
	if !ok || sub == "" {
		return fmt.Errorf("missing sub in event.request.userAttributes")
	}

	lastLogin := time.Now().UTC()

	err := h.usersdb.RecordLogin(ctx, sub, lastLogin)
	if err != nil {
		return fmt.Errorf("error recording login: %w", err)
	}

	return nil
}
