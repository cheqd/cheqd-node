package pricefeeder

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/cheqd/cheqd-node/x/oracle/types"

	"github.com/ojo-network/price-feeder/config"
	"github.com/ojo-network/price-feeder/oracle"
	"github.com/ojo-network/price-feeder/oracle/client"
	v1 "github.com/ojo-network/price-feeder/router/v1"
)

type PriceFeeder struct {
	Oracle    *oracle.Oracle
	AppConfig AppConfig
}

func (pf *PriceFeeder) Start(currentBlockHeight int64, oracleParams types.Params) error {
	logWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	logLevel, err := zerolog.ParseLevel(pf.AppConfig.LogLevel)
	if err != nil {
		return err
	}
	logger := zerolog.New(logWriter).Level(logLevel).With().Timestamp().Logger()
	cfg, err := config.LoadConfigFromFlags(pf.AppConfig.ConfigPath, "")
	if err != nil {
		return err
	}
	// listen for and trap any OS signal to gracefully shutdown and exit
	ctx, cancel := context.WithCancel(context.TODO())
	g, ctx := errgroup.WithContext(ctx)

	trapSignal(cancel, logger)

	providerTimeout, err := time.ParseDuration(cfg.ProviderTimeout)
	if err != nil {
		return fmt.Errorf("failed to parse provider timeout: %w", err)
	}

	providers := oracle.CreatePairProvidersFromCurrencyPairProvidersList(oracleParams.CurrencyPairProviders)
	deviations, err := oracle.CreateDeviationsFromCurrencyDeviationThresholdList(oracleParams.CurrencyDeviationThresholds)
	if err != nil {
		return err
	}
	o := oracle.New(
		logger,
		client.OracleClient{},
		providers,
		providerTimeout,
		deviations,
		cfg.ProviderEndpointsMap(),
		true,
	)
	pf.SetOracle(o)

	var metrics *telemetry.Metrics

	// check if app telemetry is already running
	if !telemetry.IsTelemetryEnabled() {
		telemetryCfg := telemetry.Config{}
		err = mapstructure.Decode(cfg.Telemetry, &telemetryCfg)
		if err != nil {
			return err
		}

		metrics, err = telemetry.New(telemetryCfg)
		if err != nil {
			return err
		}
	} else {
		// disable price feeder telemetry and use app telemetry instead
		cfg.Telemetry.Enabled = false
		metrics = nil
	}

	g.Go(func() error {
		// start the process that observes and publishes exchange prices
		return startPriceFeeder(ctx, logger, cfg, pf.Oracle, metrics)
	})

	// Block main process until all spawned goroutines have gracefully exited and
	// signal has been captured in the main process or if an error occurs.
	return g.Wait()
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func trapSignal(cancel context.CancelFunc, logger zerolog.Logger) {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		logger.Info().Str("signal", sig.String()).Msg("caught signal; shutting down...")
		cancel()
	}()
}

func (pf *PriceFeeder) SetOracle(o *oracle.Oracle) {
	pf.Oracle = o
}

func (pf *PriceFeeder) GetOracle() *oracle.Oracle {
	return pf.Oracle
}

// startPriceFeeder starts the price feeder server which listens to websocket connections
// from price providers.
func startPriceFeeder(
	ctx context.Context,
	logger zerolog.Logger,
	cfg config.Config,
	oracle *oracle.Oracle,
	metrics *telemetry.Metrics,
) error {
	rtr := mux.NewRouter()
	v1Router := v1.New(logger, cfg, oracle, metrics)
	v1Router.RegisterRoutes(rtr, v1.APIPathPrefix)

	writeTimeout, err := time.ParseDuration(cfg.Server.WriteTimeout)
	if err != nil {
		return err
	}
	readTimeout, err := time.ParseDuration(cfg.Server.ReadTimeout)
	if err != nil {
		return err
	}

	srvErrCh := make(chan error, 1)
	srv := &http.Server{
		Handler:           rtr,
		Addr:              cfg.Server.ListenAddr,
		WriteTimeout:      writeTimeout,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readTimeout,
	}

	go func() {
		logger.Info().Str("listen_addr", cfg.Server.ListenAddr).Msg("starting price-feeder server...")
		srvErrCh <- srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			logger.Info().Str("listen_addr", cfg.Server.ListenAddr).Msg("shutting down price-feeder server...")
			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Error().Err(err).Msg("failed to gracefully shutdown price-feeder server")
				return err
			}

			return nil

		case err := <-srvErrCh:
			logger.Error().Err(err).Msg("failed to start price-feeder server")
			return err
		}
	}
}
