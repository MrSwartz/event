package eventservice

import (
	"context"
	"event/internal/config"
	"fmt"
	"net/http"
	"time"
)

type HttpSrv struct {
	httpSrv *http.Server
}

func (s *HttpSrv) Run(cnf *config.Service, h http.Handler) error {
	s.httpSrv = &http.Server{
		Addr:           fmt.Sprintf(":%d", cnf.Port),
		Handler:        h,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(cnf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cnf.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(cnf.IdleTimeout) * time.Second,
	}

	// return
	s.httpSrv.ListenAndServe()
	return nil
}

func (s *HttpSrv) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
