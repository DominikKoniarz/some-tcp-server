package server

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	s      *net.Listener
	logger *log.Logger
}

type ParsedRequest struct {
	// length          int
	protocolVersion string
	messageType     string
	data            string
}

func (pr ParsedRequest) Stringify() string {
	return fmt.Sprintf("Protocol version: %s, Message type: %s, Data: %s", pr.protocolVersion, pr.messageType, pr.data)
}

func NewServer() *Server {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	return &Server{
		s:      &server,
		logger: log.New(log.Writer(), "Server: ", log.LstdFlags),
	}
}

func (srv *Server) Start() {
	srv.logger.Println("Server started")

	for {
		conn, err := (*srv.s).Accept()
		if err != nil {
			srv.logger.Println("Error accepting connection:", err)
			continue
		}

		go srv.handleConnection(conn)
	}
}

func (srv *Server) Stop() {
	(*srv.s).Close()
}

func (srv *Server) handleConnection(conn net.Conn) {
	srv.logger.Println("New connection:", conn.RemoteAddr())

	defer func() {
		conn.Close()
	}()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				srv.logger.Println("Connection closed by client:", conn.RemoteAddr())
				break
			}

			srv.logger.Println("Error reading data:", err)
			break
		}

		fmt.Println("Received data:", string(buf[:n]))

		parsedRequest, err := srv.ParseRequest(buf[:n])
		if err != nil {
			srv.logger.Println("Error parsing request:", err)

			if _, err := conn.Write([]byte(err.Error())); err != nil {
				srv.logger.Println("Error writing response:", err)
				break
			}

			continue
		}

		fmt.Println("Parsed request:", parsedRequest)

		if _, err := conn.Write([]byte(parsedRequest.Stringify())); err != nil {
			srv.logger.Println("Error writing response:", err)
		}

	}

}

func (srv *Server) ParseRequest(rawData []byte) (ParsedRequest, error) {
	if len(rawData) < 5 {
		return ParsedRequest{}, errors.New("invalid data format")
	}

	// 4 bytes for length of whole command
	// length := rawData[:4]
	// 4 bytes for protocol version
	protocolVersion := rawData[:4]
	// 4 bytes for message type
	messageType := rawData[4:5]
	// rest of the bytes for data
	data := rawData[5:]

	if string(protocolVersion) != "0001" {
		return ParsedRequest{}, errors.New("invalid protocol version")
	}

	// intLength, err := strconv.Atoi(string(length))
	// if err != nil {
	// 	log.Println("Error converting length to integer:", err)
	// 	return ParsedRequest{}, errors.New("error converting length to integer")
	// }

	stringMessageType := string(messageType)

	return ParsedRequest{
		// length:          intLength,
		protocolVersion: string(protocolVersion),
		messageType:     stringMessageType,
		data:            string(data),
	}, nil

}
