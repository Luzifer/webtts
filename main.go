package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/sirupsen/logrus"

	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
	"github.com/Luzifer/webtts/pkg/synth"
	"github.com/Luzifer/webtts/pkg/synth/google"
)

var (
	cfg = struct {
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		SignatureKey   string `flag:"signature-key" default:"" description:"Key to sign requests with" validate:"nonzero"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() (err error) {
	rconfig.AutoEnv(true)
	if err = rconfig.ParseAndValidate(&cfg); err != nil {
		return fmt.Errorf("parsing CLI options: %w", err)
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing application")
	}

	if cfg.VersionAndExit {
		fmt.Printf("webtts %s\n", version) //nolint:forbidigo
		os.Exit(0)
	}

	http.HandleFunc("/tts.ogg", handleTTS)

	var h http.Handler = http.DefaultServeMux
	h = httpHelper.NewHTTPLogHandler(h)

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           h,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithField("addr", cfg.Listen).Info("starting server")
	if err = server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("starting server")
	}
}

func handleTTS(w http.ResponseWriter, r *http.Request) {
	var (
		text      = r.FormValue("text")
		lang      = r.FormValue("lang")
		provider  = r.FormValue("provider")
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
		logrus.WithError(err).Error("Signature not validated")
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}

	var p synth.Provider
	switch provider {
	case "google", "gcp":
		if p, err = google.New(); err != nil {
			logrus.WithError(err).Error("creating google provider")
			http.Error(w, "creating provider", http.StatusInternalServerError)
			return
		}

	default:
		logrus.WithField("provider", provider).Error("Invalid provider")
		http.Error(w, "invalid provider", http.StatusBadRequest)
		return
	}

	audio, err := p.GenerateAudio(r.Context(), defaultVal(voice, "en-US-Wavenet-D"), defaultVal(lang, "en-US"), text)
	if err != nil {
		logrus.WithError(err).Error("generating audio")
		http.Error(w, "audio generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/ogg")
	w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, _ = w.Write(audio)
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
		fmt.Fprintf(hash, "%s=%s\n", k, v) //nolint:errcheck
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
