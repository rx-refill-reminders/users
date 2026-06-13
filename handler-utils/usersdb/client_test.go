package usersdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rx-refill-reminders/users/handler-utils/model"
	dynamoifaceMocks "github.com/rx-refill-reminders/users/handler-utils/usersdb/dynamoiface/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ClientTestHarness struct {
	*client

	MockDynamo *dynamoifaceMocks.Client
}

func NewClientTestHarness(
	t *testing.T,
	cfg Config,
) *ClientTestHarness {
	h := &ClientTestHarness{}

	h.MockDynamo = dynamoifaceMocks.NewClient(t)

	h.client = &client{
		Config: cfg,

		dynamo: h.MockDynamo,
	}

	return h
}

func TestClient_CreateUser(t *testing.T) {
	errInjected := fmt.Errorf("injected")

	t.Run("success", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		user := model.User{
			ID: uuid.NewString(),
		}

		h.MockDynamo.EXPECT().
			PutItem(mock.Anything, mock.Anything).
			Return(nil, nil)

		err := h.CreateUser(t.Context(), user)

		require.NoError(t, err)
	})

	t.Run("err-user-already-exists", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		user := model.User{
			ID: uuid.NewString(),
		}

		h.MockDynamo.EXPECT().
			PutItem(mock.Anything, mock.Anything).
			Return(nil, &types.ConditionalCheckFailedException{})

		err := h.CreateUser(t.Context(), user)

		require.ErrorIs(t, err, ErrUserAlreadyExists)
	})

	t.Run("err-marshal-failed", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		user := model.User{
			ID: uuid.NewString(),
		}

		originalMarshalMap := marshalMap
		t.Cleanup(func() { marshalMap = originalMarshalMap })
		marshalMap = func(_ any) (map[string]types.AttributeValue, error) {
			return nil, errInjected
		}

		err := h.CreateUser(t.Context(), user)

		require.ErrorIs(t, err, errInjected)
	})

	t.Run("err-put-item-failed", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		user := model.User{
			ID: uuid.NewString(),
		}

		h.MockDynamo.EXPECT().
			PutItem(mock.Anything, mock.Anything).
			Return(nil, errInjected)

		err := h.CreateUser(t.Context(), user)

		require.ErrorIs(t, err, errInjected)
	})
}

func TestClient_ConfirmUser(t *testing.T) {
	errInjected := fmt.Errorf("injected")

	t.Run("success", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		h.MockDynamo.EXPECT().
			UpdateItem(mock.Anything, mock.Anything).
			Return(nil, nil)

		err := h.ConfirmUser(t.Context(), "user-sub-123", time.Now().UTC())

		require.NoError(t, err)
	})

	t.Run("err-update-item-failed", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		h.MockDynamo.EXPECT().
			UpdateItem(mock.Anything, mock.Anything).
			Return(nil, errInjected)

		err := h.ConfirmUser(t.Context(), "user-sub-123", time.Now().UTC())

		require.ErrorIs(t, err, errInjected)
	})
}

func TestClient_RecordLogin(t *testing.T) {
	errInjected := fmt.Errorf("injected")

	t.Run("success", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		h.MockDynamo.EXPECT().
			UpdateItem(mock.Anything, mock.Anything).
			Return(nil, nil)

		err := h.RecordLogin(t.Context(), "user-sub-123", time.Now().UTC())

		require.NoError(t, err)
	})

	t.Run("err-update-item-failed", func(t *testing.T) {
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		h.MockDynamo.EXPECT().
			UpdateItem(mock.Anything, mock.Anything).
			Return(nil, errInjected)

		err := h.RecordLogin(t.Context(), "user-sub-123", time.Now().UTC())

		require.ErrorIs(t, err, errInjected)
	})
}
