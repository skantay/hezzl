package controller

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	service "github.com/skantay/service-2/internal/domain"
	"go.uber.org/zap"
)

type Controller interface {
	Run() error
}

type controller struct {
	serviceGood service.Service
	log         *zap.Logger
}

func New(service service.Service, log *zap.Logger) Controller {
	return controller{service, log}
}

func (c controller) Run() error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		c.serviceGood.Nc.Subscribe("good", func(msg *nats.Msg) {
			c.log.Info("got message")
			err := c.serviceGood.GoodService.Create(msg.Data)
			if err != nil {
				c.log.Error(err.Error())
			}
		})
	}()

	<-shutdown
	c.log.Info("Server gracefully shut down")
	return nil
}
