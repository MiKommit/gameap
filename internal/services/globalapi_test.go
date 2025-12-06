package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gameap/gameap/internal/config"
	"github.com/gameap/gameap/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalAPIService_Games(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   any
		mockStatusCode int
		wantErr        bool
		errContains    string
		validate       func(t *testing.T, games []domain.GlobalAPIGame)
	}{
		{
			name:           "successful_response_with_full_game_data",
			mockStatusCode: http.StatusOK,
			mockResponse: domain.GlobalAPIResponse[[]domain.GlobalAPIGame]{
				Data: []domain.GlobalAPIGame{
					{
						Code:                  "cstrike",
						StartCode:             "cstrike",
						Name:                  "Counter-Strike 1.6",
						Engine:                "GoldSource",
						EngineVersion:         "1",
						SteamAppIDLinux:       90,
						SteamAppIDWindows:     0,
						RemoteRepositoryLinux: "http://files.gameap.ru/cstrike-1.6/hlcs_base.tar.xz",
						Mods: []domain.GlobalAPIGameMod{
							{
								ID:            3,
								GameCode:      "cstrike",
								Name:          "Classic (Standart)",
								StartCmdLinux: "./hlds_run -game cstrike +ip {ip} +port {port}",
								KickCmd:       "kick #{id}",
								FastRcon: domain.GameModFastRconList{
									{Info: "Status", Command: "status"},
									{Info: "Stats", Command: "stats"},
								},
								Vars: domain.GameModVarList{
									{Var: "default_map", Default: "de_dust2", Info: "Default Map", AdminVar: false},
									{Var: "maxplayers", Default: "32", Info: "Maximum players", AdminVar: false},
								},
							},
						},
					},
				},
				Message: "Games retrieved successfully",
				Success: true,
			},
			wantErr: false,
			validate: func(t *testing.T, games []domain.GlobalAPIGame) {
				t.Helper()

				require.Len(t, games, 1)

				game := games[0]
				assert.Equal(t, "cstrike", game.Code)
				assert.Equal(t, "cstrike", game.StartCode)
				assert.Equal(t, "Counter-Strike 1.6", game.Name)
				assert.Equal(t, "GoldSource", game.Engine)
				assert.Equal(t, "1", game.EngineVersion)
				assert.Equal(t, uint(90), game.SteamAppIDLinux)
				assert.Equal(t, uint(0), game.SteamAppIDWindows)
				assert.Equal(t, "http://files.gameap.ru/cstrike-1.6/hlcs_base.tar.xz", game.RemoteRepositoryLinux)

				require.Len(t, game.Mods, 1)
				mod := game.Mods[0]
				assert.Equal(t, uint(3), mod.ID)
				assert.Equal(t, "cstrike", mod.GameCode)
				assert.Equal(t, "Classic (Standart)", mod.Name)
				assert.Equal(t, "./hlds_run -game cstrike +ip {ip} +port {port}", mod.StartCmdLinux)
				assert.Equal(t, "kick #{id}", mod.KickCmd)

				require.Len(t, mod.FastRcon, 2)
				assert.Equal(t, "Status", mod.FastRcon[0].Info)
				assert.Equal(t, "status", mod.FastRcon[0].Command)

				require.Len(t, mod.Vars, 2)
				assert.Equal(t, "default_map", mod.Vars[0].Var)
				assert.Equal(t, domain.GameModVarDefault("de_dust2"), mod.Vars[0].Default)
			},
		},
		{
			name:           "successful_response_with_multiple_games",
			mockStatusCode: http.StatusOK,
			mockResponse: domain.GlobalAPIResponse[[]domain.GlobalAPIGame]{
				Data: []domain.GlobalAPIGame{
					{
						Code:            "7d2d",
						StartCode:       "7d2d",
						Name:            "7 Days to Die",
						Engine:          "7d2d",
						EngineVersion:   "1.0",
						SteamAppIDLinux: 294420,
						Mods: []domain.GlobalAPIGameMod{
							{
								ID:            11,
								GameCode:      "7d2d",
								Name:          "Default",
								StartCmdLinux: "./startserver.sh -configfile=serverconfig.xml",
							},
						},
					},
					{
						Code:            "arma3",
						StartCode:       "arma3",
						Name:            "Arma 3",
						Engine:          "armedassault3",
						EngineVersion:   "3",
						SteamAppIDLinux: 233780,
						Mods:            []domain.GlobalAPIGameMod{},
					},
					{
						Code:                    "minecraft",
						StartCode:               "minecraft",
						Name:                    "Minecraft",
						Engine:                  "Minecraft",
						EngineVersion:           "1",
						SteamAppIDLinux:         0,
						RemoteRepositoryLinux:   "http://packages.gameap.com/mcrun/mcrun-v1.2.0-linux-amd64.tar.gz",
						RemoteRepositoryWindows: "http://packages.gameap.com/mcrun/mcrun-v1.2.0-windows-amd64.zip",
						Mods: []domain.GlobalAPIGameMod{
							{
								ID:       5,
								GameCode: "minecraft",
								Name:     "Multicore",
								Vars: domain.GameModVarList{
									{Var: "version", Default: "1.20.4", Info: "Minecraft version", AdminVar: false},
									{Var: "memory", Default: "1G", Info: "Memory. Max heap size (Xmx)", AdminVar: false},
								},
							},
						},
					},
				},
				Message: "Games retrieved successfully",
				Success: true,
			},
			wantErr: false,
			validate: func(t *testing.T, games []domain.GlobalAPIGame) {
				t.Helper()

				require.Len(t, games, 3)

				assert.Equal(t, "7d2d", games[0].Code)
				assert.Equal(t, "7 Days to Die", games[0].Name)
				require.Len(t, games[0].Mods, 1)

				assert.Equal(t, "arma3", games[1].Code)
				assert.Equal(t, "Arma 3", games[1].Name)
				assert.Empty(t, games[1].Mods)

				assert.Equal(t, "minecraft", games[2].Code)
				assert.Equal(t, "Minecraft", games[2].Name)
				assert.Equal(t, "http://packages.gameap.com/mcrun/mcrun-v1.2.0-linux-amd64.tar.gz", games[2].RemoteRepositoryLinux)
				assert.Equal(t, "http://packages.gameap.com/mcrun/mcrun-v1.2.0-windows-amd64.zip", games[2].RemoteRepositoryWindows)
				require.Len(t, games[2].Mods, 1)
				require.Len(t, games[2].Mods[0].Vars, 2)
			},
		},
		{
			name:           "successful_response_game_with_empty_mods",
			mockStatusCode: http.StatusOK,
			mockResponse: domain.GlobalAPIResponse[[]domain.GlobalAPIGame]{
				Data: []domain.GlobalAPIGame{
					{
						Code:            "arma2",
						StartCode:       "arma2",
						Name:            "Arma 2",
						Engine:          "armedassault2",
						EngineVersion:   "2",
						SteamAppIDLinux: 33905,
						Mods:            []domain.GlobalAPIGameMod{},
					},
				},
				Message: "Games retrieved successfully",
				Success: true,
			},
			wantErr: false,
			validate: func(t *testing.T, games []domain.GlobalAPIGame) {
				t.Helper()

				require.Len(t, games, 1)
				assert.Equal(t, "arma2", games[0].Code)
				assert.Equal(t, "Arma 2", games[0].Name)
				assert.Empty(t, games[0].Mods)
			},
		},
		{
			name:           "successful_response_empty_games_list",
			mockStatusCode: http.StatusOK,
			mockResponse: domain.GlobalAPIResponse[[]domain.GlobalAPIGame]{
				Data:    []domain.GlobalAPIGame{},
				Message: "No games found",
				Success: true,
			},
			wantErr: false,
			validate: func(t *testing.T, games []domain.GlobalAPIGame) {
				t.Helper()
				assert.Empty(t, games)
			},
		},
		{
			name:           "API_returns_success_false",
			mockStatusCode: http.StatusOK,
			mockResponse: domain.GlobalAPIResponse[[]domain.GlobalAPIGame]{
				Data:    nil,
				Message: "Internal server error",
				Success: false,
			},
			wantErr:     true,
			errContains: "API error",
		},
		{
			name:           "HTTP_error_status_500",
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   nil,
			wantErr:        true,
			errContains:    "unexpected HTTP status code: 500",
		},
		{
			name:           "HTTP_error_status_404",
			mockStatusCode: http.StatusNotFound,
			mockResponse:   nil,
			wantErr:        true,
			errContains:    "unexpected HTTP status code: 404",
		},
		{
			name:           "invalid_JSON_response",
			mockStatusCode: http.StatusOK,
			mockResponse:   "invalid json",
			wantErr:        true,
			errContains:    "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/games", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)

				if tt.mockResponse != nil {
					if str, ok := tt.mockResponse.(string); ok {
						_, _ = w.Write([]byte(str))
					} else {
						_ = json.NewEncoder(w).Encode(tt.mockResponse)
					}
				}
			}))
			defer server.Close()

			// Create service with test server URL
			cfg := &config.Config{}
			cfg.GlobalAPI.URL = server.URL

			service := NewGlobalAPIService(cfg)

			// Execute test
			games, err := service.Games(context.Background())

			// Validate results
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, games)
				}
			}
		})
	}
}

func TestGlobalAPIService_SendBug(t *testing.T) {
	tests := []struct {
		name           string
		report         BugReport
		mockStatusCode int
		wantErr        bool
		errContains    string
		validateReq    func(t *testing.T, r *http.Request)
	}{
		{
			name: "successful bug report",
			report: BugReport{
				Version:     "1.0.0",
				Summary:     "Test bug",
				Description: "This is a test bug",
				Environment: "Test environment\n",
			},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
			validateReq: func(t *testing.T, r *http.Request) {
				t.Helper()

				assert.Equal(t, "/bugs", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var payload map[string]string
				err := json.NewDecoder(r.Body).Decode(&payload)
				require.NoError(t, err)

				assert.Equal(t, "1.0.0", payload["version"])
				assert.Equal(t, "Test bug", payload["summary"])
				assert.Equal(t, "This is a test bug", payload["description"])
				assert.Contains(t, payload["environment"], "Test environment")
				assert.Contains(t, payload["environment"], "Go version:")
				assert.Contains(t, payload["environment"], "OS/Arch:")
			},
		},
		{
			name: "HTTP error status",
			report: BugReport{
				Version:     "1.0.0",
				Summary:     "Test bug",
				Description: "This is a test bug",
			},
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
			errContains:    "unexpected HTTP status code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.validateReq != nil {
					tt.validateReq(t, r)
				}

				w.WriteHeader(tt.mockStatusCode)
			}))
			defer server.Close()

			// Create service with test server URL
			cfg := &config.Config{}
			cfg.GlobalAPI.URL = server.URL

			service := NewGlobalAPIService(cfg)

			// Execute test
			err := service.SendBug(context.Background(), tt.report)

			// Validate results
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGlobalAPIService_Games_ContextCancellation(t *testing.T) {
	// Create a server that delays the response
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.GlobalAPI.URL = server.URL

	service := NewGlobalAPIService(cfg)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Execute test
	_, err := service.Games(ctx)

	// Should fail with context error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}
