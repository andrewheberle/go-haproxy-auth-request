package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/andrewheberle/go-haproxy-auth-request/internal/pkg/logger"
	"github.com/andrewheberle/go-haproxy-auth-request/pkg/auth"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// command line flags
	pflag.String("listen", "127.0.0.1:3000", "Listen address")
	pflag.String("url", "http://127.0.0.1:9091/api/authz/forward-auth", "URL to perform verification against")
	pflag.String("method", http.MethodHead, "HTTP method for verification request")
	pflag.StringSlice("headers", []string{"authorization", "proxy-authorization", "remote-user", "remote-groups", "remote-name", "remote-email"}, "HTTP Headers to return on success")
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

	// set up SPOA handler
	h, err := auth.NewHandler(viper.GetString("url"), viper.GetString("method"), viper.GetDuration("timeout"), viper.GetStringSlice("headers"))
	if err != nil {
		slog.Error("error creating auth handler", "error", err)
		os.Exit(1)
	}

	// set up SPOA
	a := agent.New(h.Handler, logger.NewLogger())

	// start serving traffic
	if err := a.Serve(l); err != nil {
		slog.Error("agent serve error", "error", err)
	}
}
