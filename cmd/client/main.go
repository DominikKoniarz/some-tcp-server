package main

import (
	"github.com/DominikKoniarz/some-tcp-server/internal/client"
)

func main() {
	client := client.NewClient()
	client.Connect()
}
