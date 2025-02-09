package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/DominikKoniarz/some-tcp-server/internal/auth"
	"github.com/DominikKoniarz/some-tcp-server/internal/connection"
	"github.com/DominikKoniarz/some-tcp-server/internal/request"
)

type Server struct {
	S           *net.Listener
	Logger      *log.Logger
	Connections map[string]*connection.Connection
	NextConnID  int
	Mu          sync.Mutex
}

func (s *Server) AddConnection(conn net.Conn) *connection.Connection {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	connID := fmt.Sprintf("%d", s.NextConnID)
	s.NextConnID++

	connection := &connection.Connection{
		ID:              connID,
		IsAuthenticated: false,
		C:               &conn,
	}

	s.Connections[connID] = connection

	return connection
}

func (s *Server) GetConnection(connID string) (*connection.Connection, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	conn, ok := s.Connections[connID]
	return conn, ok
}

func (s *Server) RemoveConnection(connID string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Connections, connID)
}

func NewServer() *Server {
	server, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	s := &Server{
		S:           &server,
		Logger:      log.New(log.Writer(), "Server: ", log.LstdFlags),
		Connections: make(map[string]*connection.Connection),
		NextConnID:  0,
		Mu:          sync.Mutex{},
	}

	return s
}

func (s *Server) Start() {
	s.Logger.Println("Server started")

	for {
		conn, err := (*s.S).Accept()
		if err != nil {
			s.Logger.Println("Error accepting connection:", err)
			continue
		}

		c := s.AddConnection(conn)

		go s.handleConnection(c)
	}
}

func (s *Server) Stop() {
	s.Logger.Println("Stopping server")
	(*s.S).Close()
}

func (s *Server) handleConnection(conn *connection.Connection) {
	s.Logger.Println("Connection established with ID:", conn.ID)
	s.Logger.Println("Client address:", (*conn.C).RemoteAddr())

	defer func() {
		(*conn.C).Close()
		s.Logger.Println("Connection closed with ID:", conn.ID)
		s.RemoveConnection(conn.ID)
	}()

	buf := make([]byte, 1024)

	for {
		n, err := (*conn.C).Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				s.Logger.Println("Connection closed by client:", (*conn.C).RemoteAddr())
				break
			}

			s.Logger.Println("Error reading data:", err)
			break
		}

		fmt.Println("Received data:", string(buf[:n]))

		parsedRequest, err := request.ParseRequest(buf[:n])
		if err != nil {
			s.Logger.Println("Error parsing request:", err)

			if _, err := (*conn.C).Write([]byte(err.Error())); err != nil {
				s.Logger.Println("Error writing response:", err)
				break
			}

			continue
		}

		fmt.Println("Parsed request:", parsedRequest)

		if !conn.IsAuthenticated {
			err := auth.HandleAuth(conn, &parsedRequest)
			if err != nil {
				s.Logger.Println("Authentication error:", err)
				break
			}
		} else {
			// write request back to client
			if _, err := (*conn.C).Write(parsedRequest.ToBytes()); err != nil {
				s.Logger.Println("Error writing response:", err)
				break
			}
		}
	}
}
