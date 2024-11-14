
// x := 10       // Declare x
// ptr := &x     // ptr now holds the address of x
// *ptr = 20     // Dereference ptr to update x
// fmt.Println(x) // Prints 20

// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	address string
	mux chi.Router //Mux is short for multiplexer. The mux is what receives an HTTP request, looks at where it should go, and directs it to the code that should give a response. 
	server *http.Server //gives the value stored at the memory address the pointer is pointing to.  
}

type Options struct {
	Host string
	Port int
}



func New(opts Options) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server {
		address: address,
		mux:	mux,
		server: &http.Server {
			Addr:	address,
			Handler:	mux,
			ReadTimeout: 5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout: 5 * time.Second,
		},
		} //& gives memory address
	
} 

//setting up routes and listening for HTTP request on the given address
func (s *Server) Start() error {
	s.setupRoutes()
	
	fmt.Println("Starting on", s.address)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServer) {
		return fmt.Errorf("error starting server: %w", err)
	}
return nil
}

// Stop the server gracefully within the timeout
func (s *Server) Stop() error { //in golang this is a "METHOD" and NOT a "FUNCTION"
	fmt.Println("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)//The syntax ctx, cancel := ... means both values returned by context.WithTimeout are assigned to ctx and cancel, respectively, in a single line.
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}
	return nil
}
