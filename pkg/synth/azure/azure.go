// Package azure provides a text-to-speech synthesis provider for Azure.
package azure

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Luzifer/webtts/pkg/synth"
	"github.com/sirupsen/logrus"
)

type (
	// Provider represents the Azure text-to-speech synthesis provider.
	Provider struct{}

	ssmlRequest struct {
		XMLName xml.Name `xml:"speak"`
		Text    string   `xml:",chardata"`
		Version string   `xml:"version,attr"`
		Lang    string   `xml:"xml:lang,attr"`
		Voice   struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"voice"`
	}
)

var _ synth.Provider = (*Provider)(nil)

// New creates a new instance of the Azure text-to-speech synthesis provider.
func New() (*Provider, error) {
	speechKey := os.Getenv("AZURE_SPEECH_RESOURCE_KEY")
	speechRegion := os.Getenv("AZURE_SPEECH_REGION")

	if speechKey == "" || speechRegion == "" {
		return nil, fmt.Errorf("missing environment variables: AZURE_SPEECH_RESOURCE_KEY and AZURE_SPEECH_REGION")
	}

	return &Provider{}, nil
}

// GenerateAudio generates audio from the given text using Azure's Text-to-Speech service.
func (p Provider) GenerateAudio(ctx context.Context, voice, language, text string) ([]byte, error) {
	speechKey := os.Getenv("AZURE_SPEECH_RESOURCE_KEY")
	speechRegion := os.Getenv("AZURE_SPEECH_REGION")

	body, err := p.requestSSML(voice, language, text)
	if err != nil {
		return nil, fmt.Errorf("generating SSML: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiURL(speechRegion), body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", speechKey)
	req.Header.Set("User-Agent", "webtts/0.x (https://github.com/Luzifer/webtts)")
	req.Header.Set("X-Microsoft-OutputFormat", "ogg-48khz-16bit-mono-opus")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting audio: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.WithError(err).Error("closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading audio data:%w", err)
	}

	return audioData, nil
}

func (Provider) apiURL(region string) string {
	return fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", region)
}

func (Provider) requestSSML(voice, language, text string) (io.Reader, error) {
	var req ssmlRequest
	req.Lang = language
	req.Version = "1.0"
	req.Voice.Name = voice
	req.Voice.Text = text

	data, err := xml.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling ssml request: %w", err)
	}

	return bytes.NewReader(data), nil
}
