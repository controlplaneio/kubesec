package server

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/controlplaneio/kubesec/v2/pkg/report"
	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/dsse"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func parseEd25519PrivateKey(data []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to PEM decode private key")
	}

	switch block.Type {
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		privKey, ok := key.(ed25519.PrivateKey)
		if !ok {
			return nil, errors.New("not an ed25519 private key")
		}
		return privKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}

type HandlerFunc func(c context.Context) error

type HTTPError interface {
	error
	Unwrap() error
	Msg() string
	Status() int
}

type httpError struct {
	err        error
	statusCode int
	msg        string // public message to provide to users
}

var (
	// check at compile time that httpError implements HTTPError correctly
	_ HTTPError = &httpError{}
)

func (he *httpError) Error() string { return he.err.Error() }

func (he *httpError) Unwrap() error { return he.err }

func (he *httpError) Wrap(err error) *httpError {
	if he == nil {
		return nil
	}
	he.err = errors.Join(he.err, err)
	return he
}

func (he *httpError) Msg() string {
	// if a public message is set
	if he.msg != "" {
		return he.msg
	}

	// fallback generic responses based on status
	switch he.statusCode {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}
func (he *httpError) Status() int { return he.statusCode }

// In wrapping error treat all httpError types as nil-able
func WrapHTTPError(err error) *httpError {
	if err == nil {
		return nil
	}
	return &httpError{
		statusCode: http.StatusBadRequest,
		err:        err,
	}
}

// In wrapping error treat all httpError types as nil-able
func NewHTTPError(msg string) *httpError {
	return &httpError{
		statusCode: http.StatusBadRequest,
		msg:        msg,
		err:        errors.New(msg),
	}
}

func (he *httpError) WithStatus(status int) *httpError {
	if he == nil {
		return nil
	}
	he.statusCode = status
	return he
}

func (he *httpError) WithMsg(msg string) *httpError {
	if he == nil {
		return nil
	}
	he.msg = msg
	return he
}

func (he *httpError) Public() *httpError {
	if he == nil {
		return nil
	}
	he.msg = he.err.Error()
	return he
}

func (he *httpError) Private() *httpError {
	if he == nil {
		return nil
	}
	he.msg = ""
	return he
}

// errorHandler allows for writing a handler which can return an error then this
// function can handle it generically conforming to http.HandlerFunc.
func errorHandler(logger *zap.SugaredLogger, f func(http.ResponseWriter, *http.Request) HTTPError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpError := f(w, r)
		if httpError != nil {
			logger.Errorf("%s %+v", httpError.Msg(), httpError.Error())
			http.Error(w, httpError.Msg(), httpError.Status())
		}
	}
}

// ListenAndServe starts a web server and waits for SIGTERM
func ListenAndServe(
	addr string,
	timeout time.Duration,
	logger *zap.SugaredLogger,
	stopCh <-chan struct{},
	keypath string,
	schemaConfig ruler.SchemaConfig,
) {

	mux := http.DefaultServeMux
	mux.Handle("/", scanHandler(logger, keypath, schemaConfig))
	mux.Handle("/scan", scanHandler(logger, keypath, schemaConfig))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", errorHandler(logger, func(w http.ResponseWriter, r *http.Request) HTTPError {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK\n"))
		return WrapHTTPError(err)
	}))

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  15 * time.Second,
	}

	logger.Infof("Starting HTTP server on %s", addr)

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

func retrieveRequestData(r *http.Request) ([]byte, error) {
	// TODO: Implement breaking change respecting header Content-Type
	// contentType := r.Header.Get("Content-Type")
	// err := r.ParseForm()
	// formData := r.Form.Get(formPrefix)

	formPrefix := "file="
	formPrefixLen := len(formPrefix)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("failed reading request body")
	}

	if string(body[:formPrefixLen]) == formPrefix {
		body = body[formPrefixLen:]
	}

	err = r.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func scanHandler(logger *zap.SugaredLogger, keypath string, schemaConfig ruler.SchemaConfig) http.Handler {
	return errorHandler(logger, func(w http.ResponseWriter, r *http.Request) HTTPError {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "https://kubesec.io", http.StatusSeeOther)
			return nil
		}

		// fail early if no in-toto signing key is configured for this server
		if r.URL.Query().Get("in-toto") != "" && keypath == "" {
			logger.Errorf("Attempted to serve an in-toto payload but no key is available")
			w.WriteHeader(http.StatusInternalServerError)
			return NewHTTPError("attempted to serve an in-toto payload but no key is available")
		}

		const fileName = "API"
		body, err := retrieveRequestData(r)
		if err != nil {
			return WrapHTTPError(err).WithStatus(http.StatusBadRequest)
		}

		// Extract rule IDs from query params
		// Supports multiple e.g. ?rule=Foo&rule=Bar, scans all if unset
		var ruleIDs []string
		if rules, ok := r.URL.Query()["rule"]; ok {
			ruleIDs = append(ruleIDs, rules...)
		}

		var payload interface{}
		ruleset, err := ruler.NewRuleset(logger, ruleIDs...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(err.Error() + "\n")); err != nil {
				logger.Errorf("Writing response failed %v", err)
			}
			return nil
		}
		reports, err := ruleset.Run(fileName, body, schemaConfig)
		if err != nil {
			return WrapHTTPError(err).WithStatus(http.StatusBadRequest).Public() // pass through report error
		}

		if r.URL.Query().Get("in-toto") != "" {
			keyData, err := os.ReadFile(keypath)
			if err != nil {
				return NewHTTPError("attempted to serve an in-toto payload but the key is unavailable").Wrap(err)
			}

			privKey, err := parseEd25519PrivateKey(keyData)
			if err != nil {
				return NewHTTPError("failed to parse private key").Wrap(err)
			}

			signer := cryptoutil.NewED25519Signer(privKey)
			keyID, err := signer.KeyID()
			if err != nil {
				return NewHTTPError("failed to compute key ID").Wrap(err)
			}

			linkMb := ruler.GenerateInTotoLink(reports, body)

			payloadBytes, err := json.Marshal(linkMb.Signed)
			if err != nil {
				return NewHTTPError("failed to marshal link payload").Wrap(err)
			}

			envelope, err := dsse.Sign(
				"https://in-toto.io/Statement/v0.1",
				bytes.NewReader(payloadBytes),
				dsse.SignWithSigners(signer),
			)
			if err != nil {
				return NewHTTPError("could not sign in-toto link").Wrap(err)
			}

			for _, sig := range envelope.Signatures {
				linkMb.Signatures = append(linkMb.Signatures, ruler.InTotoSignature{
					KeyID:       keyID,
					Sig:         hex.EncodeToString(sig.Signature),
					Certificate: "",
				})
			}

			payload = map[string]interface{}{
				"reports": reports,
				"link":    linkMb,
			}
		} else {
			payload = reports
		}

		res, err := json.Marshal(payload)
		if err != nil {
			return NewHTTPError("failed to marshal JSON").Wrap(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		formattedOutput, err := report.PrettyJSON(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return NewHTTPError("failed to pretty format the JSON report").Wrap(err)
		}
		_, err = w.Write([]byte(string(formattedOutput) + "\n"))
		if err != nil {
			return WrapHTTPError(err)
		}
		return nil
	})
}
