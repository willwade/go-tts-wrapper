// pkg/tts/registry.go
package tts

import "sync"

type ProviderFactory func(config TTSConfig) (TTSProvider, error)

var (
    providerFactories = make(map[ProviderType]ProviderFactory)
    providersMu sync.RWMutex
)

func RegisterProvider(name ProviderType, factory ProviderFactory) {
    providersMu.Lock()
    defer providersMu.Unlock()
    providerFactories[name] = factory
}