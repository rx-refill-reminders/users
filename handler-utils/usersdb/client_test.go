package usersdb

import (
	"fmt"
	"testing"

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
		// Create a new ClientTestHarness
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		// Create an example model.User instance
		user := model.User{
			ID: uuid.NewString(),
		}

		// Mock the DynamoDB PutItem call
		h.MockDynamo.EXPECT().
			PutItem(mock.Anything, mock.Anything).
			Return(nil, nil)

		// Call CreateUser
		err := h.CreateUser(t.Context(), user)

		// Expect the call to succeed
		require.NoError(t, err)
	})

	t.Run("err-marshal-failed", func(t *testing.T) {
		// Create a new ClientTestHarness
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		// Create an example model.User instance
		user := model.User{
			ID: uuid.NewString(),
		}

		// Simulate a marshalMap failure
		originalMarshalMap := marshalMap
		t.Cleanup(func() { marshalMap = originalMarshalMap })
		marshalMap = func(_ any) (map[string]types.AttributeValue, error) {
			return nil, errInjected
		}

		// Call CreateUser
		err := h.CreateUser(t.Context(), user)

		// Expect the call to fail, and the result to be errInjected
		require.ErrorIs(t, err, errInjected)
	})

	t.Run("err-put-item-failed", func(t *testing.T) {
		// Create a new ClientTestHarness
		h := NewClientTestHarness(t, Config{
			UsersTable: "users-table",
		})

		// Create an example model.User instance
		user := model.User{
			ID: uuid.NewString(),
		}

		// Mock the DynamoDB PutItem call (errInjected)
		h.MockDynamo.EXPECT().
			PutItem(mock.Anything, mock.Anything).
			Return(nil, errInjected)

		// Call CreateUser
		err := h.CreateUser(t.Context(), user)

		// Expect the call to fail, and the result to be errInjected
		require.ErrorIs(t, err, errInjected)
	})
}
