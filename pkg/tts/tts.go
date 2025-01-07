package tts

import (
	"context"
	"fmt"
	"io"
)

// Voice represents a TTS voice with standardized properties
type Voice struct {
	ID          string
	Name        string
	Language    string
	Gender      string
	Provider    string
	NativeVoice any // Original provider-specific voice object
}

// TTSConfig holds configuration for TTS providers
type TTSConfig struct {
	APIKey       string
	Region       string
	LanguageCode string
	VoiceID      string
	OutputFormat string
	Engine       string
}

// AudioConfig contains settings for audio output
type AudioConfig struct {
	Rate     float64 // Speech rate (1.0 is normal)
	Pitch    float64 // Voice pitch (1.0 is normal)
	Volume   float64 // Volume level (1.0 is normal)
	DeviceID string  // Output device ID
}

// TTSProvider defines the interface that all TTS providers must implement
type TTSProvider interface {
	// Plain text methods
	Speak(ctx context.Context, text string) error
	SynthToFile(ctx context.Context, text, filename string) error
	SpeakStreamed(ctx context.Context, text string, w io.Writer) error

	// SSML methods
	SpeakSSML(ctx context.Context, ssml string) error
	SynthSSMLToFile(ctx context.Context, ssml, filename string) error
	SpeakSSMLStreamed(ctx context.Context, ssml string, w io.Writer) error

	// Configuration methods
	SetProperty(property string, value interface{}) error
	GetVoices(ctx context.Context) ([]Voice, error)
	Connect(eventName string, callback func(interface{})) error

	// Audio control methods
	PauseAudio() error
	ResumeAudio() error
	StopAudio() error
	SetOutputDevice(deviceID string) error

	// Validation methods
	CheckCredentials(ctx context.Context) bool
	ValidateSSML(ssml string) error
}

// BaseProvider implements common functionality for all providers
type BaseProvider struct {
	config      TTSConfig
	audioConfig AudioConfig
}

// NewBaseProvider creates a new base provider with the given config
func NewBaseProvider(config TTSConfig) *BaseProvider {
	return &BaseProvider{
		config: config,
		audioConfig: AudioConfig{
			Rate:   1.0,
			Pitch:  1.0,
			Volume: 1.0,
		},
	}
}

// SetProperty implements basic property setting
func (b *BaseProvider) SetProperty(property string, value interface{}) error {
	switch property {
	case "rate":
		if rate, ok := value.(float64); ok {
			b.audioConfig.Rate = rate
			return nil
		}
	case "pitch":
		if pitch, ok := value.(float64); ok {
			b.audioConfig.Pitch = pitch
			return nil
		}
	case "volume":
		if volume, ok := value.(float64); ok {
			b.audioConfig.Volume = volume
			return nil
		}
	}
	return fmt.Errorf("invalid property or value type: %s", property)
}

// Common provider errors
var (
	ErrNotImplemented    = fmt.Errorf("method not implemented by provider")
	ErrInvalidProperty   = fmt.Errorf("invalid property")
	ErrInvalidCredential = fmt.Errorf("invalid credentials")
	ErrNoVoicesFound     = fmt.Errorf("no voices found")
	ErrInvalidSSML       = fmt.Errorf("invalid SSML")
)
