package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple_password",
			password: "password123",
		},
		{
			name:     "complex_password",
			password: "P@ssw0rd!#$%^&*()",
		},
		{
			name:     "empty_password",
			password: "",
		},
		{
			name:     "long_password",
			password: strings.Repeat("a", 72),
		},
		{
			name:     "unicode_password",
			password: "пароль密码🔐",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			require.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash)

			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
			assert.NoError(t, err)
		})
	}
}

func TestHashPassword_UniqueHashes(t *testing.T) {
	password := "samepassword"

	hash1, err := HashPassword(password)
	require.NoError(t, err)

	hash2, err := HashPassword(password)
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "hashes should be unique due to different salts")
}

func TestHashPassword_ExceedsBcryptLimit(t *testing.T) {
	password := strings.Repeat("a", 100)

	_, err := HashPassword(password)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
}

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		wrongPassword  string
		expectErrOnBad bool
	}{
		{
			name:           "correct_password",
			password:       "correctpassword",
			wrongPassword:  "wrongpassword",
			expectErrOnBad: true,
		},
		{
			name:           "empty_password",
			password:       "",
			wrongPassword:  "notempty",
			expectErrOnBad: true,
		},
		{
			name:           "complex_password",
			password:       "C0mpl3x!@#$%",
			wrongPassword:  "simple",
			expectErrOnBad: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			require.NoError(t, err)

			err = VerifyPassword(hash, tt.password)
			assert.NoError(t, err, "correct password should verify successfully")

			err = VerifyPassword(hash, tt.wrongPassword)
			assert.Error(t, err, "wrong password should not verify")
			assert.Contains(t, err.Error(), "password verification failed")
		})
	}
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	tests := []struct {
		name       string
		hashValue  string
		password   string
		wantErrMsg string
	}{
		{
			name:       "invalid_hash_format",
			hashValue:  "notavalidhash",
			password:   "anypassword",
			wantErrMsg: "password verification failed",
		},
		{
			name:       "empty_hash",
			hashValue:  "",
			password:   "anypassword",
			wantErrMsg: "password verification failed",
		},
		{
			name:       "truncated_hash",
			hashValue:  "$2a$10$truncated",
			password:   "anypassword",
			wantErrMsg: "password verification failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hashValue, tt.password)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrMsg)
		})
	}
}

func TestVerifyPassword_CaseSensitive(t *testing.T) {
	password := "CaseSensitive"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	err = VerifyPassword(hash, "casesensitive")
	assert.Error(t, err, "password verification should be case-sensitive")

	err = VerifyPassword(hash, "CASESENSITIVE")
	assert.Error(t, err, "password verification should be case-sensitive")

	err = VerifyPassword(hash, "CaseSensitive")
	assert.NoError(t, err, "exact password should verify")
}

func TestDefaultBcryptCost(t *testing.T) {
	assert.Equal(t, bcrypt.DefaultCost, DefaultBcryptCost)
}
