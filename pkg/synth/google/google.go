// Package google provides a TTS provider for Google's Text-to-Speech service.
package google

import (
	"context"
	"fmt"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"

	"github.com/Luzifer/webtts/pkg/synth"
)

type (
	// Provider represents a TTS provider for Google's Text-to-Speech service.
	Provider struct{}
)

var _ synth.Provider = (*Provider)(nil)

// New creates a new instance of the Google TTS provider.
func New() (*Provider, error) {
	return &Provider{}, nil
}

// GenerateAudio generates audio from the given text using Google's Text-to-Speech service.
func (Provider) GenerateAudio(ctx context.Context, voice, language, text string) ([]byte, error) {
	ttsClient, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating TTS client: %w", err)
	}
	defer ttsClient.Close() //nolint:errcheck

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: language,
			Name:         voice,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	}

	resp, err := ttsClient.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("generating audio: %w", err)
	}

	return resp.AudioContent, nil
}
