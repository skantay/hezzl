package v

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NC interface {
	SendJSON(ctx context.Context, data any) error
}

type natsSend struct {
	nc *nats.Conn
}

func New(nc *nats.Conn) NC {
	return natsSend{nc}
}

func (n natsSend) SendJSON(ctx context.Context, data any) error {
	ec, err := nats.NewEncodedConn(n.nc, "json")
	if err != nil {
		return fmt.Errorf("trouble with encoding nats: %w", err)
	}

	if err := ec.Publish("Goods.Collection", data); err != nil {
		return fmt.Errorf("trouble publish to nats: %w", err)
	}

	return nil
}
