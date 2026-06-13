package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/handler-utils/config"
	usersdbmocks "github.com/rx-refill-reminders/users/handler-utils/usersdb/mocks"
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

func authenticationEvent(
	attrs map[string]string,
) events.CognitoEventUserPoolsPostAuthentication {
	return events.CognitoEventUserPoolsPostAuthentication{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: "PostAuthentication_Authentication",
		},
		Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
			UserAttributes: attrs,
		},
	}
}

func TestHandler_Handle(t *testing.T) {
	errInjected := fmt.Errorf("injected")

	t.Run("success", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := authenticationEvent(map[string]string{
			"sub": "user-sub-123",
		})

		h.MockUsersDB.EXPECT().
			RecordLogin(mock.Anything, "user-sub-123", mock.MatchedBy(func(ts time.Time) bool {
				return !ts.IsZero()
			})).
			Return(nil)

		err := h.Handle(t.Context(), event)

		require.NoError(t, err)
	})

	t.Run("skip-wrong-trigger-source", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := events.CognitoEventUserPoolsPostAuthentication{
			CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
				TriggerSource: "SomethingElse",
			},
			Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
				UserAttributes: map[string]string{
					"sub": "user-sub-123",
				},
			},
		}

		err := h.Handle(t.Context(), event)

		require.NoError(t, err)
	})

	t.Run("err-missing-sub", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := authenticationEvent(map[string]string{
			"email": "test@example.com",
		})

		err := h.Handle(t.Context(), event)

		require.ErrorContains(t, err, "missing sub")
	})

	t.Run("err-record-login-failed", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := authenticationEvent(map[string]string{
			"sub": "user-sub-123",
		})

		h.MockUsersDB.EXPECT().
			RecordLogin(mock.Anything, mock.Anything, mock.Anything).
			Return(errInjected)

		err := h.Handle(t.Context(), event)

		require.ErrorIs(t, err, errInjected)
	})
}
