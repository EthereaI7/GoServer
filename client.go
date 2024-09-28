package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//create a new client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	//connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn
	return client
}

var ServerIp string
var ServerPort int

// ./server -ip 192.168.205.156 -port 8888
func init() {
	flag.StringVar(&ServerIp, "ip", "192.168.205.156", "set server IP(default:192.168.205.156)")
	flag.IntVar(&ServerPort, "port", 8888, "set server port(default:8888)")
}

func main() {
	flag.Parse()

	client := NewClient(ServerIp, ServerPort)
	if client == nil {
		fmt.Println("Failed to connect to server")
		return
	}

	fmt.Println("Connected to server")

}
