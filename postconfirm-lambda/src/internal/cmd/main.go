package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rx-refill-reminders/lambda-go/logs"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/handler"
)

func run(
	ctx context.Context,
	event events.CognitoEventUserPoolsPostConfirmation,
) (events.CognitoEventUserPoolsPostConfirmation, error) {
	logger := logs.NewLogger(logs.LoggerOpts{
		Out: os.Stdout,
	})

	h, err := handler.NewHandlerFromEnv(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error initializing handler: %w", err)
		return event, nil
	}

	err = h.Handle(ctx, event)
	if err != nil {
		return event, err
	}

	return event, nil
}

func main() {
	lambda.Start(run)
}
