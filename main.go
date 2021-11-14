package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	log "github.com/sirupsen/logrus"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		SignatureKey   string `flag:"signature-key" default:"" description:"Key to sign requests with" validate:"nonzero"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	ttsClient *texttospeech.Client

	version = "dev"
)

func init() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("webtts %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	var err error
	if ttsClient, err = texttospeech.NewClient(context.Background()); err != nil {
		log.WithError(err).Fatal("Unable to create TTS client")
	}
	defer ttsClient.Close()

	http.HandleFunc("/tts.ogg", handleTTS)

	var h http.Handler = http.DefaultServeMux
	h = httpHelper.NewHTTPLogHandler(h)

	http.ListenAndServe(cfg.Listen, h)
}

func handleTTS(w http.ResponseWriter, r *http.Request) {
	var (
		text      = r.FormValue("text")
		lang      = r.FormValue("lang")
		signature = r.FormValue("signature")
		validTo   = r.FormValue("valid-to")
		voice     = r.FormValue("voice")
	)

	if text == "" {
		http.Error(w, "no text given", http.StatusBadRequest)
		return
	}

	expiry, err := time.Parse(time.RFC3339, validTo)
	if err != nil || time.Now().After(expiry) {
		http.Error(w, "invalid or expired validity", http.StatusBadRequest)
		return
	}

	if err = checkSignature(signature, r); err != nil {
		log.WithError(err).Error("Signature not validated")
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: defaultVal(lang, "en-US"),
			Name:         defaultVal(voice, "en-US-Wavenet-D"),
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	}

	resp, err := ttsClient.SynthesizeSpeech(r.Context(), &req)
	if err != nil {
		log.WithError(err).Error("Unable to synthesize speech")
		http.Error(w, "unable to synthesize speech", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/ogg")
	w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(resp.AudioContent)
}

func checkSignature(signature string, r *http.Request) error {
	keys := []string{}
	for k := range r.URL.Query() {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	hash := hmac.New(sha256.New, []byte(cfg.SignatureKey))
	for _, k := range keys {
		v := r.URL.Query().Get(k)
		if v == "" {
			continue
		}
		fmt.Fprintf(hash, "%s=%s\n", k, v)
	}

	if signature != fmt.Sprintf("%x", hash.Sum(nil)) {
		return errors.New("signature mismatch")
	}

	return nil
}

func defaultVal(s string, d string) string {
	if s != "" {
		return s
	}
	return d
}
