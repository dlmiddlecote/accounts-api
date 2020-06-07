package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/blendle/zapdriver"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	account "github.com/dlmiddlecote/accounts-api"
	"github.com/dlmiddlecote/accounts-api/pkg/endpoints"
	"github.com/dlmiddlecote/accounts-api/pkg/service"
	"github.com/dlmiddlecote/kit/api"
)

const (
	buildVersion = "dev"
	namespace    = "ACCOUNTS"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error : %s", err)
		os.Exit(1)
	}
}

func run() error {
	// Configuration
	var cfg struct {
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:8080"`
			DebugHost       string        `conf:"default:0.0.0.0:8090"`
			EnableDebug     bool          `conf:"default:true"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:password,noprint"`
			Host       string `conf:"default:0.0.0.0"`
			Name       string `conf:"default:ometria_core"`
			DisableTLS bool   `conf:"default:false"`
		}
	}

	if err := conf.Parse(os.Args[1:], namespace, &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage(namespace, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}

			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// Logging
	var logger *zap.SugaredLogger
	{
		if l, err := zapdriver.NewProduction(); err != nil {
			return errors.Wrap(err, "creating logger")
		} else {
			logger = l.Sugar()
		}
	}
	defer logger.Sync()

	logger.Infow("Application starting", "version", buildVersion)
	defer logger.Info("Application finished")

	// Debug listener
	if cfg.Web.EnableDebug {

		http.Handle("/metrics", promhttp.Handler())
		go func() {
			logger.Infow("Debug listener starting", "addr", cfg.Web.DebugHost)
			err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux)
			logger.Infow("Debug listener closed", "err", err)
		}()
	}

	// main application setup
	var svc account.Service
	{
		svc = service.NewService(logger.Named("service"))
	}

	var e api.Endpointer
	{
		e = endpoints.NewAccountEndpoints(logger.Named("endpoints"), svc)
	}

	var srv http.Handler
	{
		srv = api.NewServer(logger.Named("server"), e)
	}

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Create our http server
	app := http.Server{
		Addr:    cfg.Web.APIHost,
		Handler: srv,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.Infow("API listener starting", "addr", cfg.Web.APIHost)
		serverErrors <- app.ListenAndServe()
	}()

	// Shutdown
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logger.Infow("Start shutdown", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := app.Shutdown(ctx)
		if err != nil {
			logger.Infow("Graceful shutdown did not complete", "err", err)
			err = app.Close()
		}

		if err != nil {
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
