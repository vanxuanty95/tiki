package main

import (
	"context"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tiki/cmd/api/config"
	"tiki/internal/api"
	"tiki/internal/pkg/logger"
	"time"
)

var tikiLogger logger.Logger

func main() {
	state := flag.String("state", "local", "state of service")
	tikiLogger = logger.WithPrefix("main")
	cfg, err := config.Load(state)
	if err != nil {
		tikiLogger.Panicln(err)
	}
	var server *http.Server
	go func() {
		server = initRestfulAPI(cfg)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			tikiLogger.Panicf("Fail to listen and server: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-signals

	tikiLogger.Info("shutting server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		tikiLogger.Errorf("Fail to listen and server: %v", err)
	}
	tikiLogger.Info("shutdown server")
}

func initRestfulAPI(cfg *config.Config) *http.Server {
	tikiLogger.Info("Start server")
	tikiLogger.Infof("%s:%s", cfg.RestfulAPI.Host, cfg.RestfulAPI.Port)
	server, err := api.CreateAPIEngine(cfg)
	if err != nil {
		tikiLogger.Panicf("Fail init server: %v", err)
		return nil
	}
	return server
}
