package main

import (
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/andrewheberle/go-haproxy-auth-request/internal/pkg/logger"
	"github.com/andrewheberle/go-haproxy-auth-request/pkg/auth"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// logging
	// command line flags
	pflag.String("listen", "127.0.0.1:3000", "Listen address")
	pflag.String("cookie", "authelia_session", "Session cookie name")
	pflag.String("url", "http://127.0.0.1:9091/api/verify", "URL to perform verification against")
	pflag.Duration("timeout", time.Second*5, "Timeout for verification")
	pflag.Bool("debug", false, "Enable debug logging")
	pflag.Parse()

	// bind to viper
	viper.BindPFlags(pflag.CommandLine)

	// load from environment
	viper.SetEnvPrefix("auth")
	viper.AutomaticEnv()

	// logging setup
	var logLevel = new(slog.LevelVar)
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logHandler))
	if viper.GetBool("debug") {
		logLevel.Set(slog.LevelDebug)
	}

	// some logging
	slog.Info("starting", "listen", viper.GetString("listen"), "url", viper.GetString("url"))

	// set up listener
	l, err := net.Listen("tcp", viper.GetString("listen"))
	if err != nil {
		slog.Error("error creating listener", "error", err)
		os.Exit(1)
	}
	defer l.Close()

	h, err := auth.NewHandler(viper.GetString("url"), viper.GetString("cookie"), viper.GetDuration("timeout"))
	if err != nil {
		slog.Error("error creating auth handler", "error", err)
		os.Exit(1)
	}

	a := agent.New(h.Handler, logger.NewLogger())

	if err := a.Serve(l); err != nil {
		slog.Error("agent serve error", "error", err)
	}
}
