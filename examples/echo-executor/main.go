// Echo executor: demonstrates a minimal Archon-compatible need worker.
//
// Usage:
//
//	go run . -bus-url nats://127.0.0.1:4222 -subject need.echo -api http://127.0.0.1:8080
//
// The bus driver is pluggable (default nats). Domain code never imports a
// broker SDK — only archon-bus.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	archonbus "github.com/diogoX451/archon-bus"
	executorruntime "github.com/diogoX451/archon-executor-runtime"
	_ "github.com/diogoX451/archon-nats-bus"
	needprotocol "github.com/diogoX451/archon-need-protocol"
)

func main() {
	busDriver := flag.String("bus-driver", env("ARCHON_BUS_DRIVER", "nats"), "bus driver")
	busURL := flag.String("bus-url", firstEnv("nats://127.0.0.1:4222", "ARCHON_BUS_URL", "NATS_URL"), "bus URL")
	legacyNATS := flag.String("nats", "", "deprecated: use -bus-url")
	subject := flag.String("subject", env("NEED_SUBJECT", "need.echo"), "need subject")
	api := flag.String("api", env("ORCHESTRATOR_URL", "http://127.0.0.1:8080"), "orchestrator base URL")
	flag.Parse()

	url := strings.TrimSpace(*busURL)
	if v := strings.TrimSpace(*legacyNATS); v != "" {
		log.Printf("warning: -nats is deprecated; use -bus-url")
		url = v
	}

	b, err := archonbus.Open(archonbus.Config{
		Driver:        *busDriver,
		URL:           url,
		MaxReconnects: -1,
		ReconnectWait: 2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_, err = b.Subscribe(*subject, func(ctx context.Context, msg archonbus.Message) error {
		ev, err := needprotocol.ParseNeedEvent(msg.Data())
		if err != nil {
			log.Printf("bad payload: %v", err)
			return msg.Ack()
		}
		if ev.Need.CorrelationID == "" {
			log.Printf("missing correlation_id")
			return msg.Ack()
		}
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
			return err
		}
		return msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("echo-executor driver=%s subject=%s api=%s", *busDriver, *subject, *api)
	<-ctx.Done()
	fmt.Println("shutdown")
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func firstEnv(fallback string, keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return fallback
}
