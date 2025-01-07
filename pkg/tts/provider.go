package tts

import "fmt"

// ProviderType represents the type of TTS provider
type ProviderType string

const (
	ProviderAWS        ProviderType = "aws"
	ProviderGoogle     ProviderType = "google"
	ProviderMicrosoft  ProviderType = "microsoft"
	ProviderIBM        ProviderType = "ibm"
	ProviderElevenLabs ProviderType = "elevenlabs"
	ProviderWitAI      ProviderType = "witai"
	ProviderESpeak     ProviderType = "espeak"
	ProviderSherpa     ProviderType = "sherpa"
)

// providerConstructor is a function that creates a new provider instance
type providerConstructor func(TTSConfig) (TTSProvider, error)

// providerRegistry stores constructor functions for each provider type
var providerRegistry = make(map[ProviderType]providerConstructor)

// RegisterProvider registers a provider constructor for a given type
func RegisterProvider(pType ProviderType, constructor providerConstructor) {
	providerRegistry[pType] = constructor
}

// NewTTSProvider creates a new TTS provider of the specified type
func NewTTSProvider(providerType ProviderType, config TTSConfig) (TTSProvider, error) {
	constructor, ok := providerRegistry[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
	return constructor(config)
}
