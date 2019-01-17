package consumer

import (
	"context"
	"log"
)

type event struct {
	Data []byte
}

// Receive func logs an event payload
func Receive(ctx context.Context, e event) error {
	log.Printf("%s", string(e.Data))

	return nil
}
