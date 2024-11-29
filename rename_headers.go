package traefik_custom_headers_plugin

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

// Rename holds one rename configuration.
type rename struct {
	HeaderName    string `json:"headerName"`
	NewHeaderName string `json:"newHeaderName"`
}

// Config holds the plugin configuration.
type Config struct {
	Rename []rename `json:"renameData"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates and returns a new rewrite body plugin instance.
type renameHeaders struct {
	name    string
	next    http.Handler
	renames []rename
}

// New creates a new Custom Header plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	slog.Info("[plugin/rename-headers] in New")
	return &renameHeaders{
		name:    name,
		next:    next,
		renames: config.Rename,
	}, nil
}

func (r *renameHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	slog.Info("[plugin/rename-headers] in ServeHTTP")
	slog.Info("[plugin/rename-headers] in ServeHTTP, request headers", "kv", req.Header)

	wrappedWriter := &responseWriter{
		writer:          rw,
		headersToRename: r.renames,
	}

	rw.Header().Set("X-Test", "aaa")
	slog.Info("[plugin/rename-headers] in ServeHTTP, response headers", "kv", rw.Header())

	r.next.ServeHTTP(wrappedWriter, req)
}

type responseWriter struct {
	writer          http.ResponseWriter
	headersToRename []rename
}

func (r *responseWriter) Header() http.Header {
	slog.Info("[plugin/rename-headers] in Header")
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	slog.Info("[plugin/rename-headers] in Write")
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	slog.Info("[plugin/rename-headers] in WriteHeader", "headersToRename", r.headersToRename)
	for _, headerToRename := range r.headersToRename {
		r.Header().Set(headerToRename.NewHeaderName, headerToRename.HeaderName)
		r.Header().Del(headerToRename.HeaderName)
	}
	slog.Info("[plugin/rename-headers] in WriteHeader, response headers", "kv", r.Header())

	r.writer.WriteHeader(statusCode)
}

func (r *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	slog.Info("[plugin/rename-headers] in Hijack")
	hijacker, ok := r.writer.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", r.writer)
	}

	return hijacker.Hijack()
}

func (r *responseWriter) Flush() {
	slog.Info("[plugin/rename-headers] in Flush")
	if flusher, ok := r.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
