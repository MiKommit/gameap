package auth

import (
	"testing"
	"time"

	"github.com/gameap/gameap/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPASETOService(t *testing.T) {
	tests := []struct {
		name      string
		secretKey []byte
		wantErr   bool
	}{
		{
			name:      "exact_32_bytes_key",
			secretKey: []byte("12345678901234567890123456789012"),
			wantErr:   false,
		},
		{
			name:      "short_key_padded",
			secretKey: []byte("shortkey"),
			wantErr:   false,
		},
		{
			name:      "long_key_trimmed",
			secretKey: []byte("12345678901234567890123456789012345678901234567890"),
			wantErr:   false,
		},
		{
			name:      "empty_key_padded",
			secretKey: []byte{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewPASETOService(tt.secretKey)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, service)
			assert.NotNil(t, service.parser)
		})
	}
}

func TestPASETOService_GenerateTokenForUser(t *testing.T) {
	service, err := NewPASETOService([]byte("12345678901234567890123456789012"))
	require.NoError(t, err)

	tests := []struct {
		name          string
		user          *domain.User
		tokenDuration time.Duration
	}{
		{
			name: "generate_token_with_1_hour_duration",
			user: &domain.User{
				ID:    1,
				Login: "testuser",
				Email: "test@example.com",
			},
			tokenDuration: time.Hour,
		},
		{
			name: "generate_token_with_24_hours_duration",
			user: &domain.User{
				ID:    2,
				Login: "anotheruser",
				Email: "another@example.com",
			},
			tokenDuration: 24 * time.Hour,
		},
		{
			name: "generate_token_with_empty_login",
			user: &domain.User{
				ID:    3,
				Login: "",
				Email: "empty@example.com",
			},
			tokenDuration: time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateTokenForUser(tt.user, tt.tokenDuration)

			require.NoError(t, err)
			assert.NotEmpty(t, token)
		})
	}
}

func TestPASETOService_ValidateToken(t *testing.T) {
	service, err := NewPASETOService([]byte("12345678901234567890123456789012"))
	require.NoError(t, err)

	user := &domain.User{
		ID:    1,
		Login: "testuser",
		Email: "test@example.com",
	}

	t.Run("valid_token", func(t *testing.T) {
		token, err := service.GenerateTokenForUser(user, time.Hour)
		require.NoError(t, err)

		claims, err := service.ValidateToken(token)

		require.NoError(t, err)
		assert.NotNil(t, claims)

		subject, err := claims.GetSubject()
		require.NoError(t, err)
		assert.Equal(t, "user:login:testuser", subject)
	})

	t.Run("invalid_token_format", func(t *testing.T) {
		claims, err := service.ValidateToken("invalid-token")

		require.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})

	t.Run("token_from_different_key", func(t *testing.T) {
		differentService, err := NewPASETOService([]byte("different_key_1234567890123456"))
		require.NoError(t, err)

		token, err := differentService.GenerateTokenForUser(user, time.Hour)
		require.NoError(t, err)

		claims, err := service.ValidateToken(token)

		require.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestPASETOService_TokenExpiration(t *testing.T) {
	service, err := NewPASETOService([]byte("12345678901234567890123456789012"))
	require.NoError(t, err)

	user := &domain.User{
		ID:    1,
		Login: "testuser",
		Email: "test@example.com",
	}

	t.Run("token_with_negative_duration_expired", func(t *testing.T) {
		token, err := service.GenerateTokenForUser(user, -time.Hour)
		require.NoError(t, err)

		claims, err := service.ValidateToken(token)

		require.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestPASETOService_GenerateUniqueTokens(t *testing.T) {
	service, err := NewPASETOService([]byte("12345678901234567890123456789012"))
	require.NoError(t, err)

	user := &domain.User{
		ID:    1,
		Login: "testuser",
		Email: "test@example.com",
	}

	token1, err := service.GenerateTokenForUser(user, time.Hour)
	require.NoError(t, err)

	token2, err := service.GenerateTokenForUser(user, time.Hour)
	require.NoError(t, err)

	assert.NotEqual(t, token1, token2, "tokens should be unique due to different JTI")
}

func TestPASETOService_SubjectFormat(t *testing.T) {
	service, err := NewPASETOService([]byte("12345678901234567890123456789012"))
	require.NoError(t, err)

	tests := []struct {
		name            string
		login           string
		expectedSubject string
	}{
		{
			name:            "simple_login",
			login:           "admin",
			expectedSubject: "user:login:admin",
		},
		{
			name:            "email_as_login",
			login:           "user@example.com",
			expectedSubject: "user:login:user@example.com",
		},
		{
			name:            "login_with_special_chars",
			login:           "user-name_123",
			expectedSubject: "user:login:user-name_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &domain.User{
				ID:    1,
				Login: tt.login,
			}

			token, err := service.GenerateTokenForUser(user, time.Hour)
			require.NoError(t, err)

			claims, err := service.ValidateToken(token)
			require.NoError(t, err)

			subject, err := claims.GetSubject()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedSubject, subject)
		})
	}
}
