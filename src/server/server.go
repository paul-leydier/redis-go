// Core server logic --------------------------------------

package server

import (
	"bytes"
	"log"
	"net"
)

type Server struct {
	storage *storage
}

func NewServer() *Server {
	return &Server{storage: newStorage()}
}

func (s *Server) Run(network string, url string, port string) {
	l, err := net.Listen(network, url+":"+port)
	if err != nil {
		log.Fatalf("could not bind to port - %s", err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Panicf("error while closing the net.Listener - %s", err)
		}
	}(l)
	s.listen(l)
}
func (s *Server) listen(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("could not accept connection - %s", err)
		}
		go s.Serve(conn)
	}
}

func (s *Server) Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("could not close connection - %s", err)
		}
	}(conn)
	for {
		msg := make([]byte, 64)
		_, err := conn.Read(msg)
		if err != nil {
			log.Fatalf("could not read message - %s", err)
		}
		msg = bytes.Trim(msg, "\x00")
		response, err := s.handleMessage(msg)
		if err != nil {
			response = handleServingError(err)
		}
		_, err = conn.Write(response)
		if err != nil {
			log.Fatalf("could not write response - %s", err)
		}
	}
}
