package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ninnemana/pi-server/handlers"
	"github.com/pkg/errors"
)

var (
	certFile = flag.String("cert", "", "Full certificate file path")
	keyFile  = flag.String("key", "", "Private TLS key")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example server.\n"))
	})

	mux.HandleFunc("/webhooks/github", handlers.Github)

	httpRedirect := http.NewServeMux()
	httpRedirect.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		target := "https://" + req.Host + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, target, http.StatusPermanentRedirect)
	})

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	secure := make(chan error)
	insecure := make(chan error)
	go func() {
		tlsSrv := &http.Server{
			Addr:         ":443",
			Handler:      mux,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		switch {
		case certFile == nil || keyFile == nil:
		default:
			_, err := tls.LoadX509KeyPair(*certFile, *keyFile)
			if err != nil {
				log.Printf("%+v\n", errors.Wrap(err, "certificate files were not valid"))
				return
			}

			secure <- errors.Wrap(
				tlsSrv.ListenAndServeTLS(
					*certFile,
					*keyFile,
				),
				"fell out of serving secure web server",
			)
		}
	}()

	go func() {
		srv := &http.Server{
			Addr:    ":80",
			Handler: httpRedirect,
		}

		insecure <- errors.Wrap(
			srv.ListenAndServe(),
			"fell out of serving insecure traffic",
		)
	}()

	select {
	case err := <-secure:
		log.Fatalf("secure server crash: %v", err)
		os.Exit(1)
	case err := <-insecure:
		log.Fatalf("insecure server crash: %v", err)
		os.Exit(1)
	}
}
