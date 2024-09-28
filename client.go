package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //current state
}

func NewClient(serverIp string, serverPort int) *Client {
	//create a new client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. public chat")
	fmt.Println("2. private chat")
	fmt.Println("3. update name")
	fmt.Println("0. exit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please input right number")
		return false
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("please input chatMsg, '/exit' for exit")
	fmt.Scanln(&chatMsg)
	for chatMsg != "/exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	fmt.Println("here is the user list: ")
	client.SelectUser()
	fmt.Println("please input the name you want to chat with, '/exit' for exit: ")
	fmt.Scanln(&remoteName)

	for remoteName != "/exit" {
		if remoteName == client.Name {
			fmt.Println("can't chat with yourself")
			continue
		}

		fmt.Println("to " + remoteName + ": ")
		fmt.Scanln(&chatMsg)
		for chatMsg != "/exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("to " + remoteName + ": ")
			fmt.Scanln(&chatMsg)
		}
		remoteName = ""
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("please input your name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return false
	}
	return true
}

func (client *Client) run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("public chat")
			client.PublicChat()
			break
		case 2:
			fmt.Println("private chat")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("update name")
			client.UpdateName()
			break
		}
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
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

	//start a goroutine to deal with server response
	go client.DealResponse()

	fmt.Println("Connected to server")
	client.run()
}
