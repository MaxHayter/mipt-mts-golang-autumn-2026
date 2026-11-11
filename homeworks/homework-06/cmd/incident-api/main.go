package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homework06/internal/client/incident"
	httpapi "homework06/internal/transport/http"
)

func main() {
	if err := run(); err != nil {
		slog.Error("incident-api stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	addr := flag.String("http-addr", "localhost:8080", "HTTP listen address")
	grpcAddr := flag.String("grpc-addr", "localhost:50051", "incident-service address")
	timeout := flag.Duration("unary-timeout", 3*time.Second, "unary RPC timeout")
	flag.Parse()
	if *timeout <= 0 {
		return fmt.Errorf("unary-timeout must be positive")
	}
	client, err := incident.New(*grpcAddr)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("close grpc client", "error", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	server := &http.Server{Addr: *addr, Handler: httpapi.New(client, *timeout), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second, BaseContext: func(net.Listener) context.Context { return ctx }}
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	slog.Info("incident-api started", "address", *addr, "grpc_address", *grpcAddr)
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			if closeErr := server.Close(); closeErr != nil {
				return errors.Join(err, closeErr)
			}
			return err
		}
		return nil
	}
}
