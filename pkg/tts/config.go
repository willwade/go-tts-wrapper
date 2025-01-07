// Package tts provides a unified interface for text-to-speech providers
package tts

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

// Voice represents a TTS voice with standardized properties
type Voice struct {
	ID          string
	Name        string
	Language    string
	Gender      string
	Provider    string
	NativeVoice any // Original provider-specific voice object
}
