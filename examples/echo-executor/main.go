// Echo executor: demonstrates a minimal Archon-compatible need worker.
//
// Usage:
//
//	go run . -nats nats://127.0.0.1:4222 -subject need.echo -api http://127.0.0.1:8080
//
// This is intentionally small so newcomers see the whole loop:
// subscribe → parse need → work → webhook → ack.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	needprotocol "github.com/diogoX451/archon-need-protocol"
	executorruntime "github.com/diogoX451/archon-executor-runtime"
	natsbus "github.com/diogoX451/archon-nats-bus"
)

func main() {
	natsURL := flag.String("nats", env("NATS_URL", "nats://127.0.0.1:4222"), "NATS URL")
	subject := flag.String("subject", env("NEED_SUBJECT", "need.echo"), "need subject")
	api := flag.String("api", env("ORCHESTRATOR_URL", "http://127.0.0.1:8080"), "orchestrator base URL")
	flag.Parse()

	bus, err := natsbus.New(natsbus.Config{URL: *natsURL, MaxReconnects: -1, ReconnectWait: 2 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_, err = bus.Subscribe(*subject, func(ctx context.Context, msg natsbus.Message) error {
		ev, err := needprotocol.ParseNeedEvent(msg.Data())
		if err != nil {
			log.Printf("bad payload: %v", err)
			return msg.Ack() // permanent poison
		}
		if ev.Need.CorrelationID == "" {
			log.Printf("missing correlation_id")
			return msg.Ack()
		}
		// Work: echo the need payload with a timestamp.
		out, _ := json.Marshal(map[string]any{
			"echo":      json.RawMessage(ev.Need.Payload),
			"type":      ev.Need.Type,
			"processed": time.Now().UTC().Format(time.RFC3339),
		})
		if err := executorruntime.PostWebhook(ctx, executorruntime.WebhookConfig{
			BaseURL:       *api,
			DeliveryToken: os.Getenv("DELIVERY_TOKEN"),
		}, ev.Need.CorrelationID, out); err != nil {
			log.Printf("webhook: %v", err)
			return err // transient → redeliver
		}
		return msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("echo-executor listening subject=%s api=%s", *subject, *api)
	<-ctx.Done()
	fmt.Println("shutdown")
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
