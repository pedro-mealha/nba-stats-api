package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WeNeedThePoh/nba-stats-api/internal/app/domain/stats"
	"github.com/WeNeedThePoh/nba-stats-api/internal/app/gateway"
	"github.com/WeNeedThePoh/nba-stats-api/internal/app/gateway/nba"
	"github.com/WeNeedThePoh/nba-stats-api/internal/app/http/rest"
	"github.com/kelseyhightower/envconfig"
	"github.com/relistan/rubberneck"
	"go.uber.org/zap"
)

type (
	config struct {
		Debug    bool
		LogLevel string `split_words:"true"`
		Web      struct {
			APIHost         string        `split_words:"true" default:"0.0.0.0:8080"`
			ReadTimeout     time.Duration `split_words:"true" default:"30s"`
			WriteTimeout    time.Duration `split_words:"true" default:"2m"`
			IdleTimeout     time.Duration `split_words:"true" default:"5s"`
			ShutdownTimeout time.Duration `split_words:"true" default:"30s"`
		}
		NBA struct {
			CDNBaseURL string        `split_words:"true" required:"true"`
			BaseURL    string        `split_words:"true" required:"true"`
			Timeout    time.Duration `default:"120s"`
		}
	}

	// Notifier holds the context and channels to listen to the notifications
	Notifier struct {
		done chan struct{}
		sig  chan os.Signal
	}
)

func main() {
	l, _ := zap.NewProduction()
	defer l.Sync() //nolint: errcheck

	logger := l.Sugar()
	ctx := context.Background()

	if err := run(ctx, logger); err != nil {
		logger.Fatal(err)
	}
}

// nolint
func run(ctx context.Context, logger *zap.SugaredLogger) error {
	defer logger.Info("completed")

	// =========================================================================
	// Configuration
	// =========================================================================
	logger.Info("Loading configs")

	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to load the env vars: %w", err)
	}

	printer := rubberneck.NewPrinter(logger.Infof, rubberneck.NoAddLineFeed)
	printer.Print(cfg)

	// =========================================================================
	// Config rest client
	// =========================================================================
	var (
		nbaClient = gateway.NewClientWithTimeout(cfg.NBA.Timeout)
		n         = nba.New(nbaClient, cfg.NBA.BaseURL, cfg.NBA.CDNBaseURL)
	)

	// =========================================================================
	// Start Server
	// =========================================================================
	logger.Info("initializing REST server")

	var (
		serverErrors = make(chan error, 1)
		rs           = stats.NewService(n)
		a            = rest.NewAPI(logger, rs)
	)

	server := &http.Server{
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      a.Routes(),
		BaseContext:  func(net.Listener) context.Context { return ctx },
	}

	go func() {
		logger.Infow("Initializing API", "host", cfg.Web.APIHost)
		serverErrors <- server.ListenAndServe()
	}()

	done := newSignal(ctx)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		logger.Infow("start shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Errorw("failed to gracefully shutdown the server", "err", err)

			if err = server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}
	}

	return nil
}

func newSignal(ctx context.Context, signals ...os.Signal) *Notifier {
	if signals == nil {
		// default signals
		signals = []os.Signal{
			os.Interrupt,
			syscall.SIGTERM,
		}
	}

	signaler := Notifier{
		done: make(chan struct{}),
		sig:  make(chan os.Signal),
	}

	signal.Notify(signaler.sig, signals...)

	go signaler.listenToSignal(ctx)

	return &signaler
}

// listenToSignal is a blocking statement that listens to two channels:
//
//   - s.sig: is the os.Signal that will the triggered by the signal.Notify once
//     the expected signals are executed by the OS in the service
//   - ctx.Done(): in case of close of context, the service should also shutdown
func (s *Notifier) listenToSignal(ctx context.Context) {
	for {
		select {
		case <-s.sig:
			s.done <- struct{}{}
			return
		case <-ctx.Done():
			s.done <- struct{}{}
			return
		}
	}
}

// Done returns the call of the done channel
func (s *Notifier) Done() <-chan struct{} { return s.done }
