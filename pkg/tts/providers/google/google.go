package pkg

import (
	"context"
	"io"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

// GoogleProvider implements TTSProvider for Google Cloud TTS
type GoogleProvider struct {
	*BaseProvider
	client      *texttospeech.Client
	audioPlayer *AudioPlayer
}

// NewGoogleProvider creates a new Google Cloud TTS provider
func NewGoogleProvider(cfg TTSConfig) (*GoogleProvider, error) {
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	audioPlayer, err := NewAudioPlayer()
	if err != nil {
		client.Close()
		return nil, err
	}

	return &GoogleProvider{
		BaseProvider: NewBaseProvider(cfg),
		client:      client,
		audioPlayer: audioPlayer,
	}, nil
}

func (p *GoogleProvider) synthesize(ctx context.Context, text string, isSSML bool) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: p.config.LanguageCode,
			Name:         p.config.VoiceID,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate: p.audioConfig.Rate,
			Pitch:        p.audioConfig.Pitch,
			VolumeGainDb: 20 * p.audioConfig.Volume, // Convert to dB scale
		},
	}

	if isSSML {
		req.Input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: text},
		}
	} else {
		req.Input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		}
	}

	return p.client.SynthesizeSpeech(ctx, req)
}

func (p *GoogleProvider) Speak(ctx context.Context, text string) error {
	resp, err := p.synthesize(ctx, text, false)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(io.NopCloser(io.NewSectionReader(resp.AudioContent, 0, int64(len(resp.AudioContent)))))
}

func (p *GoogleProvider) SpeakSSML(ctx context.Context, ssml string) error {
	if err := p.ValidateSSML(ssml); err != nil {
		return err
	}
	resp, err := p.synthesize(ctx, ssml, true)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(io.NopCloser(io.NewSectionReader(resp.AudioContent, 0, int64(len(resp.AudioContent)))))
}

func (p *GoogleProvider) GetVoices(ctx context.Context) ([]Voice, error) {
	resp, err := p.client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{})
	if err != nil {
		return nil, err
	}

	voices := make([]Voice, 0, len(resp.Voices))
	for _, v := range resp.Voices {
		if len(v.LanguageCodes) == 0 {
			continue
		}
		voices = append(voices, Voice{
			ID:          v.Name,
			Name:        v.Name,
			Language:    v.LanguageCodes[0],
			Gender:      v.SsmlGender.String(),
			Provider:    "Google",
			NativeVoice: v,
		})
	}
	return voices, nil
}

// Audio control methods
func (p *GoogleProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *GoogleProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *GoogleProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *GoogleProvider) SetOutputDevice(deviceID string) error {
	// Not implemented for Google Cloud TTS
	return ErrNotImplemented
}

func (p *GoogleProvider) CheckCredentials(ctx context.Context) bool {
	_, err := p.client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{})
	return err == nil
}
