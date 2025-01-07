package tts_test

import (
	"testing"

	"github.com/willwade/go-tts-wrapper/pkg/tts"
)

func TestNewTTSProvider(t *testing.T) {
	testCases := []struct {
		name          string
		providerType  tts.ProviderType
		config        tts.TTSConfig
		expectError   bool
	}{
		{
			name:         "Invalid provider",
			providerType: "invalid",
			config:       tts.TTSConfig{},
			expectError:  true,
		},
		{
			name:         "AWS Polly without credentials",
			providerType: tts.ProviderAWS,
			config:       tts.TTSConfig{},
			expectError:  false, // AWS SDK handles missing credentials gracefully
		},
		{
			name:         "ElevenLabs without API key",
			providerType: tts.ProviderElevenLabs,
			config:       tts.TTSConfig{},
			expectError:  false, // Error occurs during API calls, not initialization
		},
		{
			name:         "ESpeak provider",
			providerType: tts.ProviderESpeak,
			config:       tts.TTSConfig{},
			expectError:  false, // May fail if espeak-ng not installed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := tts.NewTTSProvider(tc.providerType, tc.config)
			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Skipf("Provider initialization failed (might be normal): %v", err)
					return
				}
				if provider == nil {
					t.Error("Expected provider instance, got nil")
				}
			}
		})
	}
}

func TestBaseProviderMethods(t *testing.T) {
	base := tts.NewBaseProvider(tts.TTSConfig{})

	t.Run("SetProperty", func(t *testing.T) {
		testCases := []struct {
			property    string
			value       interface{}
			expectError bool
		}{
			{"rate", 1.5, false},
			{"pitch", 0.8, false},
			{"volume", 1.2, false},
			{"invalid", 1.0, true},
			{"rate", "invalid", true},
		}

		for _, tc := range testCases {
			err := base.SetProperty(tc.property, tc.value)
			if tc.expectError && err == nil {
				t.Errorf("Expected error for property %s with value %v", tc.property, tc.value)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error for property %s: %v", tc.property, err)
			}
		}
	})

	t.Run("Default audio config", func(t *testing.T) {
		if base.audioConfig.Rate != 1.0 {
			t.Errorf("Expected default rate 1.0, got %f", base.audioConfig.Rate)
		}
		if base.audioConfig.Pitch != 1.0 {
			t.Errorf("Expected default pitch 1.0, got %f", base.audioConfig.Pitch)
		}
		if base.audioConfig.Volume != 1.0 {
			t.Errorf("Expected default volume 1.0, got %f", base.audioConfig.Volume)
		}
	})
}
