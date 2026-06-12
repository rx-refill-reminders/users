package handler

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/config"
	"github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/model"
	usersdbmocks "github.com/rx-refill-reminders/users/postconfirm-lambda/src/internal/usersdb/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type HandlerTestHarness struct {
	*handler

	MockUsersDB *usersdbmocks.Client
}

func NewHandlerTestHarness(
	t *testing.T,
	cfg config.Config,
) *HandlerTestHarness {
	h := &HandlerTestHarness{}

	h.MockUsersDB = usersdbmocks.NewClient(t)

	h.handler = &handler{
		Config: cfg,

		usersdb: h.MockUsersDB,
	}

	return h
}

func confirmSignUpEvent(
	attrs map[string]string,
) events.CognitoEventUserPoolsPostConfirmation {
	return events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: triggerSourceConfirmSignUp,
		},
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: attrs,
		},
	}
}

func TestHandler_Handle(t *testing.T) {
	errInjected := fmt.Errorf("injected")

	t.Run("success", func(t *testing.T) {
		// Create a new HandlerTestHarness
		h := NewHandlerTestHarness(t, config.Config{})

		// Create an example Cognito post-confirmation event
		event := events.CognitoEventUserPoolsPostConfirmation{
			CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
				TriggerSource: triggerSourceConfirmSignUp,
			},
			Request: events.CognitoEventUserPoolsPostConfirmationRequest{
				UserAttributes: map[string]string{
					"sub":         "user-sub-123",
					"email":       "test@example.com",
					"given_name":  "Jane",
					"family_name": "Doe",
				},
			},
		}

		// Mock the usersdb CreateUser call
		h.MockUsersDB.EXPECT().
			CreateUser(mock.Anything, mock.MatchedBy(func(user model.User) bool {
				return user.ID == "user-sub-123" &&
					user.Email == "test@example.com" &&
					user.FirstName == "Jane" &&
					user.LastName == "Doe" &&
					!user.CreatedAt.IsZero() &&
					!user.UpdatedAt.IsZero()
			})).
			Return(nil)

		// Call Handle
		err := h.Handle(t.Context(), event)

		// Expect the call to succeed
		require.NoError(t, err)
	})

	t.Run("skip-wrong-trigger-source", func(t *testing.T) {
		// Create a new HandlerTestHarness
		h := NewHandlerTestHarness(t, config.Config{})

		// Create an event with an unexpected trigger source
		event := events.CognitoEventUserPoolsPostConfirmation{
			CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
				TriggerSource: "PostConfirmation_SomethingElse",
			},
			Request: events.CognitoEventUserPoolsPostConfirmationRequest{
				UserAttributes: map[string]string{
					"sub": "user-sub-123",
				},
			},
		}

		// Call Handle
		err := h.Handle(t.Context(), event)

		// Expect the call to succeed without creating a user
		require.NoError(t, err)
	})

	t.Run("err-missing-sub", func(t *testing.T) {
		// Create a new HandlerTestHarness
		h := NewHandlerTestHarness(t, config.Config{})

		// Create an event without a sub attribute
		event := events.CognitoEventUserPoolsPostConfirmation{
			CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
				TriggerSource: triggerSourceConfirmSignUp,
			},
			Request: events.CognitoEventUserPoolsPostConfirmationRequest{
				UserAttributes: map[string]string{
					"email": "test@example.com",
				},
			},
		}

		// Call Handle
		err := h.Handle(t.Context(), event)

		// Expect the call to fail
		require.ErrorContains(t, err, "missing sub")
	})

	t.Run("err-create-user-failed", func(t *testing.T) {
		// Create a new HandlerTestHarness
		h := NewHandlerTestHarness(t, config.Config{})

		// Create an example Cognito post-confirmation event
		event := confirmSignUpEvent(map[string]string{
			"sub": "user-sub-123",
		})

		// Mock the usersdb CreateUser call (errInjected)
		h.MockUsersDB.EXPECT().
			CreateUser(mock.Anything, mock.Anything).
			Return(errInjected)

		// Call Handle
		err := h.Handle(t.Context(), event)

		// Expect the call to fail, and the result to be errInjected
		require.ErrorIs(t, err, errInjected)
	})
}
