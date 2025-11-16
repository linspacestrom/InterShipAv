package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/InterShipAv/internal/config"
	"github.com/linspacestrom/InterShipAv/internal/services"
	"go.uber.org/zap"
)

type Server struct {
	srv  *http.Server
	logg *zap.Logger
}

func NewServer(cfg *config.Config, teamSvc services.TeamSer, userSvc services.UserSer, logg *zap.Logger) *Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// инициализация хэндлеров
	NewTeamHandler(router, teamSvc, logg)
	NewUserHandler(router, userSvc, logg)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
		//ReadTimeout:  cfg.ServerReadTimeout,
		//WriteTimeout: cfg.ServerWriteTimeout,
		//IdleTimeout:  cfg.ServerIdleTimeout,
	}

	return &Server{
		srv:  srv,
		logg: logg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.logg.Info("starting HTTP server", zap.String("addr", s.srv.Addr))

	errc := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- err
		} else {
			errc <- nil
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.logg.Info("HTTP server shutting down")
		return s.srv.Shutdown(shutdownCtx)
	case err := <-errc:
		if err != nil {
			s.logg.Error("HTTP server error", zap.Error(err))
		}
		return err
	}
}
