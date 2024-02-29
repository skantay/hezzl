package v

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/skantay/service-2/internal/usecase"
	"go.uber.org/zap"
)

type NC interface {
	Serve(ctx context.Context) error
}

type natsServe struct {
	nc      *nats.Conn
	log     *zap.Logger
	service usecase.Service
}

func New(nc *nats.Conn, log *zap.Logger, ser usecase.Service) NC {
	return natsServe{nc, log, ser}
}

func (n natsServe) Serve(ctx context.Context) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	if !n.nc.IsConnected() {
		n.log.Error("NATS connection is not established")
		return errors.New("NATS connection is not established")
	}
	go func() {
		n.nc.Subscribe("goods", func(msg *nats.Msg) {
			n.log.Sugar().Infof("got a message")
			if err := n.service.Good.Create(msg.Data); err != nil {
				n.log.Error(err.Error())
				return
			}
			n.log.Sugar().Infof("created a log")
		})
	}()

	<-shutdown
	n.log.Info("Server gracefully shut down")
	return nil
}
