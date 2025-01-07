package ibm

import (
	"context"
	"fmt"
	"io"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/watson-developer-cloud/go-sdk/texttospeechv1"
	tts "github.com/willwade/go-tts-wrapper"
)

type WatsonProvider struct {
	*tts.BaseProvider
	client      *texttospeechv1.TextToSpeechV1
	audioPlayer *tts.AudioPlayer
}

func NewWatsonProvider(cfg tts.TTSConfig) (*WatsonProvider, error) {
	authenticator := &core.IamAuthenticator{
		ApiKey: cfg.APIKey,
	}

	options := &texttospeechv1.TextToSpeechV1Options{
		Authenticator: authenticator,
	}

	if cfg.Region != "" {
		options.URL = fmt.Sprintf("https://api.%s.text-to-speech.watson.cloud.ibm.com", cfg.Region)
	}

	client, err := texttospeechv1.NewTextToSpeechV1(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create Watson client: %w", err)
	}

	audioPlayer, err := tts.NewAudioPlayer()
	if err != nil {
		return nil, err
	}

	return &WatsonProvider{
		BaseProvider: tts.NewBaseProvider(cfg),
		client:      client,
		audioPlayer: audioPlayer,
	}, nil
}

func (p *WatsonProvider) synthesize(ctx context.Context, text string, isSSML bool) (io.ReadCloser, error) {
	var synthesizeOptions *texttospeechv1.SynthesizeOptions

	if isSSML {
		synthesizeOptions = p.client.NewSynthesizeOptions(text).
			SetAccept("audio/mp3").
			SetVoice(p.config.VoiceID).
			SetTextType("ssml")
	} else {
		synthesizeOptions = p.client.NewSynthesizeOptions(text).
			SetAccept("audio/mp3").
			SetVoice(p.config.VoiceID)
	}

	result, _, err := p.client.Synthesize(synthesizeOptions)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ... (rest of implementation similar to previous file)

func init() {
	tts.RegisterProvider(tts.ProviderIBM, func(cfg tts.TTSConfig) (tts.TTSProvider, error) {
		return NewWatsonProvider(cfg)
	})
} 