package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skantay/hezzl/config"
	"github.com/skantay/hezzl/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type Controller interface {
	Serve(ctx context.Context) error
}

type ginController struct {
	service   usecase.Service
	log       *zap.Logger
	cfg       config.Config
	validator *validator.Validate
}

func New(
	service usecase.Service,
	log *zap.Logger,
	cfg config.Config,
	validator *validator.Validate,
) Controller {
	return ginController{
		service:   service,
		log:       log,
		cfg:       cfg,
		validator: validator,
	}
}

func (g ginController) Serve(ctx context.Context) error {
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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		g.log.Sugar().Fatalf("listen: %s\n", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		g.log.Sugar().Infof("Graceful shutdown failed: %s\n", err)
		return err
	}

	return nil
}
