func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer server.Close()

	fmt.Println("Server is running on port 8080")

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Client connected")

		go handleClient(&client)
	}
}

func handleClient(clientPtr *net.Conn) {
	client := *clientPtr

	for {
		buff := make([]byte, 1024)
		_, err := client.Read(buff)

		if err != nil {
			if err.Error() == "EOF" {
				client.Close()
				fmt.Println("Client disconnected")
				break
			}

			fmt.Println(err)
			return
		}

		fmt.Println(string(buff))

		serverResponse := append([]byte("Message received: "), buff...)
		_, err = client.Write(serverResponse)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
