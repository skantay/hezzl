package v

import (
	"context"
	"fmt"
	"runtime"

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
	_, err := n.nc.Subscribe("Goods.Collection", func(msg *nats.Msg) {
		n.log.Sugar().Infof("got a message")

		fmt.Println(string(msg.Data))
		if err := n.service.Good.Create(msg.Data); err != nil {
			n.log.Error(err.Error())
			return
		}

		n.log.Sugar().Infof("created a log")
	})
	if err != nil {
		return err
	}

	runtime.Goexit()

	<-ctx.Done()
	return ctx.Err()
}
