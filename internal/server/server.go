package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/DominikKoniarz/some-tcp-server/internal/message"
	"github.com/DominikKoniarz/some-tcp-server/internal/request"
)

type Connection struct {
	ID              string
	IsAuthenticated bool
	C               *net.Conn
}

type Server struct {
	s           *net.Listener
	logger      *log.Logger
	connections map[string]*Connection
	nextConnID  int
	mu          sync.Mutex
}

func (s *Server) AddConnection(conn net.Conn) *Connection {
	s.mu.Lock()
	defer s.mu.Unlock()

	connID := fmt.Sprintf("%d", s.nextConnID)
	s.nextConnID++

	connection := &Connection{
		ID:              connID,
		IsAuthenticated: false,
		C:               &conn,
	}

	s.connections[connID] = connection

	return connection
}

func (s *Server) GetConnection(connID string) (*Connection, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.connections[connID]
	return conn, ok
}

func (s *Server) RemoveConnection(connID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.connections, connID)
}

func NewServer() *Server {
	server, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	s := &Server{
		s:           &server,
		logger:      log.New(log.Writer(), "Server: ", log.LstdFlags),
		connections: make(map[string]*Connection),
		nextConnID:  0,
		mu:          sync.Mutex{},
	}

	return s
}

func (s *Server) Start() {
	s.logger.Println("Server started")

	for {
		conn, err := (*s.s).Accept()
		if err != nil {
			s.logger.Println("Error accepting connection:", err)
			continue
		}

		c := s.AddConnection(conn)

		go s.handleConnection(c)
	}
}

func (s *Server) Stop() {
	s.logger.Println("Stopping server")
	(*s.s).Close()
}

func (s *Server) handleConnection(conn *Connection) {
	s.logger.Println("Connection established with ID:", conn.ID)
	s.logger.Println("Client address:", (*conn.C).RemoteAddr())

	defer func() {
		(*conn.C).Close()
		s.logger.Println("Connection closed with ID:", conn.ID)
		s.RemoveConnection(conn.ID)
	}()

	buf := make([]byte, 1024)

	for {
		n, err := (*conn.C).Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				s.logger.Println("Connection closed by client:", (*conn.C).RemoteAddr())
				break
			}

			s.logger.Println("Error reading data:", err)
			break
		}

		fmt.Println("Received data:", string(buf[:n]))

		parsedRequest, err := request.ParseRequest(buf[:n])
		if err != nil {
			s.logger.Println("Error parsing request:", err)

			if _, err := (*conn.C).Write([]byte(err.Error())); err != nil {
				s.logger.Println("Error writing response:", err)
				break
			}

			continue
		}

		fmt.Println("Parsed request:", parsedRequest)

		// if _, err := (*conn.C).Write(parsedRequest.ToBytes()); err != nil {
		// 	s.logger.Println("Error writing response:", err)
		// }

		if !conn.IsAuthenticated {
			if parsedRequest.MessageType != request.AUTH_MESSAGE_TYPE {
				s.logger.Println("Client not authenticated")
				if err := (*conn.C).Close(); err != nil {
					s.logger.Println("Error closing connection:", err)
				}
				break
			}

			authMessage, err := message.ParseAuthMessage(parsedRequest.Data)
			if err != nil {
				s.logger.Println("Error parsing auth message:", err)
				if _, err := (*conn.C).Write([]byte(err.Error())); err != nil {
					s.logger.Println("Error writing response:", err)
				}
				continue
			}

			if authMessage.Username == "root" && authMessage.Password == "root" {
				conn.IsAuthenticated = true
				s.logger.Println("Client authenticated")
				if _, err := (*conn.C).Write([]byte("Authenticated")); err != nil {
					s.logger.Println("Error writing response:", err)
				}
			} else {
				s.logger.Println("Invalid credentials")
				if _, err := (*conn.C).Write([]byte("Invalid credentials")); err != nil {
					s.logger.Println("Error writing response:", err)
				}

				s.logger.Println("Invalid authentication attempt for user:", authMessage.Username)

				break
			}
		} else {
			// write request back to client
			if _, err := (*conn.C).Write(parsedRequest.ToBytes()); err != nil {
				s.logger.Println("Error writing response:", err)
			}
		}
	}
}
