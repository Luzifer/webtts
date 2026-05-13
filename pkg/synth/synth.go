// Package synth provides an interface for audio synthesis.
package synth

import "context"

type (
	// Provider defines the methods required to generate audio.
	Provider interface {
		// GenerateAudio transforms the given voice, language and text into OGG Audio
		GenerateAudio(ctx context.Context, voice, language, text string) ([]byte, error)
	}
)
