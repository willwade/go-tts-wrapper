package tts

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWitAIProvider(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/synthesize":
			// Return mock MP3 data
			w.Header().Set("Content-Type", "audio/mpeg")
			w.Write([]byte("mock mp3 data"))
		case "/synthesize/voices":
			// Return mock voices data
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[
				{
					"id": "voice1",
					"name": "Test Voice",
					"locale": "en-US",
					"gender": "female",
					"style": "natural",
					"category": "standard"
				}
			]`))
		}
	}))
	defer server.Close()

	// Create provider with mock server URL
	witAIBaseURL = server.URL + "/synthesize"
	provider, err := NewWitAIProvider(TTSConfig{
		APIKey:  "test-key",
		VoiceID: "voice1",
	})
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	t.Run("GetVoices", func(t *testing.T) {
		voices, err := provider.GetVoices(context.Background())
		if err != nil {
			t.Errorf("GetVoices failed: %v", err)
		}
		if len(voices) != 1 {
			t.Errorf("Expected 1 voice, got %d", len(voices))
		}
		if voices[0].ID != "voice1" {
			t.Errorf("Expected voice ID 'voice1', got '%s'", voices[0].ID)
		}
	})

	t.Run("SSML not supported", func(t *testing.T) {
		err := provider.SpeakSSML(context.Background(), "<speak>Hello</speak>")
		if err == nil {
			t.Error("Expected error for SSML, got nil")
		}
	})

	t.Run("CheckCredentials", func(t *testing.T) {
		if !provider.CheckCredentials(context.Background()) {
			t.Error("CheckCredentials returned false, expected true")
		}
	})
}

func TestWitAIProviderErrors(t *testing.T) {
	// Mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	witAIBaseURL = server.URL + "/synthesize"
	provider, err := NewWitAIProvider(TTSConfig{
		APIKey:  "invalid-key",
		VoiceID: "voice1",
	})
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	t.Run("Invalid credentials", func(t *testing.T) {
		if provider.CheckCredentials(context.Background()) {
			t.Error("CheckCredentials returned true with invalid key")
		}
	})

	t.Run("Synthesis failure", func(t *testing.T) {
		err := provider.Speak(context.Background(), "Hello")
		if err == nil {
			t.Error("Expected error for synthesis with invalid key, got nil")
		}
	})
}
