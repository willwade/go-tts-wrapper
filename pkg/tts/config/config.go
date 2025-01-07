// pkg/tts/config/config.go
package tts

import (
	"encoding/json"
	"os"

	"github.com/willwade/go-tts-wrapper/pkg/tts"
)

type Config struct {
    Providers map[tts.ProviderType]tts.TTSConfig
    Audio     tts.AudioConfig
    Cache     tts.CacheConfig
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}