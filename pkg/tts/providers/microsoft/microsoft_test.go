package pkg

import (
	"context"
	"testing"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func TestMicrosoftProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Microsoft Azure TTS tests in short mode")
	}

	cfg := TTSConfig{
		APIKey:  "test-key",
		Region:  "westus",
		VoiceID: "en-US-JennyNeural",
	}

	provider, err := NewMicrosoftProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: could not initialize provider: %v", err)
		return
	}

	t.Run("Provider configuration", func(t *testing.T) {
		tests := []struct {
			name     string
			property string
			value    interface{}
			wantErr  bool
		}{
			{"set rate", "rate", 1.5, false},
			{"set pitch", "pitch", 0.8, false},
			{"set volume", "volume", 1.2, false},
			{"invalid property", "unknown", 1.0, true},
			{"invalid value", "rate", "invalid", true},
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
			{"valid SSML", "<speak version='1.0'><voice name='en-US-JennyNeural'>Hello</voice></speak>", false},
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
		tests := []struct {
			name    string
			fn      func() error
			wantErr bool
		}{
			{"pause", provider.PauseAudio, false},
			{"resume", provider.ResumeAudio, false},
			{"stop", provider.StopAudio, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.fn()
				if (err != nil) != tt.wantErr {
					t.Errorf("%s() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				}
			})
		}
	})
}

func TestMicrosoftProviderSynthesis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping synthesis tests in short mode")
	}

	cfg := TTSConfig{
		APIKey:  "test-key",
		Region:  "westus",
		VoiceID: "en-US-JennyNeural",
	}

	provider, err := NewMicrosoftProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: could not initialize provider: %v", err)
		return
	}

	ctx := context.Background()
	tests := []struct {
		name    string
		text    string
		isSSML  bool
		wantErr bool
	}{
		{"simple text", "Hello, world!", false, false},
		{"empty text", "", false, true},
		{"valid SSML", "<speak version='1.0'><voice name='en-US-JennyNeural'>Hello</voice></speak>", true, false},
		{"invalid SSML", "<speak>Invalid", true, true},
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

func TestMicrosoftProviderCredentials(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		region    string
		wantValid bool
	}{
		{"empty credentials", "", "", false},
		{"invalid key", "invalid-key", "westus", false},
		{"missing region", "test-key", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TTSConfig{
				APIKey: tt.apiKey,
				Region: tt.region,
			}
			provider, err := NewMicrosoftProvider(cfg)
			if err != nil {
				if tt.wantValid {
					t.Errorf("NewMicrosoftProvider() unexpected error = %v", err)
				}
				return
			}
			valid := provider.CheckCredentials(context.Background())
			if valid != tt.wantValid {
				t.Errorf("CheckCredentials() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}
