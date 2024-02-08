package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/cossim/kit/pkg/config"
	"github.com/cossim/kit/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type HttpService struct {
	server  *http.Server
	logger  *zap.Logger
	handler gin.IRouter
	cfg     *config.Config

	svc HTTPService
}

func NewHttpService(c *config.Config, svc HTTPService) *HttpService {
	s := &HttpService{
		logger: log.NewDevLogger(svc.Name()),
		cfg:    c,

		svc: svc,
	}

	handler := gin.Default()
	s.handler = handler
	s.server = &http.Server{
		Handler:           handler,
		Addr:              s.cfg.HTTP.Addr(),
		MaxHeaderBytes:    1 << 20,
		IdleTimeout:       90 * time.Second, // matches http.DefaultTransport keep-alive timeout
		ReadHeaderTimeout: 32 * time.Second,
	}
	return s
}

func (s *HttpService) Start(ctx context.Context) error {
	if err := s.svc.Init(s.cfg); err != nil {
		return err
	}

	s.svc.Registry(s.handler)

	s.logger.Info(fmt.Sprintf("%s http service start", s.cfg.Register.Name), zap.String("addr", s.cfg.HTTP.Addr()))
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.logger.Info(fmt.Sprintf("%s http service stop", s.cfg.Register.Name), zap.String("addr", s.cfg.HTTP.Addr()))
			return nil
		}
		return fmt.Errorf("启动 [%s] http服务失败：%v", s.cfg.Register.Name, err)
	}
	return nil
}

func (s *HttpService) Stop() error {
	s.logger.Info(fmt.Sprintf("%s http service try stop", s.cfg.Register.Name), zap.String("addr", s.cfg.HTTP.Addr()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Warn(fmt.Sprintf("%s http service try stop", s.cfg.Register.Name), zap.String("addr", s.cfg.HTTP.Addr()))
		return err
	}
	return nil
}

func (s *HttpService) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start(context.Background())
}
