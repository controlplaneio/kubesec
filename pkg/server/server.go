package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/controlplaneio/kubesec/pkg/ruler"
	"github.com/in-toto/in-toto-golang/in_toto"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ListenAndServe starts a web server and waits for SIGTERM
func ListenAndServe(port string, timeout time.Duration, logger *zap.SugaredLogger, stopCh <-chan struct{}, keypath string) {
	mux := http.DefaultServeMux
	mux.Handle("/", scanHandler(logger, keypath))
	mux.Handle("/scan", scanHandler(logger, keypath))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  15 * time.Second,
	}

	logger.Infof("Starting HTTP server on port %s", port)

	// run server in background
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("HTTP server crashed %v", err)
		}
	}()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server graceful shutdown failed %v", err)
	} else {
		logger.Info("HTTP server stopped")
	}
}

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}

func PrettyJSON(b []byte) string {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return out.String()
}

func retrieveRequestData(r *http.Request) ([]byte, error) {
	// TODO: Implement breaking change respecting header Content-Type
	// contentType := r.Header.Get("Content-Type")
	// err := r.ParseForm()
	// formData := r.Form.Get(formPrefix)

	formPrefix := "file="
	formPrefixLen := len(formPrefix)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("Error reading request body")
	}
	defer r.Body.Close()

	if string(body[:formPrefixLen]) == formPrefix {
		body = body[formPrefixLen:]
	}

	return body, nil
}

func scanHandler(logger *zap.SugaredLogger, keypath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "https://kubesec.io", http.StatusSeeOther)
			return
		}

		// fail early if no in-toto signing key is configured for this server
		if r.URL.Query().Get("in-toto") != "" && keypath == "" {
			logger.Errorf("Attempted to serve an in-toto payload but no key is available")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := retrieveRequestData(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error() + "\n"))
			return
		}

		var payload interface{}
		reports, err := ruler.NewRuleset(logger).Run(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error() + "\n"))
			return
		}

		if r.URL.Query().Get("in-toto") != "" {
			json_key, err := ioutil.ReadFile(keypath)
			if err != nil {
				logger.Errorf("Attempted to serve an in-toto payload but the key is unavailable: %v",
					err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			key, err := in_toto.ParseEd25519FromPrivateJSON(string(json_key))
			if err != nil {
				logger.Errorf("Attempted to serve an in-toto payload but the key is unavailable: %v",
					err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			link := ruler.GenerateInTotoLink(reports, body)
			err = link.Sign(key)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			payload = map[string]interface{}{
				"reports": reports,
				"link":    link,
			}
		} else {
			payload = reports
		}

		res, err := json.Marshal(payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(PrettyJSON(res) + "\n"))
	})
}
