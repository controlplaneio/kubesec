package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/sublimino/kubesec/pkg/ruler"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ListenAndServe starts a web server and waits for SIGTERM
func ListenAndServe(port string, timeout time.Duration, logger *zap.SugaredLogger, stopCh <-chan struct{}) {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "https://kubesec.io", http.StatusSeeOther)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data []byte
		isJson := json.Valid(body)
		if isJson {
			data = body
		} else {
			data, err = yaml.YAMLToJSON(body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		report := ruler.NewRuleset(logger).Run(data)
		res, err := json.Marshal(report)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(PrettyJSON(res)))
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
	return string(out.Bytes())
}
