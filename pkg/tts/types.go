package tts

// ProviderType represents the type of TTS provider
type ProviderType string

const (
    ProviderAWS         ProviderType = "aws"
    ProviderGoogle      ProviderType = "google"
    ProviderMicrosoft   ProviderType = "microsoft"
    ProviderIBM         ProviderType = "ibm"
    ProviderElevenLabs  ProviderType = "elevenlabs"
    ProviderESpeak      ProviderType = "espeak"
    ProviderSherpaONNX  ProviderType = "sherpa-onnx"
)

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

// CacheConfig defines caching behavior
type CacheConfig struct {
    Enabled      bool
    Directory    string
    MaxSize      int64  // Maximum cache size in bytes
    TTL          int64  // Time-to-live in seconds
    FilePattern  string // Pattern for cache files
} 