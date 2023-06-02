package main

import (
	"context"
	"event/internal/config"
	"event/pkg/eventservice"
	"event/pkg/eventservice/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cnf, err := config.ReadConfig()
	if err != nil {
		logrus.Fatalf("can't read config: %v", err.Error())
	}
	logrus.Info("config loaded")

	ctx := context.Background()

	srvc, err := service.NewService(ctx, *cnf)
	if err != nil {
		logrus.Errorf("can't create service layer, error: %v", err.Error())
		return
	}

	handler := &eventservice.Handler{
		Service: srvc,
	}

	httpSrv := new(eventservice.HttpSrv)
	go func() {
		if err := httpSrv.Run(&cnf.Service, handler.InitRoutes()); err != nil {
			logrus.Errorf("can't to start http server %v", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := handler.Service.CloseData(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
		return
	}
	logrus.Info("data closed")

	if err := httpSrv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
		return
	}
	logrus.Info("service stopped")
}
