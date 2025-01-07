package aws_test

import (
	"context"
	"testing"

	tts "github.com/willwade/go-tts-wrapper"
)

func TestPollyProvider(t *testing.T) {
	cfg := tts.TTSConfig{
		Region:       "us-west-2",
		LanguageCode: "en-US",
		VoiceID:      "Joanna",
	}

	provider, err := NewPollyProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: could not initialize provider: %v", err)
		return
	}

	ctx := context.Background()

	t.Run("GetVoices", func(t *testing.T) {
		voices, err := provider.GetVoices(ctx)
		if err != nil {
			t.Errorf("GetVoices failed: %v", err)
		}
		if len(voices) == 0 {
			t.Error("Expected voices, got empty list")
		}
		for _, voice := range voices {
			if voice.Provider != "AWS Polly" {
				t.Errorf("Expected provider 'AWS Polly', got '%s'", voice.Provider)
			}
		}
	})

	t.Run("Synthesis", func(t *testing.T) {
		tests := []struct {
			name    string
			text    string
			isSSML  bool
			wantErr bool
		}{
			{"plain text", "Hello, world!", false, false},
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
	})

	t.Run("Audio controls", func(t *testing.T) {
		if err := provider.PauseAudio(); err != nil {
			t.Errorf("PauseAudio failed: %v", err)
		}
		if err := provider.ResumeAudio(); err != nil {
			t.Errorf("ResumeAudio failed: %v", err)
		}
		if err := provider.StopAudio(); err != nil {
			t.Errorf("StopAudio failed: %v", err)
		}
	})
}

func TestPollyProviderCredentials(t *testing.T) {
	cfg := tts.TTSConfig{
		Region:  "us-west-2",
		APIKey:  "invalid-key",
		VoiceID: "Joanna",
	}

	provider, err := NewPollyProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: could not initialize provider: %v", err)
		return
	}

	if provider.CheckCredentials(context.Background()) {
		t.Error("Expected CheckCredentials to return false with invalid credentials")
	}
} 