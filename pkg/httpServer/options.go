package httpserver

import (
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func Address(address string) (Option, error) {
	return func(s *Server) {
		s.server.Addr = address
	}, nil
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) (Option, error) {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}, nil
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) (Option, error) {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}, nil
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) (Option, error) {
	return func(s *Server) {
		s.ShutdownTimeout = timeout
	}, nil
}
