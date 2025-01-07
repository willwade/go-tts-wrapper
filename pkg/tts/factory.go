package tts

import "fmt"

// NewTTSProvider creates a new TTS provider of the specified type
func NewTTSProvider(providerType ProviderType, config TTSConfig) (TTSProvider, error) {
    factory, exists := providerFactories[providerType]
    if !exists {
        return nil, fmt.Errorf("unsupported provider type: %s", providerType)
    }
    
    return factory(config)
} 