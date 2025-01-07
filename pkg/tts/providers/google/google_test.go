package google_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

func TestGoogleProvider(t *testing.T) {
	// Skip in short mode or without credentials
	if testing.Short() {
		t.Skip("Skipping Google Cloud TTS tests in short mode")
	}

	// Mock server for credentials
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token": "mock-token"}`))
	}))
	defer server.Close()

	// Create test client with mock credentials
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx,
		option.WithEndpoint(server.URL),
		option.WithTokenSource(mockTokenSource{}),
	)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer client.Close()

	provider := &GoogleProvider{
		BaseProvider: NewBaseProvider(TTSConfig{
			LanguageCode: "en-US",
			VoiceID:      "en-US-Standard-A",
		}),
		client: client,
	}

	t.Run("Voice configuration", func(t *testing.T) {
		tests := []struct {
			name     string
			property string
			value    interface{}
			wantErr  bool
		}{
			{"valid rate", "rate", 1.5, false},
			{"valid pitch", "pitch", 0.8, false},
			{"valid volume", "volume", 1.2, false},
			{"invalid property", "invalid", 1.0, true},
			{"invalid value type", "rate", "invalid", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := provider.SetProperty(tt.property, tt.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetProperty() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	t.Run("SSML validation", func(t *testing.T) {
		tests := []struct {
			name    string
			ssml    string
			wantErr bool
		}{
			{"valid SSML", "<speak>Hello</speak>", false},
			{"invalid SSML", "Hello", true},
			{"empty SSML", "", true},
			{"malformed SSML", "<speak>Hello", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := provider.ValidateSSML(tt.ssml)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateSSML() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	t.Run("Audio controls", func(t *testing.T) {
		if err := provider.PauseAudio(); err != nil {
			t.Errorf("PauseAudio() error = %v", err)
		}
		if err := provider.ResumeAudio(); err != nil {
			t.Errorf("ResumeAudio() error = %v", err)
		}
		if err := provider.StopAudio(); err != nil {
			t.Errorf("StopAudio() error = %v", err)
		}
	})
}

// Mock token source for testing
type mockTokenSource struct{}

func (m mockTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "mock-token",
	}, nil
}

func TestGoogleProviderSynthesis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping synthesis tests in short mode")
	}

	ctx := context.Background()
	provider, err := NewGoogleProvider(TTSConfig{
		LanguageCode: "en-US",
		VoiceID:      "en-US-Standard-A",
	})
	if err != nil {
		t.Skipf("Skipping test: could not initialize provider: %v", err)
		return
	}

	tests := []struct {
		name    string
		text    string
		isSSML  bool
		wantErr bool
	}{
		{"simple text", "Hello, world!", false, false},
		{"empty text", "", false, true},
		{"valid SSML", "<speak>Hello, world!</speak>", true, false},
		{"invalid SSML", "<speak>Hello", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.isSSML {
				err = provider.SpeakSSML(ctx, tt.text)
			} else {
				err = provider.Speak(ctx, tt.text)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Synthesis error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
