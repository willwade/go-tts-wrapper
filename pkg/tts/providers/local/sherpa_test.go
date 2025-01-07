package pkg

import (
	"context"
	"testing"
)

func TestSherpaProvider(t *testing.T) {
	// Skip if model file not available
	modelPath := "path/to/model.onnx"
	config := TTSConfig{
		Engine: modelPath, // Model path
	}
	sherpaConfig := SherpaConfig{
		VitsModelPath: modelPath,
		SamplingRate: 22050,
		NoiseScale:   0.667,
		NoiseScaleW:  0.8,
		LengthScale:  1.0,
	}

	provider, err := NewSherpaProvider(config, sherpaConfig)
	if err != nil {
		t.Skipf("Skipping test: could not initialize Sherpa-ONNX: %v", err)
		return
	}
	defer provider.Close()

	t.Run("GetVoices", func(t *testing.T) {
		voices, err := provider.GetVoices(context.Background())
		if err != nil {
			t.Errorf("GetVoices failed: %v", err)
		}
		if len(voices) != 1 {
			t.Errorf("Expected 1 voice, got %d", len(voices))
		}
		if voices[0].Provider != "Sherpa-ONNX" {
			t.Errorf("Expected provider 'Sherpa-ONNX', got '%s'", voices[0].Provider)
		}
	})

	t.Run("SSML not supported", func(t *testing.T) {
		err := provider.SpeakSSML(context.Background(), "<speak>Hello</speak>")
		if err == nil {
			t.Error("Expected error for SSML, got nil")
		}
	})

	t.Run("CheckCredentials", func(t *testing.T) {
		if !provider.CheckCredentials(context.Background()) {
			t.Error("CheckCredentials returned false, expected true")
		}
	})
}

func TestSherpaProviderErrors(t *testing.T) {
	config := TTSConfig{
		Engine: "nonexistent.onnx",
	}
	sherpaConfig := SherpaConfig{
		VitsModelPath: "nonexistent.onnx",
	}

	_, err := NewSherpaProvider(config, sherpaConfig)
	if err == nil {
		t.Error("Expected error with invalid model path, got nil")
	}
}
