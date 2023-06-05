package eventservice

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrSwartz/event/internal/config"
	"github.com/MrSwartz/event/pkg/eventservice/service"

	"github.com/sirupsen/logrus"
)

func Run() error {
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
		return err
	}

	handler := &Handler{
		Service: srvc,
	}

	httpSrv := new(HttpSrv)
	go func() {
		if err := httpSrv.Run(&cnf.Service, handler.InitRoutes(cnf.Service.ExposeSwagger)); err != nil {
			logrus.Errorf("can't to start http server %v", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// сначала закрываю http сервер, но не закрываю слой данных,
	// чтобы дать шанс выгрузить данные из буфера
	if err := httpSrv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
		return err
	}

	if err := handler.Service.CloseData(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
		return err
	}
	logrus.Info("data closed")

	logrus.Info("service stopped")
	return nil
}
