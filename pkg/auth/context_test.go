package auth_test

import (
	"context"
	"testing"

	"github.com/gameap/gameap/internal/domain"
	"github.com/gameap/gameap/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionFromContext(t *testing.T) {
	t.Run("with_session_in_context", func(t *testing.T) {
		testUser := &domain.User{
			ID:    1,
			Login: "testuser",
			Email: "test@example.com",
		}

		ctx := auth.ContextWithSession(context.Background(), &auth.Session{
			Login: testUser.Login,
			Email: testUser.Email,
		})

		session := auth.SessionFromContext(ctx)

		require.NotNil(t, session)
		assert.Equal(t, testUser.Login, session.Login)
		assert.Equal(t, testUser.Email, session.Email)
	})

	t.Run("without_session_in_context", func(t *testing.T) {
		ctx := context.Background()

		session := auth.SessionFromContext(ctx)

		assert.Nil(t, session)
	})

	t.Run("with_wrong_type_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.SessionKey{}, "not a session")

		session := auth.SessionFromContext(ctx)

		assert.Nil(t, session)
	})
}

func TestContextWithSession(t *testing.T) {
	t.Run("stores_session_in_context", func(t *testing.T) {
		session := &auth.Session{
			ID:    "session-123",
			Login: "admin",
			Email: "admin@example.com",
			User: &domain.User{
				ID:    1,
				Login: "admin",
			},
		}

		ctx := auth.ContextWithSession(context.Background(), session)

		retrieved := auth.SessionFromContext(ctx)
		require.NotNil(t, retrieved)
		assert.Equal(t, session.ID, retrieved.ID)
		assert.Equal(t, session.Login, retrieved.Login)
		assert.Equal(t, session.Email, retrieved.Email)
		assert.Equal(t, session.User, retrieved.User)
	})

	t.Run("stores_nil_session", func(t *testing.T) {
		ctx := auth.ContextWithSession(context.Background(), nil)

		session := auth.SessionFromContext(ctx)

		assert.Nil(t, session)
	})
}

func TestSession_IsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		session  *auth.Session
		expected bool
	}{
		{
			name:     "nil_session",
			session:  nil,
			expected: false,
		},
		{
			name:     "session_with_nil_user",
			session:  &auth.Session{Login: "test"},
			expected: false,
		},
		{
			name:     "session_with_zero_id_user",
			session:  &auth.Session{User: &domain.User{ID: 0}},
			expected: false,
		},
		{
			name:     "authenticated_session",
			session:  &auth.Session{User: &domain.User{ID: 1, Login: "admin"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsAuthenticated()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSession_IsTokenSession(t *testing.T) {
	tests := []struct {
		name     string
		session  *auth.Session
		expected bool
	}{
		{
			name:     "nil_session",
			session:  nil,
			expected: false,
		},
		{
			name:     "session_without_token",
			session:  &auth.Session{Login: "test"},
			expected: false,
		},
		{
			name:     "session_with_zero_id_token",
			session:  &auth.Session{Token: &domain.PersonalAccessToken{ID: 0}},
			expected: false,
		},
		{
			name:     "token_session",
			session:  &auth.Session{Token: &domain.PersonalAccessToken{ID: 1, Name: "api-token"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsTokenSession()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDaemonSessionFromContext(t *testing.T) {
	t.Run("with_daemon_session_in_context", func(t *testing.T) {
		node := &domain.Node{
			ID:   1,
			Name: "test-node",
		}
		daemonSession := &auth.DaemonSession{Node: node}

		ctx := auth.ContextWithDaemonSession(context.Background(), daemonSession)

		retrieved := auth.DaemonSessionFromContext(ctx)

		require.NotNil(t, retrieved)
		assert.Equal(t, node, retrieved.Node)
	})

	t.Run("without_daemon_session_in_context", func(t *testing.T) {
		ctx := context.Background()

		session := auth.DaemonSessionFromContext(ctx)

		assert.Nil(t, session)
	})

	t.Run("with_wrong_type_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.DaemonSessionKey{}, "not a daemon session")

		session := auth.DaemonSessionFromContext(ctx)

		assert.Nil(t, session)
	})
}

func TestContextWithDaemonSession(t *testing.T) {
	t.Run("stores_daemon_session_in_context", func(t *testing.T) {
		node := &domain.Node{
			ID:   5,
			Name: "production-node",
		}
		daemonSession := &auth.DaemonSession{Node: node}

		ctx := auth.ContextWithDaemonSession(context.Background(), daemonSession)

		retrieved := auth.DaemonSessionFromContext(ctx)
		require.NotNil(t, retrieved)
		assert.Equal(t, node.ID, retrieved.Node.ID)
		assert.Equal(t, node.Name, retrieved.Node.Name)
	})

	t.Run("stores_nil_daemon_session", func(t *testing.T) {
		ctx := auth.ContextWithDaemonSession(context.Background(), nil)

		session := auth.DaemonSessionFromContext(ctx)

		assert.Nil(t, session)
	})
}

func TestSessionAndDaemonSessionCoexist(t *testing.T) {
	userSession := &auth.Session{
		ID:    "user-session",
		Login: "admin",
		User:  &domain.User{ID: 1},
	}
	daemonSession := &auth.DaemonSession{
		Node: &domain.Node{ID: 2, Name: "node-1"},
	}

	ctx := context.Background()
	ctx = auth.ContextWithSession(ctx, userSession)
	ctx = auth.ContextWithDaemonSession(ctx, daemonSession)

	retrievedUser := auth.SessionFromContext(ctx)
	retrievedDaemon := auth.DaemonSessionFromContext(ctx)

	require.NotNil(t, retrievedUser)
	require.NotNil(t, retrievedDaemon)
	assert.Equal(t, userSession.ID, retrievedUser.ID)
	assert.Equal(t, daemonSession.Node.ID, retrievedDaemon.Node.ID)
}
