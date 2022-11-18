package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
)

func InitService(
	ctx context.Context,
	config *Config,
	logger log.Logger,
	exitFn func(),
	diContainter *DIContainer,
) *Service {
	return &Service{
		config:      config,
		logger:      logger,
		exitFn:      exitFn,
		diContainer: diContainter,
	}
}

type Service struct {
	config      *Config
	logger      log.Logger
	exitFn      func()
	diContainer *DIContainer
}

func (s *Service) Start() {
	s.logger.Log("init rest", "Start rest server ...")
	s.logger.Log("init rest", http.ListenAndServe(s.config.HttpPort, s.diContainer.infra.httpHandler))
}

func (s *Service) WaitForStopSignal() {
	s.logger.Log("init rest", "waiting for stop signal ...")

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-stopSignalChannel

	if _, ok := sig.(os.Signal); ok {
		s.logger.Log("init rest", "received", sig)
		close(stopSignalChannel)
		s.shutdown()
	}
}

func (s *Service) shutdown() {
	s.logger.Log("init rest", "shutdown: stopping services ...")
	// stop service properly
	// like db connection

	s.logger.Log("init rest", "shutdown: all services stopped!")

	s.exitFn()
}
