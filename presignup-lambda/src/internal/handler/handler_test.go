package handler

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rx-refill-reminders/users/handler-utils/config"
	"github.com/rx-refill-reminders/users/handler-utils/model"
	"github.com/rx-refill-reminders/users/handler-utils/usersdb"
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

func signUpEvent(
	attrs map[string]string,
) events.CognitoEventUserPoolsPreSignup {
	return events.CognitoEventUserPoolsPreSignup{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: "PreSignUp_SignUp",
		},
		Request: events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: attrs,
		},
	}
}

func TestHandler_Handle(t *testing.T) {
	errInjected := fmt.Errorf("injected")

	t.Run("success", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := signUpEvent(map[string]string{
			"sub":         "user-sub-123",
			"email":       "test@example.com",
			"given_name":  "Jane",
			"family_name": "Doe",
		})

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

		err := h.Handle(t.Context(), event)

		require.NoError(t, err)
	})

	t.Run("skip-wrong-trigger-source", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := events.CognitoEventUserPoolsPreSignup{
			CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
				TriggerSource: "SomethingElse",
			},
			Request: events.CognitoEventUserPoolsPreSignupRequest{
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

		event := signUpEvent(map[string]string{
			"email": "test@example.com",
		})

		err := h.Handle(t.Context(), event)

		require.ErrorContains(t, err, "missing sub")
	})

	t.Run("err-user-already-exists", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := signUpEvent(map[string]string{
			"sub": "user-sub-123",
		})

		h.MockUsersDB.EXPECT().
			CreateUser(mock.Anything, mock.Anything).
			Return(fmt.Errorf("%w: injected", usersdb.ErrUserAlreadyExists))

		err := h.Handle(t.Context(), event)

		require.ErrorIs(t, err, usersdb.ErrUserAlreadyExists)
	})

	t.Run("err-create-user-failed", func(t *testing.T) {
		h := NewHandlerTestHarness(t, config.Config{})

		event := signUpEvent(map[string]string{
			"sub": "user-sub-123",
		})

		h.MockUsersDB.EXPECT().
			CreateUser(mock.Anything, mock.Anything).
			Return(errInjected)

		err := h.Handle(t.Context(), event)

		require.ErrorIs(t, err, errInjected)
	})
}
