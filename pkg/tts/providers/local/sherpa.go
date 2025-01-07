package tts

import (
	"bytes"
	"context"
	"fmt"

	"github.com/k2-fsa/sherpa-onnx-go/offline_tts"
)

// SherpaProvider implements TTSProvider for Sherpa-ONNX TTS
type SherpaProvider struct {
	*BaseProvider
	tts         *offline_tts.OfflineTts
	audioPlayer *AudioPlayer
}

// SherpaConfig contains Sherpa-ONNX specific configuration
type SherpaConfig struct {
	VitsModelPath   string  // Path to ONNX model file
	VitsLexiconPath string  // Optional lexicon file path
	SamplingRate    int     // e.g., 22050
	NoiseScale      float32 // Default: 0.667
	NoiseScaleW     float32 // Default: 0.8
	LengthScale     float32 // Default: 1.0
}

// NewSherpaProvider creates a new Sherpa-ONNX TTS provider
func NewSherpaProvider(cfg TTSConfig, sherpaConfig SherpaConfig) (*SherpaProvider, error) {
	config := &offline_tts.OfflineTtsConfig{
		Model: offline_tts.OfflineTtsModelConfig{
			Vits: offline_tts.VitsModelConfig{
				Model:   sherpaConfig.VitsModelPath,
				Lexicon: sherpaConfig.VitsLexiconPath,
			},
		},
		NoiseScale:  sherpaConfig.NoiseScale,
		NoiseScaleW: sherpaConfig.NoiseScaleW,
		LengthScale: sherpaConfig.LengthScale,
		SampleRate:  sherpaConfig.SamplingRate,
	}

	tts, err := offline_tts.NewOfflineTts(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Sherpa-ONNX: %w", err)
	}

	audioPlayer, err := NewAudioPlayer()
	if err != nil {
		tts.Delete()
		return nil, err
	}

	return &SherpaProvider{
		BaseProvider: NewBaseProvider(cfg),
		tts:          tts,
		audioPlayer:  audioPlayer,
	}, nil
}

func (p *SherpaProvider) synthesize(ctx context.Context, text string) ([]float32, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return p.tts.GenerateWithCallback(text, func(progress float32) {
			select {
			case <-ctx.Done():
				// Context cancelled
			default:
				// Could emit progress updates here
			}
		})
	}
}

func (p *SherpaProvider) Speak(ctx context.Context, text string) error {
	samples, err := p.synthesize(ctx, text)
	if err != nil {
		return err
	}

	// Convert float32 samples to MP3
	audioData, err := encodePCMToMP3(samples, p.tts.SampleRate())
	if err != nil {
		return err
	}

	return p.audioPlayer.PlayMP3Stream(bytes.NewReader(audioData))
}

func (p *SherpaProvider) SpeakSSML(ctx context.Context, ssml string) error {
	return fmt.Errorf("SSML not supported by Sherpa-ONNX")
}

func (p *SherpaProvider) GetVoices(ctx context.Context) ([]Voice, error) {
	// Sherpa-ONNX uses model files directly, so we just return the currently loaded model
	voices := []Voice{
		{
			ID:       "default",
			Name:     "Sherpa-ONNX Model",
			Provider: "Sherpa-ONNX",
			Language: "unknown", // Model-dependent
		},
	}
	return voices, nil
}

// Audio control methods
func (p *SherpaProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *SherpaProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *SherpaProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *SherpaProvider) SetOutputDevice(deviceID string) error {
	return ErrNotImplemented
}

func (p *SherpaProvider) CheckCredentials(ctx context.Context) bool {
	// No credentials needed, just check if the TTS engine is initialized
	return p.tts != nil
}

// Close cleans up resources
func (p *SherpaProvider) Close() error {
	if err := p.audioPlayer.Close(); err != nil {
		return err
	}
	p.tts.Delete()
	return nil
}
