package ginV1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skantay/hezzl/config"
	"github.com/skantay/hezzl/internal/controller"
	"github.com/skantay/hezzl/internal/domain"
	"go.uber.org/zap"
)

type ginController struct {
	service domain.Service
	log     *zap.Logger
	cfg     config.Config
}

func New(
	service domain.Service,
	log *zap.Logger,
	cfg config.Config,
) controller.RunController {
	return ginController{
		service: service,
		log:     log,
		cfg:     cfg,
	}
}

func (g ginController) Run() error {
	r := gin.Default()

	r.GET("/goods/list", g.goodsListHandler)
	r.PATCH("/good/reprioritize", g.reprioritizeGoodHandler)
	r.PATCH("/good/update", g.updateGoodHandler)
	r.DELETE("/good/remove", g.removeGoodHandler)
	r.POST("/good/create", g.createGoodHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", g.cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			g.log.Sugar().Fatalf("listen: %s\n", err)
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		g.log.Sugar().Infof("Graceful shutdown failed: %s\n", err)
		return err
	}

	g.log.Info("Server gracefully shut down")
	return nil
}
