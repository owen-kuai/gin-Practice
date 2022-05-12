package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"practice/pkg/app/api"
	"syscall"
	"time"

	"practice/pkg/conf"
	"practice/pkg/logger"
)

//RunServer run the backend server
func RunServer(configPath string) {
	// init config
	conf.InitConfig(configPath, "config", "yaml")

	// init logger
	logger.InitLogger()
	mainLogger := logger.Logger("main")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go cancelOnInterrupt(ctx, cancelFunc, mainLogger)

	if err := runWithContext(ctx, mainLogger); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		mainLogger.FatalError(err, "run api server failed")
	}
	return
}

//cancelOnInterrupt cancel server by listen to the sys signal
func cancelOnInterrupt(ctx context.Context, f context.CancelFunc, logger logger.AppLogger) {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-term:
			logger.Info("received SIGTERM, exiting gracefully...")
			f()
			os.Exit(0)
		case <-ctx.Done():
			os.Exit(0)
		}
	}
}

//runWithContext run gin server with context
func runWithContext(ctx context.Context, logger logger.AppLogger) error {

	// init gin server
	router := api.New()

	config := conf.GetConfig()
	bind := fmt.Sprintf("%s:%d", config.Host, config.Port)

	srv := &http.Server{
		Addr:    bind,
		Handler: router,
	}
	logger.Infof("start to server, bind : %s ", bind)

	go func() {
		<-ctx.Done()
		// default 3 second to gracefully shut down the server
		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if srv.Shutdown(c) != nil {
			srv.Close()
		}
	}()

	return srv.ListenAndServe()
}
