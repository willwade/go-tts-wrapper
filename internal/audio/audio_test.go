package pkg

import (
	"bytes"
	"testing"
	"time"
)

func TestNewAudioPlayer(t *testing.T) {
	player, err := NewAudioPlayer()
	if err != nil {
		t.Skipf("Skipping test: could not initialize audio: %v", err)
		return
	}
	defer player.Close()

	if player.stream != nil {
		t.Error("Expected initial stream to be nil")
	}
	if player.playing {
		t.Error("Expected initial playing state to be false")
	}
	if player.paused {
		t.Error("Expected initial paused state to be false")
	}
}

func TestAudioPlayerStates(t *testing.T) {
	player, err := NewAudioPlayer()
	if err != nil {
		t.Skipf("Skipping test: could not initialize audio: %v", err)
		return
	}
	defer player.Close()

	t.Run("Initial state", func(t *testing.T) {
		if player.IsPlaying() {
			t.Error("Expected IsPlaying() to be false initially")
		}
	})

	t.Run("State transitions", func(t *testing.T) {
		// Test with empty audio (should error but not crash)
		err := player.PlayMP3Stream(bytes.NewReader([]byte{}))
		if err == nil {
			t.Error("Expected error playing empty stream")
		}

		if err := player.Pause(); err == nil {
			t.Error("Expected error when pausing with no active playback")
		}

		if err := player.Resume(); err == nil {
			t.Error("Expected error when resuming with no active playback")
		}

		if err := player.Stop(); err != nil {
			t.Errorf("Unexpected error when stopping with no active playback: %v", err)
		}
	})
}

func TestWaitForCompletion(t *testing.T) {
	player, err := NewAudioPlayer()
	if err != nil {
		t.Skipf("Skipping test: could not initialize audio: %v", err)
		return
	}
	defer player.Close()

	t.Run("No active playback", func(t *testing.T) {
		if err := player.WaitForCompletion(); err == nil {
			t.Error("Expected error waiting for completion with no active playback")
		}
	})

	t.Run("Stop during wait", func(t *testing.T) {
		// Start a wait operation
		waitDone := make(chan struct{})
		go func() {
			player.WaitForCompletion()
			close(waitDone)
		}()

		// Stop should trigger completion
		time.Sleep(100 * time.Millisecond)
		player.Stop()

		select {
		case <-waitDone:
			// Success
		case <-time.After(time.Second):
			t.Error("WaitForCompletion did not return after Stop")
		}
	})
}

func TestMP3ToPCMConversion(t *testing.T) {
	testCases := []struct {
		name        string
		input       []byte
		expectError bool
	}{
		{
			name:        "Empty input",
			input:       []byte{},
			expectError: true,
		},
		{
			name:        "Invalid MP3 data",
			input:       []byte{0x01, 0x02, 0x03},
			expectError: true,
		},
		{
			name: "Minimal MP3 header",
			input: []byte{
				0xFF, 0xFB, 0x90, 0x64, // MPEG-1 Layer 3 header
				0x00, 0x00, 0x00, 0x00, // Some dummy data
			},
			expectError: true, // Valid header but incomplete MP3
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pcm, rate, err := mp3ToPCM(bytes.NewReader(tc.input))
			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if pcm == nil {
					t.Error("Expected PCM data, got nil")
				}
				if rate <= 0 {
					t.Errorf("Expected positive sample rate, got %f", rate)
				}
			}
		})
	}
}

func TestAudioPlayerResourceCleanup(t *testing.T) {
	player, err := NewAudioPlayer()
	if err != nil {
		t.Skipf("Skipping test: could not initialize audio: %v", err)
		return
	}

	// Test multiple close calls
	if err := player.Close(); err != nil {
		t.Errorf("First close failed: %v", err)
	}

	if err := player.Close(); err != nil {
		t.Errorf("Second close should succeed: %v", err)
	}
}
