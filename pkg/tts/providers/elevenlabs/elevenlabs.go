package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const elevenLabsBaseURL = "https://api.elevenlabs.io/v1"

// ElevenLabsProvider implements TTSProvider for ElevenLabs
type ElevenLabsProvider struct {
	*BaseProvider
	apiKey      string
	client      *http.Client
	audioPlayer *AudioPlayer
}

// NewElevenLabsProvider creates a new ElevenLabs provider
func NewElevenLabsProvider(cfg TTSConfig) (*ElevenLabsProvider, error) {
	audioPlayer, err := NewAudioPlayer()
	if err != nil {
		return nil, err
	}

	return &ElevenLabsProvider{
		BaseProvider: NewBaseProvider(cfg),
		apiKey:      cfg.APIKey,
		client:      &http.Client{},
		audioPlayer: audioPlayer,
	}, nil
}

type synthesisRequest struct {
	Text     string           `json:"text"`
	ModelID  string           `json:"model_id,omitempty"`
	VoiceSettings voiceSettings `json:"voice_settings,omitempty"`
}

type voiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
}

func (p *ElevenLabsProvider) synthesize(ctx context.Context, text string, isSSML bool) ([]byte, error) {
	if isSSML {
		return nil, fmt.Errorf("SSML not supported by ElevenLabs")
	}

	reqBody := synthesisRequest{
		Text: text,
		VoiceSettings: voiceSettings{
			Stability:       0.75,
			SimilarityBoost: 0.75,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/text-to-speech/%s", elevenLabsBaseURL, p.config.VoiceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *ElevenLabsProvider) Speak(ctx context.Context, text string) error {
	audioData, err := p.synthesize(ctx, text, false)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(bytes.NewReader(audioData))
}

func (p *ElevenLabsProvider) SpeakSSML(ctx context.Context, ssml string) error {
	return fmt.Errorf("SSML not supported by ElevenLabs")
}

func (p *ElevenLabsProvider) GetVoices(ctx context.Context) ([]Voice, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", elevenLabsBaseURL+"/voices", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("xi-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get voices: status %d", resp.StatusCode)
	}

	var result struct {
		Voices []struct {
			VoiceID  string `json:"voice_id"`
			Name     string `json:"name"`
			Category string `json:"category"`
		} `json:"voices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	voices := make([]Voice, 0, len(result.Voices))
	for _, v := range result.Voices {
		voices = append(voices, Voice{
			ID:       v.VoiceID,
			Name:     v.Name,
			Provider: "ElevenLabs",
		})
	}

	return voices, nil
}

// Audio control methods
func (p *ElevenLabsProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *ElevenLabsProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *ElevenLabsProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *ElevenLabsProvider) SetOutputDevice(deviceID string) error {
	return ErrNotImplemented
}

func (p *ElevenLabsProvider) CheckCredentials(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", elevenLabsBaseURL+"/user/subscription", nil)
	if err != nil {
		return false
	}

	req.Header.Set("xi-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
