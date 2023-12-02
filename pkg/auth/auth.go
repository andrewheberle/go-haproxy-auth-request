package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andrewheberle/go-haproxy-auth-request/internal/pkg/spop"
	"github.com/negasus/haproxy-spoe-go/action"
	"github.com/negasus/haproxy-spoe-go/request"
)

type AuthHandler struct {
	client  *http.Client
	headers []string
	method  string
	timeout time.Duration
	url     *url.URL
}

func NewHandler(endpoint, method string, timeout time.Duration, headers []string) (*AuthHandler, error) {
	// parse url
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("problem parsing url: %w", err)
	}

	// set up client to not follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// handle nil headers
	if headers == nil {
		headers = make([]string, 0)
	}

	return &AuthHandler{client, headers, method, timeout, u}, nil
}

func (auth *AuthHandler) Handler(req *request.Request) {
	logger := slog.With("sid", req.StreamID)
	msg, err := req.Messages.GetByName("auth-request")
	if err != nil {
		logger.Info("auth-request message not found")
		return
	}

	// grab headers arg from message
	h, ok := msg.KV.Get("headers")
	if !ok {
		logger.Warn("headers key not found")
		return
	}

	// check that its the correct type
	b, ok := h.([]byte)
	if !ok {
		logger.Error("value of the headers key was not the expected type")
		return
	}

	// parse into http.Headers struct
	headers, err := spop.ParseBinaryHeader(b)
	if err != nil {
		logger.Error("could not parse headers key value", "error", err)
		return
	}

	// debugging of headers
	slog.Debug("request headers", "headers", headers)

	// set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), auth.timeout)
	defer cancel()

	// set up http request
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, auth.url.String(), nil)
	if err != nil {
		logger.Error("creating auth request failed", "error", err)
		return
	}
	r.Header = headers

	// do request
	res, err := auth.client.Do(r)
	if err != nil {
		logger.Error("performing auth request failed", "error", err)
		return
	}

	// return status
	req.Actions.SetVar(action.ScopeTransaction, "response_code", res.StatusCode)
	logger = logger.With("response_code", res.StatusCode)

	// success if ok
	if res.StatusCode == http.StatusOK {
		// signal successful auth
		req.Actions.SetVar(action.ScopeTransaction, "response_successful", true)
		logger = logger.With("response_successful", true)

		// set variables in response
		for _, h := range auth.headers {
			if v := res.Header.Get(h); v != "" {
				k := fmt.Sprintf("response_header.%s", normalise(h))
				req.Actions.SetVar(action.ScopeRequest, k, v)
				logger = logger.With(k, v)
			}
		}

		logger.Info("message handled")
		return
	}

	// otherwise auth is not successful
	logger = logger.With("response_successful", false)

	// handle or access denied redirect
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusFound || res.StatusCode == http.StatusSeeOther {
		// check if location is provided
		if location := res.Header.Get("location"); location != "" {
			logger = logger.With("response_redirect", true, "response_location", location)
			req.Actions.SetVar(action.ScopeTransaction, "response_redirect", true)
			req.Actions.SetVar(action.ScopeTransaction, "response_location", location)
		}
	}

	// all other responses
	req.Actions.SetVar(action.ScopeTransaction, "response_successful", false)
	logger.Info("message handled")
}

// normalise converts the provided string to lowercase and replaces dashes with underscores
func normalise(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), "-", "_")
}
