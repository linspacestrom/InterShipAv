package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func NewServer(cfg *config.Config, teamSvc services.TeamSer, userSvc services.UserSer, prSvc services.PRSer, logg *zap.Logger) *Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	NewTeamHandler(router, teamSvc, logg)
	NewUserHandler(router, userSvc, logg)
	NewPullRequestHandler(router, prSvc, logg)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		s.logg.Info("context canceled, shutting down HTTP server")
	case sig := <-sigChan:
		s.logg.Info("received OS signal, shutting down HTTP server", zap.String("signal", sig.String()))
	case err := <-errc:
		if err != nil {
			s.logg.Error("HTTP server error", zap.Error(err))
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(shutdownCtx)
}
