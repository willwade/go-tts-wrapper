// pkg/tts/audio/devices.go
package audio

type AudioDevice struct {
    ID          string
    Name        string
    IsDefault   bool
    SampleRates []int
}

func ListAudioDevices() ([]AudioDevice, error) {
    // Implementation needed
}

func SetDefaultDevice(deviceID string) error {
    // Implementation needed
}